package sssp

import (
	"container/heap"
	"math"

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

	return &Solver{
		G:    g,
		Dist: make(DistMap, g.V),
		K:    k,
		T:    t,
	}
}

func (s *Solver) Run(source int) []float64 {
	for i := range s.Dist {
		s.Dist[i] = Infinity
	}
	s.Dist[source] = 0

	// Calculate Max Level l = ceil(log n / t)
	n := float64(s.G.V)
	l := int(math.Ceil(math.Log(n) / float64(s.T)))

	// Initial call
	// S = {source}, B = Infinity
	S := []int{source}
	s.BMSSP(l, Infinity, S)

	return s.Dist
}

// BMSSP (Bounded Multi-Source Shortest Path) - Algorithm 3
func (s *Solver) BMSSP(l int, B float64, S []int) (float64, []int) {
	// Base Case
	if l == 0 {
		// In base case, S is singleton {x} (Lemma 3.7 proof implies decomposition, but algo 3 line 1 says check l=0)
		// Wait, Algorithm 3 Line 2 checks l=0.
		// If l=0, S should be small? The paper says "S is a singleton".
		// But in recursion we might pass a set S_i.
		// However, the paper implies we reduce down to base cases.
		// If S is not singleton here, we iterate?
		// Actually, standard BaseCase takes a set S? Algorithm 2 signature: BASECASE(B, S).
		// Requirement 1 of BaseCase: S={x}.
		// So if we reach l=0 with |S| > 1, we must have split it?
		// No, the recursion structure ensures |S| reduces.
		// Actually, at l=0, we just run mini-Dijkstra from ALL sources in S?
		// Let's assume S is handled by running BaseCase on the single source, or we iterate.
		// Algorithm 3 calls BASECASE(B,S) directly. Algorithm 2 expects singleton.
		// Implication: At level 0, |S| is small enough or just 1.
		// Let's implement BaseCase to handle set S generally (Multi-source Dijkstra) to be safe.
		return s.BaseCase(B, S)
	}

	// 1. Find Pivots
	P, W := s.FindPivots(B, S)

	// 2. Initialize DS
	// M = 2^((l-1)t)
	pow := math.Pow(2, float64((l-1)*s.T))
	M := int(pow)
	if M < 1 {
		M = 1
	}

	D := ds.NewDataStructure(M)

	// Insert P into D
	for _, x := range P {
		D.Insert(x, s.Dist[x])
	}

	// 3. Loop
	U := make(map[int]bool) // Using map for set U
	if len(P) > 0 {
		minP := Infinity
		for _, x := range P {
			if s.Dist[x] < minP {
				minP = s.Dist[x]
			}
		}

		limit := s.K * int(math.Pow(2, float64(l*s.T))) // k * 2^(lt)

		for len(U) < limit && D.Count > 0 {
			// Pull
			Si_items, Bi := D.Pull()
			Si := make([]int, len(Si_items))
			for idx, item := range Si_items {
				Si[idx] = item.Key
			}

			// Recursive Call
			Bi_prime, Ui := s.BMSSP(l-1, Bi, Si)

			// Add Ui to U
			for _, u := range Ui {
				U[u] = true
			}

			// Relax edges from Ui
			// K set for batch prepend
			var K []ds.Item

			for _, u := range Ui {
				for _, edge := range s.G.Adj[u] {
					v := edge.To
					w := edge.Weight
					newDist := s.Dist[u] + w

					if newDist <= s.Dist[v] { // Relax
						s.Dist[v] = newDist

						// Insert logic
						if newDist >= Bi && newDist < B {
							D.Insert(v, newDist)
						} else if newDist >= Bi_prime && newDist < Bi {
							K = append(K, ds.Item{Key: v, Value: newDist})
						}
					}
				}
			}

			// Batch Prepend K + specific Si
			// "Batch Prepend all records in K and <x, d[x]> for x in Si with d[x] in [Bi_prime, Bi)"
			var batch []ds.Item
			batch = append(batch, K...)
			for _, x := range Si {
				if s.Dist[x] >= Bi_prime && s.Dist[x] < Bi {
					batch = append(batch, ds.Item{Key: x, Value: s.Dist[x]})
				}
			}
			D.BatchPrepend(batch)

			// Check large workload
			if len(U) > limit {
				B_final := Bi_prime
				// Add W filtered
				finalU := make([]int, 0, len(U))
				for u := range U {
					finalU = append(finalU, u)
				}
				for _, w := range W {
					if s.Dist[w] < B_final {
						if !U[w] {
							finalU = append(finalU, w)
						}
					}
				}
				return B_final, finalU
			}
		}

	} // close len(P)>0 guard

	// Successful execution
	B_final := B

	// Add W filtered by B_final (which is B)
	finalU := make([]int, 0, len(U))
	for u := range U {
		finalU = append(finalU, u)
	}
	for _, w := range W {
		if s.Dist[w] < B_final {
			if !U[w] {
				finalU = append(finalU, w)
			}
		}
	}

	return B_final, finalU
}

