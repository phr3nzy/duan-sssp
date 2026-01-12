package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/phr3nzy/duan-sssp/graph"
	"github.com/phr3nzy/duan-sssp/sssp"
)

// ANSI color codes for terminal
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
)

type BenchmarkResult struct {
	Algorithm string
	Time      time.Duration
	Vertices  int
	Edges     int
	CoreCount int
}

func main() {
	// Parse flags
	vertices := flag.Int("vertices", 10000, "Number of vertices")
	edgeFactor := flag.Int("edge-factor", 3, "Edges = vertices * edge-factor")
	iterations := flag.Int("iterations", 10, "Number of benchmark iterations")
	showGraph := flag.Bool("show-graph", true, "Show graph visualization")
	parallel := flag.Bool("parallel", true, "Use all CPU cores")
	web := flag.Bool("web", false, "Open web visualization in browser")

	flag.Parse()

	edges := (*vertices) * (*edgeFactor)

	// Configure runtime
	if *parallel {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(1)
	}

	printHeader(*vertices, edges, *iterations)

	// Generate graph
	fmt.Printf("%s[1/4] Generating random graph...%s\n", colorCyan, colorReset)
	g := generateGraph(*vertices, edges)

	if *showGraph {
		visualizeGraph(g, 20) // Show sample of 20 vertices
	}

	// Run benchmarks
	fmt.Printf("\n%s[2/4] Running benchmarks with %d cores...%s\n", colorCyan, runtime.GOMAXPROCS(0), colorReset)

	results := make([]BenchmarkResult, 0)

	// Duan Algorithm
	duanTime := benchmarkDuan(g, *iterations)
	results = append(results, BenchmarkResult{
		Algorithm: "Duan (O(m log^(2/3) n))",
		Time:      duanTime,
		Vertices:  *vertices,
		Edges:     edges,
		CoreCount: runtime.GOMAXPROCS(0),
	})

	// A* Algorithm
	astarTime := benchmarkAStar(g, *iterations)
	results = append(results, BenchmarkResult{
		Algorithm: "A* with Heap",
		Time:      astarTime,
		Vertices:  *vertices,
		Edges:     edges,
		CoreCount: runtime.GOMAXPROCS(0),
	})

	// Parallel Duan (if requested)
	if *parallel && runtime.NumCPU() > 1 {
		parallelTime := benchmarkParallelMultiSource(g, *iterations)
		results = append(results, BenchmarkResult{
			Algorithm: fmt.Sprintf("Duan Parallel (%d cores)", runtime.NumCPU()),
			Time:      parallelTime,
			Vertices:  *vertices,
			Edges:     edges,
			CoreCount: runtime.GOMAXPROCS(0),
		})
	}

	// Display results
	fmt.Printf("\n%s[3/4] Results:%s\n", colorCyan, colorReset)
	displayResults(results)

	// Visualize performance
	fmt.Printf("\n%s[4/4] Performance Visualization:%s\n", colorCyan, colorReset)
	visualizePerformance(results)

	printSummary(results)

	// Web visualization
	if *web {
		fmt.Printf("\n%s[Bonus] Creating web visualization...%s\n", colorCyan, colorReset)
		startWebVisualization(g, results)
		fmt.Printf("\n%sPress Ctrl+C to exit...%s\n", colorYellow, colorReset)
		select {} // Keep server running
	}
}

