# Quick Start Guide

Get started with the Duan et al. (2025) SSSP implementation in 5 minutes.

## Installation

```bash
# Clone the repository
git clone https://github.com/phr3nzy/duan-sssp.git
cd duan-sssp

# Or install as a module
go get github.com/phr3nzy/duan-sssp
```

## Basic Usage

### 1. Create a Simple Graph

```go
package main

import (
    "fmt"
    "github.com/phr3nzy/duan-sssp/graph"
    "github.com/phr3nzy/duan-sssp/sssp"
)

func main() {
    // Create a graph with 5 vertices
    g := graph.NewGraph(5)
    
    // Add edges: AddEdge(from, to, weight)
    g.AddEdge(0, 1, 10.0)
    g.AddEdge(0, 2, 5.0)
    g.AddEdge(1, 2, 2.0)
    g.AddEdge(1, 3, 1.0)
    g.AddEdge(2, 3, 9.0)
    g.AddEdge(2, 4, 2.0)
    g.AddEdge(3, 4, 4.0)
    
    // Transform to constant-degree graph
    tg := g.ToConstantDegree()
    
    // Solve SSSP from vertex 0
    solver := sssp.NewSolver(tg.G)
    rawDist := solver.Run(tg.OriginalTo[0])
    
    // Map distances back to original graph
    distances := tg.MapDistances(rawDist)
    
    // Print results
    fmt.Println("Shortest distances from vertex 0:")
    for i, d := range distances {
        fmt.Printf("  to %d: %.2f\n", i, d)
    }
}
```

**Output**:
```
Shortest distances from vertex 0:
  to 0: 0.00
  to 1: 7.00
  to 2: 5.00
  to 3: 8.00
  to 4: 7.00
```

### 2. Generate Random Graph

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
    "github.com/phr3nzy/duan-sssp/graph"
    "github.com/phr3nzy/duan-sssp/sssp"
)

func main() {
    // Parameters
    numVertices := 1000
    numEdges := 3000  // Sparse graph: m = 3n
    
    // Create random graph
    g := graph.NewGraph(numVertices)
    rand.Seed(time.Now().UnixNano())
    
    for i := 0; i < numEdges; i++ {
        from := rand.Intn(numVertices)
        to := rand.Intn(numVertices)
        weight := rand.Float64() * 100.0
        g.AddEdge(from, to, weight)
    }
    
    // Solve SSSP
    start := time.Now()
    tg := g.ToConstantDegree()
    transformTime := time.Since(start)
    
    solver := sssp.NewSolver(tg.G)
    start = time.Now()
    rawDist := solver.Run(tg.OriginalTo[0])
    solveTime := time.Since(start)
    
    distances := tg.MapDistances(rawDist)
    
    // Statistics
    reachable := 0
    for _, d := range distances {
        if d < sssp.Infinity {
            reachable++
        }
    }
    
    fmt.Printf("Graph: %d vertices, %d edges\n", numVertices, numEdges)
    fmt.Printf("Transform time: %v\n", transformTime)
    fmt.Printf("Solve time: %v\n", solveTime)
    fmt.Printf("Reachable vertices: %d/%d\n", reachable, numVertices)
}
```

### 3. Read Graph from File

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
    "github.com/phr3nzy/duan-sssp/graph"
    "github.com/phr3nzy/duan-sssp/sssp"
)

func loadGraph(filename string) (*graph.Graph, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    
    // First line: number of vertices
    scanner.Scan()
    numVertices, _ := strconv.Atoi(scanner.Text())
    g := graph.NewGraph(numVertices)
    
    // Remaining lines: from to weight
    for scanner.Scan() {
        parts := strings.Fields(scanner.Text())
        if len(parts) != 3 {
            continue
        }
        
        from, _ := strconv.Atoi(parts[0])
        to, _ := strconv.Atoi(parts[1])
        weight, _ := strconv.ParseFloat(parts[2], 64)
        
        g.AddEdge(from, to, weight)
    }
    
    return g, scanner.Err()
}

func main() {
    // Load graph from file
    g, err := loadGraph("graph.txt")
    if err != nil {
        fmt.Printf("Error loading graph: %v\n", err)
        return
    }
    
    // Solve SSSP
    tg := g.ToConstantDegree()
    solver := sssp.NewSolver(tg.G)
    rawDist := solver.Run(tg.OriginalTo[0])
    distances := tg.MapDistances(rawDist)
    
    // Print results
    for i, d := range distances {
        if d < sssp.Infinity {
            fmt.Printf("%d: %.2f\n", i, d)
        }
    }
}
```

**File format** (`graph.txt`):
```
5
0 1 10.0
0 2 5.0
1 2 2.0
1 3 1.0
2 3 9.0
2 4 2.0
3 4 4.0
```

## Running the Example

```bash
# Run the main example
go run main.go

# Expected output:
# Initializing High-Performance SSSP (Duan et al., 2025)...
# Generating graph V=10000, E=30000...
# Transforming to Constant Degree Graph...
# Transformation done in 3.361613ms. New V=60024
# Running BMSSP...
# Execution Time: 28.638µs
# Distance to node 10: 179.23
# Done.
```

