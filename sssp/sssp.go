package sssp

import (
	"container/heap"
	"math"
	"runtime"
	"sync"

	"github.com/phr3nzy/duan-sssp/ds"
	"github.com/phr3nzy/duan-sssp/graph"
)

// Algorithm Constants
const Infinity = math.MaxFloat64

// DistMap holds current distance estimates.
type DistMap []float64

// PriorityQueue for BaseCase
type PQItem struct {
	u        int
	priority float64
	index    int
}
type PriorityQueue []*PQItem

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i]; pq[i].index = i; pq[j].index = j }
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*PQItem)
	item.index = len(*pq)
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Solver encapsulates the algorithm state.
type Solver struct {
	G    *graph.Graph
	Dist DistMap
	K    int
	T    int

	// Pre-allocated buffers for performance
	bufInt   []int
	bufItem  []ds.Item
	bufBatch []ds.Item

	// Parallel processing
	workerPool chan struct{}
	numWorkers int

	// Event listener for visualization
	listener EventListener
}

func NewSolver(g *graph.Graph) *Solver {
	n := float64(g.V)
	logN := math.Log(n)
	// k = floor(log^(1/3) n)
	k := int(math.Floor(math.Pow(logN, 1.0/3.0)))
	if k < 2 {
		k = 2
	}

	// t = floor(log^(2/3) n)
	t := int(math.Floor(math.Pow(logN, 2.0/3.0)))
	if t < 2 {
		t = 2
	}

	numWorkers := runtime.NumCPU()
	if numWorkers > 8 {
		numWorkers = 8 // Cap to avoid excessive contention
	}

	return &Solver{
		G:          g,
		Dist:       make(DistMap, g.V),
		K:          k,
		T:          t,
		bufInt:     make([]int, 0, 1000),
		bufItem:    make([]ds.Item, 0, 1000),
		bufBatch:   make([]ds.Item, 0, 1000),
		workerPool: make(chan struct{}, numWorkers),
		numWorkers: numWorkers,
		listener:   &NoOpListener{},
	}
}

// SetEventListener sets the event listener for visualization
func (s *Solver) SetEventListener(listener EventListener) {
	if listener != nil {
		s.listener = listener
	} else {
		s.listener = &NoOpListener{}
	}
}

func (s *Solver) Run(source int) []float64 {
	for i := range s.Dist {
		s.Dist[i] = Infinity
	}
	s.Dist[source] = 0
	s.listener.OnNodeDiscovered(source, 0)

	// Calculate Max Level l = ceil(log n / t)
	n := float64(s.G.V)
	l := int(math.Ceil(math.Log(n) / float64(s.T)))

	// Initial call
	// S = {source}, B = Infinity
	S := []int{source}
	s.listener.OnPhaseChange("BMSSP", l)
	s.BMSSP(l, Infinity, S)

	return s.Dist
}

// BMSSP (Bounded Multi-Source Shortest Path) - Algorithm 3
func (s *Solver) BMSSP(l int, B float64, S []int) (float64, []int) {
	s.listener.OnPhaseChange("BMSSP", l)

	if l == 0 {
		s.listener.OnPhaseChange("BaseCase", 0)
		return s.BaseCase(B, S)
	}

	s.listener.OnPhaseChange("FindPivots", l)
	P, W := s.FindPivots(B, S)

	if len(P) == 0 {
		return s.finalizeBMSSP(B, W, make(map[int]bool))
	}

	D := s.initializeDataStructure(l, P)
	U, _ := s.processMainLoop(l, B, D, W)

	return s.finalizeBMSSP(B, W, U)
}

// initializeDataStructure creates and populates the data structure for BMSSP
func (s *Solver) initializeDataStructure(l int, P []int) *ds.DataStructure {
	M := int(math.Pow(2, float64((l-1)*s.T)))
	if M < 1 {
		M = 1
	}

	D := ds.NewDataStructure(M)
	for _, x := range P {
		D.Insert(x, s.Dist[x])
	}

	return D
}

// processMainLoop handles the main iteration loop of BMSSP
func (s *Solver) processMainLoop(l int, B float64, D *ds.DataStructure, W []int) (map[int]bool, int) {
	U := make(map[int]bool)
	limit := s.K * int(math.Pow(2, float64(l*s.T)))

	for len(U) < limit && D.Count > 0 {
		Si, Bi := s.pullAndExtract(D)
		Bi_prime, Ui := s.BMSSP(l-1, Bi, Si)

		s.addToSet(U, Ui)
		K := s.relaxEdges(Ui, Bi, Bi_prime, B, D)
		s.batchPrepend(D, K, Si, Bi_prime, Bi)

		if len(U) > limit {
			return s.buildFinalSet(U, W, Bi_prime), limit
		}
	}

	return U, limit
}

