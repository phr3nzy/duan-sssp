package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/phr3nzy/duan-sssp/graph"
	"github.com/phr3nzy/duan-sssp/sssp"
)

func main() {
	fmt.Println("Initializing High-Performance SSSP (Duan et al., 2025)...")

	// 1. Generate a Sparse Random Graph
	V := 10000
	E := V * 3
	fmt.Printf("Generating graph V=%d, E=%d...\n", V, E)

	g := graph.NewGraph(V)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < E; i++ {
		u := rand.Intn(V)
		v := rand.Intn(V)
		w := rand.Float64() * 100.0
		g.AddEdge(u, v, w)
	}

	// 2. Transform (Critical Step)
	fmt.Println("Transforming to Constant Degree Graph...")
	startT := time.Now()
	tg := g.ToConstantDegree()
	fmt.Printf("Transformation done in %v. New V=%d\n", time.Since(startT), tg.G.V)

	// 3. Run Algorithm
	fmt.Println("Running BMSSP...")
	solver := sssp.NewSolver(tg.G)
	start := time.Now()
	rawDist := solver.Run(tg.OriginalTo[0]) // Run from mapped source 0
	duration := time.Since(start)

	fmt.Printf("Execution Time: %v\n", duration)

	// 4. Verification (Spot Check)
	mapped := tg.MapDistances(rawDist)
	fmt.Printf("Distance to node 10: %f\n", mapped[10])
	fmt.Println("Done.")
}