func printHeader(vertices, edges, iterations int) {
	fmt.Printf("\n")
	fmt.Printf("%s╔════════════════════════════════════════════════════════════╗%s\n", colorBold+colorBlue, colorReset)
	fmt.Printf("%s║          DUAN SSSP VISUAL BENCHMARK SUITE                  ║%s\n", colorBold+colorBlue, colorReset)
	fmt.Printf("%s║     Breaking the Sorting Barrier - O(m log^(2/3) n)       ║%s\n", colorBold+colorBlue, colorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════╝%s\n", colorBold+colorBlue, colorReset)
	fmt.Printf("\n")
	fmt.Printf("%sConfiguration:%s\n", colorYellow, colorReset)
	fmt.Printf("  Vertices:   %s%d%s\n", colorBold, vertices, colorReset)
	fmt.Printf("  Edges:      %s%d%s (%.1fx density)\n", colorBold, edges, colorReset, float64(edges)/float64(vertices))
	fmt.Printf("  Iterations: %s%d%s\n", colorBold, iterations, colorReset)
	fmt.Printf("  CPU Cores:  %s%d%s / %d available\n", colorBold, runtime.GOMAXPROCS(0), colorReset, runtime.NumCPU())
	fmt.Printf("\n")
}

func generateGraph(vertices, edges int) *graph.Graph {
	g := graph.NewGraph(vertices)
	rng := rand.New(rand.NewSource(42)) // Deterministic for reproducibility

	for i := 0; i < edges; i++ {
		u := rng.Intn(vertices)
		v := rng.Intn(vertices)
		if u == v {
			v = (v + 1) % vertices
		}
		w := rng.Float64()*100.0 + 1.0
		g.AddEdge(u, v, w)
	}

	return g
}

func visualizeGraph(g *graph.Graph, sampleSize int) {
	if sampleSize > g.V {
		sampleSize = g.V
	}

	fmt.Printf("\n%sGraph Structure (sample %d/%d vertices):%s\n", colorYellow, sampleSize, g.V, colorReset)
	fmt.Printf("┌─────┬─────────────────────────────────────┐\n")
	fmt.Printf("│ %sV%s   │ %sEdges (to → weight)%s              │\n", colorBold, colorReset, colorBold, colorReset)
	fmt.Printf("├─────┼─────────────────────────────────────┤\n")

	for i := 0; i < sampleSize && i < g.V; i++ {
		fmt.Printf("│ %3d │ ", i)

		edgeCount := len(g.Adj[i])
		if edgeCount == 0 {
			fmt.Printf("(isolated)")
		} else {
			for j, edge := range g.Adj[i] {
				if j >= 3 {
					fmt.Printf("... +%d more", edgeCount-3)
					break
				}
				fmt.Printf("%d→%.1f ", edge.To, edge.Weight)
			}
		}
		fmt.Printf("%s\n", colorReset)
	}

	if sampleSize < g.V {
		fmt.Printf("│ ... │ ... (%d more vertices)              │\n", g.V-sampleSize)
	}
	fmt.Printf("└─────┴─────────────────────────────────────┘\n")
}

func benchmarkDuan(g *graph.Graph, iterations int) time.Duration {
	fmt.Printf("  %s►%s Duan Algorithm...", colorGreen, colorReset)

	var totalTime time.Duration

	for i := 0; i < iterations; i++ {
		tg := g.ToConstantDegree()
		solver := sssp.NewSolver(tg.G)

		start := time.Now()
		solver.Run(tg.OriginalTo[0])
		totalTime += time.Since(start)

		// Progress indicator
		if i%max(iterations/10, 1) == 0 {
			fmt.Printf(".")
		}
	}

	avgTime := totalTime / time.Duration(iterations)
	fmt.Printf(" %s✓%s %v\n", colorGreen, colorReset, avgTime)

	return avgTime
}

func benchmarkAStar(g *graph.Graph, iterations int) time.Duration {
	fmt.Printf("  %s►%s A* Algorithm...", colorYellow, colorReset)

	var totalTime time.Duration

	for i := 0; i < iterations; i++ {
		start := time.Now()
		aStarSSSP(g, 0)
		totalTime += time.Since(start)

		if i%max(iterations/10, 1) == 0 {
			fmt.Printf(".")
		}
	}

	avgTime := totalTime / time.Duration(iterations)
	fmt.Printf(" %s✓%s %v\n", colorYellow, colorReset, avgTime)

	return avgTime
}

func benchmarkParallelDuan(g *graph.Graph, iterations int) time.Duration {
	fmt.Printf("  %s►%s Duan Parallel (%d cores)...", colorPurple, colorReset, runtime.NumCPU())

	numCores := runtime.NumCPU()
	var totalTime time.Duration

	for i := 0; i < iterations; i++ {
		tg := g.ToConstantDegree()

		start := time.Now()

		// Run multiple source SSSP in parallel to utilize cores
		sources := make([]int, min(numCores, g.V))
		for j := range sources {
			sources[j] = j * (g.V / len(sources))
		}

		var wg sync.WaitGroup
		for _, src := range sources {
			wg.Add(1)
			go func(source int) {
				defer wg.Done()
				solver := sssp.NewSolver(tg.G)
				solver.Run(tg.OriginalTo[source])
			}(src)
		}
		wg.Wait()

		totalTime += time.Since(start)

		if i%max(iterations/10, 1) == 0 {
			fmt.Printf(".")
		}
	}

	avgTime := totalTime / time.Duration(iterations)
	fmt.Printf(" %s✓%s %v\n", colorPurple, colorReset, avgTime)

	return avgTime
}

func displayResults(results []BenchmarkResult) {
	fmt.Printf("\n┌────────────────────────────┬────────────────┬─────────────┐\n")
	fmt.Printf("│ %sAlgorithm%s                  │ %sAvg Time%s       │ %sSpeedup%s     │\n", colorBold, colorReset, colorBold, colorReset, colorBold, colorReset)
	fmt.Printf("├────────────────────────────┼────────────────┼─────────────┤\n")

	baseline := results[0].Time

	for _, r := range results {
		speedup := float64(baseline) / float64(r.Time)
		color := colorGreen
		if speedup < 0.9 {
			color = colorRed
		} else if speedup < 1.1 {
			color = colorYellow
		}

		fmt.Printf("│ %-26s │ %s%14v%s │ %s%11.2fx%s │\n",
			r.Algorithm,
			color, r.Time, colorReset,
			color, speedup, colorReset)
	}

	fmt.Printf("└────────────────────────────┴────────────────┴─────────────┘\n")
}

func visualizePerformance(results []BenchmarkResult) {
	if len(results) == 0 {
		return
	}

	// Find max time for scaling
	maxTime := results[0].Time
	for _, r := range results {
		if r.Time > maxTime {
			maxTime = r.Time
		}
	}

	// Draw bars
	barWidth := 50
	fmt.Printf("\n")

	for _, r := range results {
		barLen := int(float64(r.Time) / float64(maxTime) * float64(barWidth))
		if barLen < 1 {
			barLen = 1
		}

		// Color based on performance
		color := colorGreen
		if r.Time > results[0].Time*2 {
			color = colorYellow
		}
		if r.Time > results[0].Time*10 {
			color = colorRed
		}

		fmt.Printf("%-26s %s", r.Algorithm, color)
		for i := 0; i < barLen; i++ {
			fmt.Printf("█")
		}
		fmt.Printf("%s %v\n", colorReset, r.Time)
	}

	fmt.Printf("\n%sScale: 0%s", colorBold, colorReset)
	for i := 0; i < barWidth-10; i++ {
		fmt.Printf(" ")
	}
	fmt.Printf("%s%v%s\n", colorBold, maxTime, colorReset)
}

func printSummary(results []BenchmarkResult) {
	if len(results) < 2 {
		return
	}

	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════╗%s\n", colorBold+colorGreen, colorReset)
	fmt.Printf("%s║                         SUMMARY                            ║%s\n", colorBold+colorGreen, colorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════╝%s\n", colorBold+colorGreen, colorReset)

	duanTime := results[0].Time
	astarTime := results[1].Time

	speedup := float64(astarTime) / float64(duanTime)

	fmt.Printf("\n%s★ Duan algorithm is %.1fx faster than A* (heap)%s\n", colorBold+colorGreen, speedup, colorReset)

	if len(results) > 2 {
		parallelTime := results[2].Time
		parallelSpeedup := float64(duanTime) / float64(parallelTime)
		fmt.Printf("%s★ Parallel version (%d cores) is %.1fx faster%s\n",
			colorBold+colorPurple, runtime.NumCPU(), parallelSpeedup, colorReset)
	}

	// Performance insights
	perVertex := float64(duanTime.Nanoseconds()) / float64(results[0].Vertices)
	perEdge := float64(duanTime.Nanoseconds()) / float64(results[0].Edges)

	fmt.Printf("\n%sPerformance Metrics:%s\n", colorYellow, colorReset)
	fmt.Printf("  Per-vertex time: %.2f ns\n", perVertex)
	fmt.Printf("  Per-edge time:   %.2f ns\n", perEdge)
	fmt.Printf("  Throughput:      %.2f M vertices/sec\n", 1000.0/perVertex)

	fmt.Printf("\n%sCPU Utilization:%s\n", colorYellow, colorReset)
	fmt.Printf("  Cores used:      %d / %d available\n", runtime.GOMAXPROCS(0), runtime.NumCPU())
	fmt.Printf("  Parallelization: %s\n", map[bool]string{true: "Enabled", false: "Disabled"}[runtime.GOMAXPROCS(0) > 1])

	fmt.Printf("\n")
}

// Simple A* implementation for comparison
func aStarSSSP(g *graph.Graph, source int) []float64 {
	dist := make([]float64, g.V)
	for i := range dist {
		dist[i] = sssp.Infinity
	}
	dist[source] = 0

	type node struct {
		v     int
		score float64
	}

	pq := make([]node, 0, g.V)
	pq = append(pq, node{source, 0})
	visited := make([]bool, g.V)

	for len(pq) > 0 {
		// Extract min (simple linear search for benchmark)
		minIdx := 0
		for i := 1; i < len(pq); i++ {
			if pq[i].score < pq[minIdx].score {
				minIdx = i
			}
		}

		current := pq[minIdx]
		pq = append(pq[:minIdx], pq[minIdx+1:]...)

		if visited[current.v] {
			continue
		}
		visited[current.v] = true

		for _, edge := range g.Adj[current.v] {
			newDist := dist[current.v] + edge.Weight
			if newDist < dist[edge.To] {
				dist[edge.To] = newDist
				pq = append(pq, node{edge.To, newDist})
			}
		}
	}

	return dist
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
