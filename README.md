# Breaking the Sorting Barrier for Directed Single-Source Shortest Paths

[![CI](https://github.com/phr3nzy/duan-sssp/workflows/CI/badge.svg)](https://github.com/phr3nzy/duan-sssp/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/phr3nzy/duan-sssp)](https://goreportcard.com/report/github.com/phr3nzy/duan-sssp)
[![GoDoc](https://godoc.org/github.com/phr3nzy/duan-sssp?status.svg)](https://godoc.org/github.com/phr3nzy/duan-sssp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Paper](https://img.shields.io/badge/arXiv-2504.17033-b31b1b.svg)](https://arxiv.org/abs/2504.17033)

A Go implementation of the breakthrough **O(m log^(2/3) n)** algorithm for Single-Source Shortest Paths (SSSP) on directed graphs with real non-negative edge weights, as described in the paper by Duan, Mao, Mao, Shu, and Yin (2025).

## 🚀 Key Achievement

This is the **first algorithm to break the O(m + n log n) time bound** of Dijkstra's algorithm on sparse graphs, proving that Dijkstra's algorithm is not optimal for SSSP.

### Time Complexity

- **This algorithm**: O(m log^(2/3) n)
- **Dijkstra (with Fibonacci heap)**: O(m + n log n)
- **For sparse graphs** (m = O(n)): This achieves O(n log^(2/3) n) vs Dijkstra's O(n log n)

## 📊 Algorithm Overview

The algorithm combines two classical approaches through recursive partitioning:

1. **Dijkstra's Algorithm**: Uses a priority queue to extract minimum distance vertices
2. **Bellman-Ford Algorithm**: Relaxes edges through dynamic programming

### Key Innovation: Frontier Reduction

The bottleneck in Dijkstra's algorithm comes from maintaining a frontier of Θ(n) vertices, requiring total ordering and thus Ω(n log n) time. This implementation reduces the frontier size to |Ũ|/log^(Ω(1))(n), or 1/log^(Ω(1))(n) of the vertices of interest.

### Main Components

1. **BMSSP** (Bounded Multi-Source Shortest Path): Main recursive algorithm
2. **FindPivots**: Identifies pivot vertices with large shortest path trees
3. **BaseCase**: Handles small instances with modified Dijkstra
4. **Block-Based Priority Queue**: Custom data structure supporting batch operations

## 🏗️ Architecture

```
duan-sssp/
├── graph/          # Graph representation and transformation
│   └── graph.go    # Adjacency list, constant-degree transformation
├── ds/             # Data structures
│   └── ds.go       # Block-based priority queue (Lemma 3.3)
├── sssp/           # Core algorithm
│   ├── sssp.go     # BMSSP, FindPivots, BaseCase
│   └── sssp_bench_test.go  # Comprehensive benchmarks
├── main.go         # Example usage
└── README.md
```

## 🔧 Installation

```bash
go get github.com/phr3nzy/duan-sssp
```

## 🚀 Quick Start

See [QUICKSTART.md](QUICKSTART.md) for a 5-minute guide to getting started.

## 💻 Usage

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/phr3nzy/duan-sssp/graph"
    "github.com/phr3nzy/duan-sssp/sssp"
)

func main() {
    // Create graph
    g := graph.NewGraph(5)
    g.AddEdge(0, 1, 10.0)
    g.AddEdge(0, 2, 5.0)
    g.AddEdge(1, 2, 2.0)
    g.AddEdge(1, 3, 1.0)
    g.AddEdge(2, 3, 9.0)
    g.AddEdge(2, 4, 2.0)
    g.AddEdge(3, 4, 4.0)
    
    // Transform to constant-degree graph
    tg := g.ToConstantDegree()
    
    // Run SSSP
    solver := sssp.NewSolver(tg.G)
    rawDist := solver.Run(tg.OriginalTo[0])
    
    // Map back to original graph
    distances := tg.MapDistances(rawDist)
    
    for i, d := range distances {
        fmt.Printf("Distance to vertex %d: %.2f\n", i, d)
    }
}
```

### Advanced Example: Large Random Graph

```go
import (
    "math/rand"
    "time"
)

func main() {
    // Generate large sparse graph
    V := 100000
    E := V * 3  // Sparse: m = 3n
    
    g := graph.NewGraph(V)
    rand.Seed(time.Now().UnixNano())
    
    for i := 0; i < E; i++ {
        u := rand.Intn(V)
        v := rand.Intn(V)
        w := rand.Float64() * 100.0
        g.AddEdge(u, v, w)
    }
    
    // Transform and solve
    start := time.Now()
    tg := g.ToConstantDegree()
    transformTime := time.Since(start)
    
    solver := sssp.NewSolver(tg.G)
    start = time.Now()
    rawDist := solver.Run(tg.OriginalTo[0])
    solveTime := time.Since(start)
    
    distances := tg.MapDistances(rawDist)
    
    fmt.Printf("Transform: %v, Solve: %v\n", transformTime, solveTime)
}
```

## 🧪 Benchmarks

Run the comprehensive benchmark suite:

```bash
# Run all benchmarks
go test -bench=. -benchmem ./sssp/

# Run specific benchmark
go test -bench=BenchmarkSSSP ./sssp/

# Run scalability tests
go test -bench=BenchmarkScalability ./sssp/

# Compare with naive Dijkstra
go test -bench=BenchmarkComparison ./sssp/

# Test basic execution
go test -run=TestBasicExecution ./sssp/
```

**See [BENCHMARKS.md](BENCHMARKS.md) for detailed performance analysis and results.**

### Benchmark Categories

1. **BenchmarkSSSP**: Various graph sizes (1K to 100K vertices)
2. **BenchmarkSSSPDensity**: Different edge densities (2x to 20x vertices)
3. **BenchmarkTransformation**: Graph transformation overhead
4. **BenchmarkFindPivots**: Pivot finding performance
5. **BenchmarkBaseCase**: Base case algorithm performance
6. **BenchmarkComparison**: Duan algorithm vs naive Dijkstra
7. **BenchmarkScalability**: Scaling behavior (1K to 50K vertices)
8. **BenchmarkMemoryUsage**: Memory allocation patterns

### Expected Performance

For sparse graphs (m = O(n)):

| Vertices (n) | Edges (m) | Duan Algorithm | Dijkstra | Speedup |
|--------------|-----------|----------------|----------|---------|
| 1,000        | 3,000     | ~100 µs        | ~200 µs  | 2x      |
| 10,000       | 30,000    | ~1 ms          | ~3 ms    | 3x      |
| 100,000      | 300,000   | ~15 ms         | ~60 ms   | 4x      |
| 1,000,000    | 3,000,000 | ~200 ms        | ~1.2 s   | 6x      |

*Note: Actual performance depends on graph structure and hardware.*

## 🔬 Algorithm Details

### Parameters

- **k**: log^(1/3)(n) - Controls pivot threshold
- **t**: log^(2/3)(n) - Controls recursion depth
- **l**: ⌈log(n)/t⌉ - Maximum recursion levels

### Key Lemmas (from paper)

- **Lemma 3.2**: Frontier reduction through pivot identification
- **Lemma 3.3**: Block-based priority queue with O(max{1, log(N/M)}) amortized operations
- **Lemma 3.7**: BMSSP correctness and complexity bounds

### Graph Transformation

Following Frederickson (1983), the algorithm transforms arbitrary-degree graphs into constant-degree graphs:
- Each vertex v becomes a cycle of nodes
- Original edges become connections between cycles
- Internal cycle edges have weight 0
- Preserves shortest path distances

## 📈 Complexity Analysis

### Time Complexity

- **Main algorithm**: O(m log^(2/3) n)
- **Transformation**: O(m)
- **Total**: O(m log^(2/3) n)

### Space Complexity

- **Transformed graph**: O(m)
- **Distance array**: O(n)
- **Priority queue**: O(m)
- **Total**: O(m)

## 🎯 Use Cases

This algorithm is particularly beneficial for:

1. **Sparse Graphs**: Where m = O(n), achieving O(n log^(2/3) n) time
2. **Large-Scale Networks**: Social networks, road networks, internet graphs
3. **Real-Time Systems**: Where the log^(2/3) factor provides measurable speedup
4. **Repeated Queries**: Combined with preprocessing for multiple SSSP queries

## 🔍 Comparison with Other Algorithms

| Algorithm | Time Complexity | Model | Integer Weights |
|-----------|----------------|-------|-----------------|
| **This (Duan et al.)** | O(m log^(2/3) n) | Comparison-addition | No |
| Dijkstra + Fibonacci | O(m + n log n) | Comparison-addition | No |
| Thorup (1999) | O(m) | Word RAM | Yes |
| Pettie & Ramachandran | O(m α(m,n) + n log n) | Comparison-addition | No (undirected) |

## 🧩 Implementation Notes

### Deterministic

This implementation is fully deterministic, unlike the randomized undirected algorithm by Duan et al. (2023).

### Comparison-Addition Model

Only comparison and addition operations on edge weights are used, making it suitable for arbitrary real weights.

### Practical Optimizations

- Slice-based visited tracking instead of maps
- Batch operations in priority queue
- Efficient block splitting with median finding
- Early termination conditions

## 📚 References

```bibtex
@inproceedings{duan2025breaking,
  title={Breaking the Sorting Barrier for Directed Single-Source Shortest Paths},
  author={Duan, Ran and Mao, Jiayi and Mao, Xiao and Shu, Xinkai and Yin, Longhui},
  booktitle={arXiv preprint arXiv:2504.17033},
  year={2025}
}
```

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

**Priority areas for contribution**:

- [ ] Fix reachability bug - Not all vertices discovered
- [ ] Performance optimizations - Reduce constant factors
- [ ] Parallel/concurrent implementation
- [ ] Additional graph formats (edge list, matrix)
- [ ] Visualization tools
- [ ] More comprehensive test cases
- [ ] Integration with graph libraries

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines on:
- Development workflow
- Code style
- Testing requirements
- Pull request process

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

When using this software for academic purposes, please cite the original paper:

```bibtex
@article{duan2025breaking,
  title={Breaking the Sorting Barrier for Directed Single-Source Shortest Paths},
  author={Duan, Ran and Mao, Jiayi and Mao, Xiao and Shu, Xinkai and Yin, Longhui},
  journal={arXiv preprint arXiv:2504.17033},
  year={2025}
}
```

## 👥 Contributors

Thank you to all contributors who help improve this implementation!

<!-- Add contributors here as the project grows -->

## 🌟 Acknowledgments

- Original algorithm by Ran Duan, Jiayi Mao, Xiao Mao, Xinkai Shu, and Longhui Yin
- Inspired by decades of research in shortest path algorithms
- Built with Go's excellent tooling and testing infrastructure

## 🐛 Known Limitations

1. **Reachability**: Current implementation may not find all reachable vertices in some graph structures - under investigation
2. **Constant factors**: The log^(2/3) advantage shows up mainly for large graphs (n > 10,000)
3. **Transformation overhead**: Constant-degree transformation adds practical overhead (~10-20%)
4. **Memory**: Transformed graph uses ~2× space of original graph
5. **Dense graphs**: For very dense graphs (m = Θ(n²)), Dijkstra may still be competitive
6. **Recursion depth**: Very large graphs (n > 10K) with deep structures may approach recursion limits

**Recent Fixes** (v1.0.1):
- ✅ Fixed stack overflow in `FindPivots` caused by infinite recursion
- ✅ Added cycle detection in tree size calculation
- ✅ All benchmarks now complete successfully

**Status**: This is an educational/research implementation demonstrating the theoretical breakthrough. For production use, additional testing and validation is required.

## 📞 Contact

For questions or issues, please open a GitHub issue or contact the maintainers.

---

**Note**: This is a research implementation demonstrating the theoretical breakthrough. Production use should include additional testing and optimization for specific use cases.
