package sssp

import (
	"container/heap"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/phr3nzy/duan-sssp/graph"
)

// BenchmarkSSSP runs benchmarks for various graph sizes
func BenchmarkSSSP(b *testing.B) {
	testCases := []struct {
		name     string
		vertices int
		edges    int
	}{
		{"Small_V1K_E3K", 1000, 3000},
		{"Medium_V5K_E15K", 5000, 15000},
		{"Large_V10K_E30K", 10000, 30000},
		{"VeryLarge_V50K_E150K", 50000, 150000},
		{"Huge_V100K_E300K", 100000, 300000},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			// Generate graph once
			g := generateRandomGraph(tc.vertices, tc.edges)
			tg := g.ToConstantDegree()
			solver := NewSolver(tg.G)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				solver.Run(tg.OriginalTo[0])
			}
		})
	}
}

// BenchmarkSSSPDensity benchmarks different graph densities
func BenchmarkSSSPDensity(b *testing.B) {
	vertices := 10000
	densities := []struct {
		name       string
		edgeFactor int // edges = vertices * edgeFactor
	}{
		{"Sparse_2x", 2},
		{"Medium_5x", 5},
		{"Dense_10x", 10},
		{"VeryDense_20x", 20},
	}

	for _, d := range densities {
		edges := vertices * d.edgeFactor
		b.Run(d.name, func(b *testing.B) {
			g := generateRandomGraph(vertices, edges)
			tg := g.ToConstantDegree()
			solver := NewSolver(tg.G)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				solver.Run(tg.OriginalTo[0])
			}
		})
	}
}

// BenchmarkTransformation benchmarks just the graph transformation step
func BenchmarkTransformation(b *testing.B) {
	testCases := []struct {
		name     string
		vertices int
		edges    int
	}{
		{"Small_V1K_E3K", 1000, 3000},
		{"Medium_V10K_E30K", 10000, 30000},
		{"Large_V50K_E150K", 50000, 150000},
		{"Huge_V100K_E300K", 100000, 300000},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			g := generateRandomGraph(tc.vertices, tc.edges)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				g.ToConstantDegree()
			}
		})
	}
}

// BenchmarkFindPivots benchmarks the pivot finding algorithm
func BenchmarkFindPivots(b *testing.B) {
	vertices := 1000 // Reduced size to prevent long-running benchmark
	edges := 3000
	g := generateRandomGraph(vertices, edges)
	tg := g.ToConstantDegree()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Create fresh solver for each iteration to avoid state pollution
		solver := NewSolver(tg.G)
		for j := range solver.Dist {
			solver.Dist[j] = Infinity
		}
		solver.Dist[0] = 0
		S := []int{0}
		b.StartTimer()

		solver.FindPivots(Infinity, S)
	}
}

// BenchmarkBaseCase benchmarks the base case algorithm
func BenchmarkBaseCase(b *testing.B) {
	vertices := 100 // Small size for isolated base case testing
	edges := 300
	g := generateRandomGraph(vertices, edges)
	tg := g.ToConstantDegree()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Create fresh solver for each iteration to avoid state pollution
		solver := NewSolver(tg.G)
		for j := range solver.Dist {
			solver.Dist[j] = Infinity
		}
		solver.Dist[0] = 0
		S := []int{0}
		b.StartTimer()

		solver.BaseCase(Infinity, S)
	}
}

// BenchmarkComparison compares different SSSP algorithms
func BenchmarkComparison(b *testing.B) {
	vertices := 10000
	edges := 30000

	b.Run("DuanAlgorithm", func(b *testing.B) {
		g := generateRandomGraph(vertices, edges)
		tg := g.ToConstantDegree()
		solver := NewSolver(tg.G)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			solver.Run(tg.OriginalTo[0])
		}
	})

	b.Run("AStar", func(b *testing.B) {
		g := generateRandomGraph(vertices, edges)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			aStarSSP(g, 0)
		}
	})

	b.Run("NaiveDijkstra", func(b *testing.B) {
		g := generateRandomGraph(vertices, edges)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			naiveDijkstra(g, 0)
		}
	})
}

// BenchmarkAlgorithmComparison provides detailed size-based comparison
func BenchmarkAlgorithmComparison(b *testing.B) {
	sizes := []struct {
		name     string
		vertices int
		edges    int
	}{
		{"Small_V1K_E3K", 1000, 3000},
		{"Medium_V5K_E15K", 5000, 15000},
		{"Large_V10K_E30K", 10000, 30000},
	}

	for _, sz := range sizes {
		b.Run(sz.name, func(b *testing.B) {
			b.Run("Duan", func(b *testing.B) {
				g := generateRandomGraph(sz.vertices, sz.edges)
				tg := g.ToConstantDegree()
				solver := NewSolver(tg.G)

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					solver.Run(tg.OriginalTo[0])
				}
			})

			b.Run("AStar", func(b *testing.B) {
				g := generateRandomGraph(sz.vertices, sz.edges)

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					aStarSSP(g, 0)
				}
			})
		})
	}
}