// pullAndExtract pulls items from data structure and extracts keys
func (s *Solver) pullAndExtract(D *ds.DataStructure) ([]int, float64) {
	items, Bi := D.Pull()
	// Return slice directly without copying - caller shouldn't modify
	Si := make([]int, len(items))
	for i, item := range items {
		Si[i] = item.Key
	}
	return Si, Bi
}

// addToSet adds elements from Ui to U
func (s *Solver) addToSet(U map[int]bool, Ui []int) {
	for _, u := range Ui {
		U[u] = true
	}
}

// relaxEdges performs edge relaxation and returns items for batch prepend
func (s *Solver) relaxEdges(Ui []int, Bi, Bi_prime, B float64, D *ds.DataStructure) []ds.Item {
	if len(Ui) == 0 {
		return nil
	}

	// For small workloads, use sequential processing
	if len(Ui) <= 4 || s.numWorkers == 1 {
		return s.relaxEdgesSequential(Ui, Bi, Bi_prime, B, D)
	}

	return s.relaxEdgesParallel(Ui, Bi, Bi_prime, B, D)
}

// relaxEdgesSequential processes edges sequentially
func (s *Solver) relaxEdgesSequential(Ui []int, Bi, Bi_prime, B float64, D *ds.DataStructure) []ds.Item {
	var K []ds.Item

	for _, u := range Ui {
		for _, edge := range s.G.Adj[u] {
			newDist := s.Dist[u] + edge.Weight

			if newDist <= s.Dist[edge.To] {
				oldDist := s.Dist[edge.To]
				s.Dist[edge.To] = newDist

				if oldDist == Infinity {
					s.listener.OnNodeDiscovered(edge.To, newDist)
				} else {
					s.listener.OnNodeRelaxed(u, edge.To, oldDist, newDist)
				}

				if newDist >= Bi && newDist < B {
					D.Insert(edge.To, newDist)
				} else if newDist >= Bi_prime && newDist < Bi {
					K = append(K, ds.Item{Key: edge.To, Value: newDist})
				}
			}
		}
	}

	return K
}

// relaxEdgesParallel processes edges in parallel using worker pool
func (s *Solver) relaxEdgesParallel(Ui []int, Bi, Bi_prime, B float64, D *ds.DataStructure) []ds.Item {
	var wg sync.WaitGroup
	results := make([][]ds.Item, len(Ui))

	// Process each vertex in parallel
	for i, u := range Ui {
		wg.Add(1)
		go func(vertexIdx, vertex int) {
			defer wg.Done()

			var localK []ds.Item
			for _, edge := range s.G.Adj[vertex] {
				newDist := s.Dist[vertex] + edge.Weight

				if newDist <= s.Dist[edge.To] {
					oldDist := s.Dist[edge.To]
					s.Dist[edge.To] = newDist

					if oldDist == Infinity {
						s.listener.OnNodeDiscovered(edge.To, newDist)
					} else {
						s.listener.OnNodeRelaxed(vertex, edge.To, oldDist, newDist)
					}

					if newDist >= Bi && newDist < B {
						D.Insert(edge.To, newDist)
					} else if newDist >= Bi_prime && newDist < Bi {
						localK = append(localK, ds.Item{Key: edge.To, Value: newDist})
					}
				}
			}
			results[vertexIdx] = localK
		}(i, u)
	}

	wg.Wait()

	// Merge results
	var totalK []ds.Item
	for _, k := range results {
		totalK = append(totalK, k...)
	}

	return totalK
}

// batchPrepend prepares and adds batch items to data structure
func (s *Solver) batchPrepend(D *ds.DataStructure, K []ds.Item, Si []int, Bi_prime, Bi float64) {
	// Reuse batch buffer
	s.bufBatch = s.bufBatch[:0]
	if cap(s.bufBatch) < len(K)+len(Si) {
		s.bufBatch = make([]ds.Item, 0, (len(K)+len(Si))*2)
	}

	s.bufBatch = append(s.bufBatch, K...)

	for _, x := range Si {
		if s.Dist[x] >= Bi_prime && s.Dist[x] < Bi {
			s.bufBatch = append(s.bufBatch, ds.Item{Key: x, Value: s.Dist[x]})
		}
	}

	D.BatchPrepend(s.bufBatch)
}

// buildFinalSet constructs the final vertex set with filtering
func (s *Solver) buildFinalSet(U map[int]bool, W []int, bound float64) map[int]bool {
	for _, w := range W {
		if s.Dist[w] < bound && !U[w] {
			U[w] = true
		}
	}
	return U
}

// finalizeBMSSP converts the result set to final format
func (s *Solver) finalizeBMSSP(B float64, W []int, U map[int]bool) (float64, []int) {
	finalU := make([]int, 0, len(U))
	for u := range U {
		finalU = append(finalU, u)
	}

	for _, w := range W {
		if s.Dist[w] < B && !U[w] {
			finalU = append(finalU, w)
		}
	}

	return B, finalU
}

