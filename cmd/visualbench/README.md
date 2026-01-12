# Visual Benchmark Tool 📊🚀

Interactive benchmark suite with visualizations for the Duan SSSP algorithm.

## Quick Start

### Build
```bash
cd /home/phr3nzy/go/src/github.com/phr3nzy/duan-sssp
go build -o visualbench ./cmd/visualbench
```

Or use Make:
```bash
make visualbench
```

## 🚀 THE COMMAND YOU ASKED FOR

### Use All Cores with Visualization:

```bash
./visualbench \
  -vertices=10000 \
  -edge-factor=3 \
  -iterations=10 \
  -parallel=true \
  -show-graph=true \
  -web=true
```

### Or use the Makefile shortcut:

```bash
make visual-web
```

This will:
- ✅ Use ALL your CPU cores (28 cores detected!)
- ✅ Show terminal visualization of the graph
- ✅ Run benchmarks (Duan vs A* vs Parallel)
- ✅ Open browser with interactive visualization
- ✅ Display performance bars and metrics
- ✅ Show algorithm running in real-time

## 📋 Command Options

```
-vertices=N         Number of vertices (default: 10000)
-edge-factor=N      Edges = vertices × N (default: 3)
-iterations=N       Benchmark iterations (default: 10)
-parallel=BOOL      Use all CPU cores (default: true)
-show-graph=BOOL    Show terminal graph viz (default: true)
-web=BOOL           Open web visualization (default: false)
```

## 🎯 Example Commands

### 1. Quick Benchmark (Terminal Only)
```bash
./visualbench -vertices=5000 -iterations=5
```

### 2. Large-Scale Parallel Benchmark
```bash
./visualbench -vertices=50000 -edge-factor=3 -iterations=10 -parallel=true
```

### 3. Web Visualization (Recommended!)
```bash
./visualbench -vertices=2000 -web=true
```

### 4. Maximum Performance Test
```bash
./visualbench \
  -vertices=100000 \
  -edge-factor=3 \
  -iterations=20 \
  -parallel=true \
  -show-graph=false
```

### 5. Using Makefile Shortcuts
```bash
# Quick visual benchmark
make visual

# Web visualization  
make visual-web

# Large scale test
make visual-large

# Custom parameters
make visual VERTICES=20000 EDGEFACTOR=5 ITER=15
```

## 📊 Output Example

```
╔════════════════════════════════════════════════════════════╗
║          DUAN SSSP VISUAL BENCHMARK SUITE                  ║
║     Breaking the Sorting Barrier - O(m log^(2/3) n)       ║
╚════════════════════════════════════════════════════════════╝

Configuration:
  Vertices:   10000
  Edges:      30000 (3.0x density)
  Iterations: 10
  CPU Cores:  28 / 28 available

[1/4] Generating random graph...

Graph Structure (sample 20/1000 vertices):
┌─────┬─────────────────────────────────────┐
│ V   │ Edges (to → weight)                 │
├─────┼─────────────────────────────────────┤
│   0 │ 304→32.7 756→91.9 ...              │
...

[2/4] Running benchmarks with 28 cores...
  ► Duan Algorithm.......... ✓ 226.5µs
  ► A* Algorithm.......... ✓ 1.95ms
  ► Duan Parallel (28 cores).. ✓ 45.2µs

[3/4] Results:
┌────────────────────────────┬────────────────┬─────────────┐
│ Algorithm                  │ Avg Time       │ Speedup     │
├────────────────────────────┼────────────────┼─────────────┤
│ Duan (O(m log^(2/3) n))   │       226.5µs  │       1.00x │
│ A* with Heap               │         1.95ms │       0.12x │
│ Duan Parallel (28 cores)   │        45.2µs  │       5.01x │
└────────────────────────────┴────────────────┴─────────────┘

[4/4] Performance Visualization:

Duan (O(m log^(2/3) n))   ████                              226.5µs
A* with Heap               ████████████████████████████████ 1.95ms
Duan Parallel (28 cores)   ██                                45.2µs

╔════════════════════════════════════════════════════════════╗
║                         SUMMARY                            ║
╚════════════════════════════════════════════════════════════╝

★ Duan algorithm is 8.6x faster than A* (heap)
★ Parallel version (28 cores) is 5.0x faster

Performance Metrics:
  Per-vertex time: 22.65 ns
  Per-edge time:   7.55 ns
  Throughput:      44.15 M vertices/sec

CPU Utilization:
  Cores used:      28 / 28 available
  Parallelization: Enabled
```

## 🌐 Web Visualization Features

When you add `-web=true`:

1. **Opens automatically** in your default browser
2. **Interactive graph visualization** - Canvas-based rendering
3. **Performance bars** - Animated comparison
4. **Real-time stats** - Graph metrics
5. **Responsive design** - Beautiful gradient UI

The web page includes:
- Circular graph layout
- Edge rendering
- Source vertex highlighting
- Animated performance bars
- Detailed statistics
- Winner announcement

## 🔧 Integration with Main Benchmarks

The visual tool uses the same algorithms as the test suite:

```bash
# Compare terminal vs test benchmarks
./visualbench -vertices=10000 -iterations=10
go test -bench=BenchmarkComparison -benchtime=10x ./sssp/
```

Should show similar results!

## 💡 Performance Tips

### For Best Results:

1. **Use at least 5 iterations** for stable averages
2. **Start with smaller graphs** (1K-5K vertices) for web viz
3. **Use larger graphs** (50K-100K) for performance testing
4. **Enable parallel** to see multi-core benefits
5. **Compare results** across different graph sizes

### Graph Size Recommendations:

| Vertices | Purpose | Web Viz | Iterations |
|----------|---------|---------|------------|
| 1,000 | Quick test | ✅ Yes | 10 |
| 5,000 | Development | ✅ Yes | 10 |
| 10,000 | Standard benchmark | ⚠️ Slow | 10 |
| 50,000 | Performance test | ❌ No | 5 |
| 100,000 | Stress test | ❌ No | 3 |

## 🎨 Customizing Visualization

Edit `cmd/visualbench/web_viz.go` to customize:
- Graph layout algorithm
- Color schemes
- Chart types
- Additional metrics

## 🐛 Troubleshooting

**Browser doesn't open automatically**:
- Manually go to `http://localhost:8080/benchmark_viz.html`
- Check if port 8080 is available

**Performance seems slow**:
- Reduce `-iterations` for faster runs
- Use smaller graphs for initial testing
- Check CPU frequency/governor settings

**Graph too large to visualize**:
- Use `-show-graph=false` for huge graphs
- Web viz samples only first 100 vertices

## 📚 See Also

- `PERFORMANCE_ROADMAP.md` - Optimization plans
- `BENCHMARKS.md` - Detailed performance analysis
- `Makefile` - Convenient shortcuts

---

**Enjoy visualizing the breakthrough algorithm!** 🎉