// FindPivots - Algorithm 1
func (s *Solver) FindPivots(B float64, S []int) ([]int, []int) {
	// Optimization: Use slices instead of maps for visited/sets
	inW := make([]bool, s.G.V)
	for _, x := range S {
		inW[x] = true
	}

	Wi_prev := S
	W_list := make([]int, len(S)) // Keep track of W as a list
	copy(W_list, S)

	// Relax k steps
	for i := 1; i <= s.K; i++ {
		Wi := make([]int, 0)
		for _, u := range Wi_prev {
			for _, edge := range s.G.Adj[u] {
				v := edge.To
				w := edge.Weight

				// Relax
				if s.Dist[u]+w < s.Dist[v] { // Strict inequality for update
					// Note: Paper says <= for validity, but < for update?
					// Standard Dijkstra relaxes on <.
					// However, we must allow equality to build the forest F later.
					// We update if strictly better.
					s.Dist[v] = s.Dist[u] + w
					// Add to Wi if < B
					if s.Dist[v] < B {
						if !inW[v] {
							Wi = append(Wi, v)
							inW[v] = true
							W_list = append(W_list, v)
						}
					}
				}
			}
		}

		if len(W_list) > s.K*len(S) {
			// Return P=S, W=W_list
			P := make([]int, len(S))
			copy(P, S)
			return P, W_list
		}
		Wi_prev = Wi
	}

	// Construct Forest F
	// F is defined on W. We need to identify roots in S.
	// Root logic: u is root if subtree size >= k.

	// We compute size of tree rooted at u.
	// Memoization: 0 = unvisited, -1 = being processed (cycle detection), positive = computed size
	memoSize := make([]int, s.G.V)

	var calcSize func(u int) int
	calcSize = func(u int) int {
		// Already computed
		if memoSize[u] > 0 {
			return memoSize[u]
		}

		// Cycle detection: if we encounter a node being processed, there's a cycle
		// In a proper shortest path forest, this shouldn't happen, but handle it gracefully
		if memoSize[u] == -1 {
			return 1 // Treat as single node to avoid infinite recursion
		}

		// Mark as being processed
		memoSize[u] = -1

		count := 1
		// Iterate outgoing edges in F
		// Edge (u, v) is in F if u,v in W and dist[v] == dist[u] + w
		for _, edge := range s.G.Adj[u] {
			v := edge.To
			if inW[v] && math.Abs(s.Dist[v]-(s.Dist[u]+edge.Weight)) < 1e-9 {
				count += calcSize(v)
			}
		}

		// Store computed size
		memoSize[u] = count
		return count
	}

	P := make([]int, 0)
	for _, u := range S {
		if calcSize(u) >= s.K {
			P = append(P, u)
		}
	}

	return P, W_list
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

		for _, edge := range s.G.Adj[u] {
			v := edge.To
			w := edge.Weight
			if s.Dist[u]+w <= s.Dist[v] && s.Dist[u]+w < B {
				s.Dist[v] = s.Dist[u] + w
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