## Running Tests

```bash
# Run basic tests
go test ./sssp/

# Run with verbose output
go test -v ./sssp/

# Run specific test
go test -run=TestBasicExecution ./sssp/
```

## Running Benchmarks

```bash
# Quick benchmark
go test -bench=BenchmarkSSSP/Small -benchtime=10x ./sssp/

# Full benchmark suite
go test -bench=. -benchmem ./sssp/

# Save results
go test -bench=. -benchmem ./sssp/ | tee benchmark_results.txt
```

## Common Patterns

### Multiple Sources

```go
sources := []int{0, 5, 10, 15}
results := make(map[int][]float64)

// Transform once
tg := g.ToConstantDegree()

// Run SSSP for each source
for _, src := range sources {
    solver := sssp.NewSolver(tg.G)
    rawDist := solver.Run(tg.OriginalTo[src])
    results[src] = tg.MapDistances(rawDist)
}
```

### Path Reconstruction

```go
// Note: Current implementation only computes distances
// To reconstruct paths, modify solver to track predecessors

type PathSolver struct {
    *sssp.Solver
    Pred []int  // Predecessor array
}

// Reconstruct path from source to target
func (ps *PathSolver) GetPath(target int) []int {
    path := []int{}
    for v := target; v != -1; v = ps.Pred[v] {
        path = append([]int{v}, path...)
    }
    return path
}
```

### Incremental Updates

```go
// For dynamic graphs with edge additions/deletions
// Current implementation doesn't support incremental updates
// Recommended: Recompute SSSP for each update

func updateGraph(g *graph.Graph, from, to int, newWeight float64) {
    // Remove old edge (if exists) - requires adjacency list modification
    // Add new edge
    g.AddEdge(from, to, newWeight)
    
    // Recompute
    tg := g.ToConstantDegree()
    solver := sssp.NewSolver(tg.G)
    distances := solver.Run(tg.OriginalTo[0])
}
```

## Performance Tips

### 1. Reuse Transformed Graph

```go
// Bad: Transform every time
for i := 0; i < numQueries; i++ {
    tg := g.ToConstantDegree()  // ❌ Expensive
    solver := sssp.NewSolver(tg.G)
    // ...
}

// Good: Transform once
tg := g.ToConstantDegree()  // ✅ Once
for i := 0; i < numQueries; i++ {
    solver := sssp.NewSolver(tg.G)
    // ...
}
```

### 2. Batch Processing

```go
// Process multiple queries in batch
sources := getAllSources()
results := make([][]float64, len(sources))

tg := g.ToConstantDegree()
for i, src := range sources {
    solver := sssp.NewSolver(tg.G)
    rawDist := solver.Run(tg.OriginalTo[src])
    results[i] = tg.MapDistances(rawDist)
}
```

### 3. Memory Management

```go
// For very large graphs, be mindful of memory
import "runtime"

func processLargeGraph(g *graph.Graph) {
    tg := g.ToConstantDegree()
    solver := sssp.NewSolver(tg.G)
    distances := solver.Run(tg.OriginalTo[0])
    
    // Process results immediately
    processDistances(distances)
    
    // Release memory
    tg = nil
    solver = nil
    runtime.GC()
}
```

## Troubleshooting

### Issue: Many vertices unreachable

**Symptom**: Most distances are Infinity

**Possible causes**:
1. Graph is actually disconnected
2. Implementation bug (known issue)
3. Source vertex has no outgoing edges

**Solution**:
- Verify graph connectivity with DFS/BFS
- Use traditional Dijkstra to validate
- Check graph construction code

### Issue: Slow performance

**Symptom**: Algorithm slower than expected

**Possible causes**:
1. Graph is too small (n < 1000)
2. Graph is very dense (m ≈ n²)
3. Debug mode enabled

**Solution**:
- Use Dijkstra for small graphs
- Compile with optimizations: `go build -ldflags="-s -w"`
- Profile with `go test -bench=. -cpuprofile=cpu.prof`

### Issue: High memory usage

**Symptom**: Out of memory errors

**Possible causes**:
1. Transformed graph doubles vertex count
2. Multiple solver instances
3. Large distance arrays

**Solution**:
- Process sources sequentially, not in parallel
- Release unused solvers
- Consider external memory algorithms for huge graphs

## Next Steps

- Read [README.md](README.md) for detailed algorithm explanation
- Check [ALGORITHM.md](ALGORITHM.md) for implementation details
- Review [BENCHMARKS.md](BENCHMARKS.md) for performance analysis
- Explore the source code in `sssp/`, `graph/`, and `ds/` directories

## Getting Help

- Open an issue on GitHub
- Check existing issues for similar problems
- Review the paper: [arXiv:2504.17033](https://arxiv.org/abs/2504.17033)

---

**Happy shortest path finding! 🚀**