// Helper function to generate random graphs
func generateRandomGraph(vertices, edges int) *graph.Graph {
	g := graph.NewGraph(vertices)
	// Use Go 1.20+ random number generation without deprecated Seed
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // Deterministic random for benchmarks

	for i := 0; i < edges; i++ {
		u := rng.Intn(vertices)
		v := rng.Intn(vertices)
		if u == v {
			v = (v + 1) % vertices // Avoid self-loops
		}
		w := rng.Float64()*100.0 + 1.0 // Weights 1-101
		g.AddEdge(u, v, w)
	}

	return g
}

// naiveDijkstra implements standard Dijkstra's algorithm for comparison
func naiveDijkstra(g *graph.Graph, source int) []float64 {
	dist := make([]float64, g.V)
	visited := make([]bool, g.V)

	for i := range dist {
		dist[i] = Infinity
	}
	dist[source] = 0

	for count := 0; count < g.V; count++ {
		// Find minimum distance vertex
		minDist := Infinity
		u := -1
		for v := 0; v < g.V; v++ {
			if !visited[v] && dist[v] < minDist {
				minDist = dist[v]
				u = v
			}
		}

		if u == -1 {
			break
		}

		visited[u] = true

		// Relax edges
		for _, edge := range g.Adj[u] {
			v := edge.To
			w := edge.Weight
			if dist[u]+w < dist[v] {
				dist[v] = dist[u] + w
			}
		}
	}

	return dist
}

// aStarNode represents a node in A* priority queue
type aStarNode struct {
	vertex int
	gScore float64 // Distance from source
	fScore float64 // gScore + heuristic
	index  int
}

type aStarPQ []*aStarNode

func (pq aStarPQ) Len() int           { return len(pq) }
func (pq aStarPQ) Less(i, j int) bool { return pq[i].fScore < pq[j].fScore }
func (pq aStarPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *aStarPQ) Push(x interface{}) {
	n := len(*pq)
	node := x.(*aStarNode)
	node.index = n
	*pq = append(*pq, node)
}
func (pq *aStarPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*pq = old[0 : n-1]
	return node
}

// aStarSSP implements A* algorithm for single-source shortest paths
// Uses zero heuristic (h(v) = 0), which makes it equivalent to Dijkstra with heap
// but demonstrates A* structure for comparison
func aStarSSP(g *graph.Graph, source int) []float64 {
	dist := make([]float64, g.V)
	for i := range dist {
		dist[i] = Infinity
	}
	dist[source] = 0

	// Priority queue ordered by f-score
	pq := &aStarPQ{}
	heap.Init(pq)
	heap.Push(pq, &aStarNode{
		vertex: source,
		gScore: 0,
		fScore: 0, // 0 + heuristic(source)
	})

	visited := make([]bool, g.V)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*aStarNode)
		u := current.vertex

		// Skip if already processed
		if visited[u] {
			continue
		}
		visited[u] = true

		// Skip if this is a stale entry
		if current.gScore > dist[u] {
			continue
		}

		// Relax edges
		for _, edge := range g.Adj[u] {
			v := edge.To
			w := edge.Weight

			if dist[u]+w < dist[v] {
				dist[v] = dist[u] + w

				// Heuristic: h(v) = 0 (admissible, consistent)
				// For a real application, this could be Euclidean distance to target
				heuristic := 0.0

				heap.Push(pq, &aStarNode{
					vertex: v,
					gScore: dist[v],
					fScore: dist[v] + heuristic,
				})
			}
		}
	}

	return dist
}

// TestBasicExecution tests that the algorithm runs without crashing
func TestBasicExecution(t *testing.T) {
	testCases := []struct {
		name     string
		vertices int
		edges    int
	}{
		{"Tiny", 10, 30},
		{"Small", 100, 300},
		{"Medium", 500, 1500},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := generateRandomGraph(tc.vertices, tc.edges)
			tg := g.ToConstantDegree()
			solver := NewSolver(tg.G)
			rawDist := solver.Run(tg.OriginalTo[0])
			distances := tg.MapDistances(rawDist)

			// Basic sanity checks
			if distances[0] != 0 {
				t.Errorf("Source distance should be 0, got %f", distances[0])
			}

			// Check that some vertices are reachable
			reachable := 0
			for _, d := range distances {
				if d < Infinity {
					reachable++
				}
			}

			if reachable == 0 {
				t.Error("No vertices reachable from source")
			}

			t.Logf("Reachable vertices: %d/%d", reachable, tc.vertices)
		})
	}
}

// BenchmarkScalability tests scalability with increasing graph sizes
func BenchmarkScalability(b *testing.B) {
	// Limited sizes to avoid stack overflow in deep recursion
	// For very large graphs, the tree size calculation can hit recursion limits
	sizes := []int{500, 1000, 2000, 5000, 10000}

	for _, size := range sizes {
		edges := size * 3 // Sparse graph
		b.Run(fmt.Sprintf("V%d_E%d", size, edges), func(b *testing.B) {
			g := generateRandomGraph(size, edges)
			tg := g.ToConstantDegree()
			solver := NewSolver(tg.G)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				solver.Run(tg.OriginalTo[0])
			}
		})
	}
}

// BenchmarkMemoryUsage tests memory usage patterns
func BenchmarkMemoryUsage(b *testing.B) {
	vertices := 10000
	edges := 30000

	b.Run("WithTransform", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			g := generateRandomGraph(vertices, edges)
			tg := g.ToConstantDegree()
			solver := NewSolver(tg.G)
			solver.Run(tg.OriginalTo[0])
		}
	})
}