// FindPivots - Algorithm 1
func (s *Solver) FindPivots(B float64, S []int) ([]int, []int) {
	inW := make([]bool, s.G.V)
	for _, x := range S {
		inW[x] = true
	}

	W_list := make([]int, len(S))
	copy(W_list, S)

	// Relax k steps
	W_list = s.relaxKSteps(B, S, inW, W_list)

	// If W grew too large, return early
	if len(W_list) > s.K*len(S) {
		P := make([]int, len(S))
		copy(P, S)
		return P, W_list
	}

	// Compute pivots from tree sizes
	P := s.computePivots(S, inW)
	return P, W_list
}

// relaxKSteps performs k relaxation steps from source set
func (s *Solver) relaxKSteps(B float64, S []int, inW []bool, W_list []int) []int {
	Wi_prev := S

	for i := 1; i <= s.K; i++ {
		Wi := make([]int, 0)

		for _, u := range Wi_prev {
			for _, edge := range s.G.Adj[u] {
				newDist := s.Dist[u] + edge.Weight

				if newDist < s.Dist[edge.To] {
					oldDist := s.Dist[edge.To]
					s.Dist[edge.To] = newDist

					if oldDist == Infinity {
						s.listener.OnNodeDiscovered(edge.To, newDist)
					} else {
						s.listener.OnNodeRelaxed(u, edge.To, oldDist, newDist)
					}

					if newDist < B && !inW[edge.To] {
						Wi = append(Wi, edge.To)
						inW[edge.To] = true
						W_list = append(W_list, edge.To)
					}
				}
			}
		}

		if len(W_list) > s.K*len(S) {
			return W_list
		}
		Wi_prev = Wi
	}

	return W_list
}

// computePivots identifies pivots based on tree sizes
func (s *Solver) computePivots(S []int, inW []bool) []int {
	memoSize := make([]int, s.G.V)

	calcSize := s.makeTreeSizeCalculator(inW, memoSize)

	P := make([]int, 0)
	for _, u := range S {
		if calcSize(u) >= s.K {
			P = append(P, u)
		}
	}

	return P
}

// makeTreeSizeCalculator creates a function to calculate tree sizes with cycle detection
func (s *Solver) makeTreeSizeCalculator(inW []bool, memoSize []int) func(int) int {
	var calcSize func(u int) int

	calcSize = func(u int) int {
		if memoSize[u] > 0 {
			return memoSize[u]
		}

		if memoSize[u] == -1 {
			return 1 // Cycle detected
		}

		memoSize[u] = -1
		count := 1 + s.countTreeChildren(u, inW, calcSize)
		memoSize[u] = count

		return count
	}

	return calcSize
}

// countTreeChildren counts children in the shortest path forest
func (s *Solver) countTreeChildren(u int, inW []bool, calcSize func(int) int) int {
	count := 0

	for _, edge := range s.G.Adj[u] {
		v := edge.To
		if inW[v] && math.Abs(s.Dist[v]-(s.Dist[u]+edge.Weight)) < 1e-9 {
			count += calcSize(v)
		}
	}

	return count
}

// BaseCase - Algorithm 2
func (s *Solver) BaseCase(B float64, S []int) (float64, []int) {
	U0 := make(map[int]bool)
	pq := &PriorityQueue{}
	heap.Init(pq)

	for _, x := range S {
		U0[x] = true
		heap.Push(pq, &PQItem{u: x, priority: s.Dist[x]})
	}

	limit := s.K + 1

	for pq.Len() > 0 && len(U0) < limit {
		item := heap.Pop(pq).(*PQItem)
		u := item.u

		// If popped distance > current dist, ignore (stale)
		if item.priority > s.Dist[u] {
			continue
		}

		U0[u] = true // Add to set
		s.listener.OnIterationComplete(len(U0))

		for _, edge := range s.G.Adj[u] {
			v := edge.To
			w := edge.Weight
			if s.Dist[u]+w <= s.Dist[v] && s.Dist[u]+w < B {
				oldDist := s.Dist[v]
				s.Dist[v] = s.Dist[u] + w

				if oldDist == Infinity {
					s.listener.OnNodeDiscovered(v, s.Dist[v])
				} else {
					s.listener.OnNodeRelaxed(u, v, oldDist, s.Dist[v])
				}

				heap.Push(pq, &PQItem{u: v, priority: s.Dist[v]})
			}
		}
	}

	uList := make([]int, 0, len(U0))
	for u := range U0 {
		uList = append(uList, u)
	}

	if len(U0) <= s.K {
		return B, uList
	}

	// Return max dist in U0 as B'
	maxD := 0.0
	for u := range U0 {
		if s.Dist[u] > maxD {
			maxD = s.Dist[u]
		}
	}

	// Filter U: {v in U0 : d[v] < B'}
	finalU := make([]int, 0)
	for u := range U0 {
		if s.Dist[u] < maxD {
			finalU = append(finalU, u)
		}
	}
	return maxD, finalU
}
