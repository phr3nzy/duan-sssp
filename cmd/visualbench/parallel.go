package main

import (
	"runtime"
	"sync"
	"time"

	"github.com/phr3nzy/duan-sssp/graph"
	"github.com/phr3nzy/duan-sssp/sssp"
)

// ParallelBenchmarkDuan runs Duan algorithm with multiple sources in parallel
func benchmarkParallelMultiSource(g *graph.Graph, iterations int) time.Duration {
	numCores := runtime.NumCPU()
	var totalTime time.Duration

	for iter := 0; iter < iterations; iter++ {
		tg := g.ToConstantDegree()

		// Select sources spread across the graph
		sources := make([]int, min(numCores*2, g.V))
		for i := range sources {
			sources[i] = (i * g.V) / len(sources)
		}

		start := time.Now()

		// Run SSSP from multiple sources in parallel
		var wg sync.WaitGroup
		results := make([][]float64, len(sources))

		for idx, src := range sources {
			wg.Add(1)
			go func(i, source int) {
				defer wg.Done()
				solver := sssp.NewSolver(tg.G)
				rawDist := solver.Run(tg.OriginalTo[source])
				results[i] = tg.MapDistances(rawDist)
			}(idx, src)
		}

		wg.Wait()
		totalTime += time.Since(start)
	}

	return totalTime / time.Duration(iterations)
}
