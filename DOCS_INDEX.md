# Documentation Index

Complete documentation for the Duan et al. (2025) SSSP implementation.

## 📚 Documentation Files

### 1. **README.md** - Main Documentation
**Purpose**: Primary entry point with comprehensive overview  
**Contents**:
- Algorithm overview and key achievement
- Installation instructions
- Basic and advanced usage examples
- Benchmark instructions
- Complexity analysis
- Comparison with other algorithms
- Known limitations

**Best for**: First-time users, general overview

---

### 2. **QUICKSTART.md** - Getting Started Guide
**Purpose**: 5-minute guide to get up and running  
**Contents**:
- Installation steps
- Three complete code examples
- Running examples and tests
- Common usage patterns
- Performance tips
- Troubleshooting guide

**Best for**: Developers who want to start coding immediately

---

### 3. **ALGORITHM.md** - Implementation Details
**Purpose**: Deep dive into the algorithm implementation  
**Contents**:
- Graph transformation walkthrough
- Data structure internals
- Detailed algorithm breakdowns (FindPivots, BaseCase, BMSSP)
- Complexity derivations
- Implementation decisions and rationale
- Edge cases and testing strategy
- Future optimization opportunities

**Best for**: Understanding how the code works, contributors, researchers

---

### 4. **BENCHMARKS.md** - Performance Analysis
**Purpose**: Comprehensive performance evaluation  
**Contents**:
- Complete benchmark results
- Scalability analysis
- Memory usage patterns
- Comparison with naive Dijkstra
- Theoretical vs actual performance
- Profiling insights
- Performance optimization tips

**Best for**: Performance analysis, optimization work, academic research

---

### 5. **SUMMARY.md** - Project Overview
**Purpose**: High-level project summary  
**Contents**:
- Project statistics
- Code structure overview
- Quick performance numbers
- Key components summary
- Usage examples
- Known issues
- Future work

**Best for**: Project managers, quick reference, citations

---

### 6. **DOCS_INDEX.md** - This File
**Purpose**: Navigation guide for all documentation  
**Contents**: You're reading it!

---

## 📂 Code Documentation

### Source Files

#### `graph/graph.go`
- **Graph**: Adjacency list representation
- **TransformedGraph**: Constant-degree transformation
- **Functions**: `NewGraph`, `AddEdge`, `ToConstantDegree`, `MapDistances`
- **Complexity**: O(m) transformation

#### `ds/ds.go`
- **DataStructure**: Block-based priority queue (Lemma 3.3)
- **Block**: Internal block structure
- **Item**: Key-value pairs
- **Functions**: `Insert`, `BatchPrepend`, `Pull`
- **Complexity**: O(log(N/M)) amortized insert

#### `sssp/sssp.go`
- **Solver**: Main SSSP algorithm state
- **Algorithms**: `BMSSP`, `FindPivots`, `BaseCase`
- **PriorityQueue**: Standard heap for base case
- **Functions**: `NewSolver`, `Run`
- **Complexity**: O(m log^(2/3) n)

#### `sssp/sssp_bench_test.go`
- **Benchmarks**: Comprehensive performance tests
- **Tests**: Basic execution validation
- **Utilities**: Random graph generation, naive Dijkstra
- **Suites**: Size scaling, density, components, comparison

#### `main.go`
- **Example**: Complete usage demonstration
- **Features**: Random graph generation, timing, verification
- **Output**: Performance metrics

---

## 🎯 Documentation by Use Case

### Use Case 1: I Want to Use This Library
**Start here**:
1. [QUICKSTART.md](QUICKSTART.md) - Get running in 5 minutes
2. [README.md](README.md) - Usage examples and API reference

### Use Case 2: I Want to Understand the Algorithm
**Read these**:
1. [README.md](README.md) - Algorithm overview
2. [ALGORITHM.md](ALGORITHM.md) - Detailed implementation
3. Original paper: [arXiv:2504.17033](https://arxiv.org/abs/2504.17033)

### Use Case 3: I Want to Benchmark/Optimize
**Start here**:
1. [BENCHMARKS.md](BENCHMARKS.md) - Performance analysis
2. `sssp/sssp_bench_test.go` - Benchmark code
3. [ALGORITHM.md](ALGORITHM.md) - Implementation decisions

### Use Case 4: I Want to Contribute
**Read these**:
1. [ALGORITHM.md](ALGORITHM.md) - Understanding the code
2. [BENCHMARKS.md](BENCHMARKS.md) - Testing and validation
3. [README.md](README.md) - Known limitations section

### Use Case 5: I Want Quick Reference
**Check**:
1. [SUMMARY.md](SUMMARY.md) - Project at a glance
2. [QUICKSTART.md](QUICKSTART.md) - Code examples

### Use Case 6: I'm Researching SSSP Algorithms
**Read**:
1. [README.md](README.md) - Complexity analysis
2. [BENCHMARKS.md](BENCHMARKS.md) - Empirical results
3. [ALGORITHM.md](ALGORITHM.md) - Implementation details
4. Original paper

---

## 📊 Documentation Statistics

| File | Lines | Words | Purpose |
|------|-------|-------|---------|
| README.md | ~450 | ~4,000 | Main documentation |
| QUICKSTART.md | ~400 | ~3,500 | Getting started |
| ALGORITHM.md | ~900 | ~7,500 | Implementation details |
| BENCHMARKS.md | ~600 | ~5,000 | Performance analysis |
| SUMMARY.md | ~300 | ~2,500 | Project overview |
| DOCS_INDEX.md | ~150 | ~1,200 | This file |
| **Total** | **~2,800** | **~23,700** | All documentation |

---

## 🔗 Quick Links

### Documentation
- [README.md](README.md) - Start here
- [QUICKSTART.md](QUICKSTART.md) - 5-minute guide
- [ALGORITHM.md](ALGORITHM.md) - Deep dive
- [BENCHMARKS.md](BENCHMARKS.md) - Performance
- [SUMMARY.md](SUMMARY.md) - Overview

### Code
- [graph/graph.go](graph/graph.go) - Graph structures
- [ds/ds.go](ds/ds.go) - Data structures
- [sssp/sssp.go](sssp/sssp.go) - Core algorithm
- [sssp/sssp_bench_test.go](sssp/sssp_bench_test.go) - Benchmarks
- [main.go](main.go) - Example usage

### External
- [Paper (arXiv)](https://arxiv.org/abs/2504.17033)
- [GitHub Repository](https://github.com/phr3nzy/duan-sssp)

---

## 📖 Reading Order

### For Beginners
1. README.md (Overview)
2. QUICKSTART.md (Hands-on)
3. Experiment with examples
4. SUMMARY.md (Reference)

### For Developers
1. QUICKSTART.md (Setup)
2. README.md (API)
3. ALGORITHM.md (Internals)
4. Source code exploration

### For Researchers
1. README.md (Context)
2. ALGORITHM.md (Implementation)
3. BENCHMARKS.md (Empirical)
4. Original paper
5. Source code analysis

### For Contributors
1. All documentation files
2. Source code review
3. Test suite examination
4. Issue tracker

---

## 🔍 Search Guide

### Find Information About...

**Installation**: README.md § Installation, QUICKSTART.md § Installation

**Usage Examples**: README.md § Usage, QUICKSTART.md § Basic Usage

**Algorithm Details**: ALGORITHM.md (entire file), README.md § Algorithm Overview

**Performance**: BENCHMARKS.md (entire file), SUMMARY.md § Benchmark Results

**Data Structures**: ALGORITHM.md § Data Structures, ds/ds.go (source)

**Graph Transformation**: ALGORITHM.md § Graph Transformation, graph/graph.go (source)

**Time Complexity**: README.md § Complexity Analysis, ALGORITHM.md § Complexity Analysis

**Known Issues**: README.md § Known Limitations, BENCHMARKS.md § Known Limitations

**Future Work**: SUMMARY.md § Future Work, ALGORITHM.md § Future Optimizations

**Testing**: BENCHMARKS.md § Running Benchmarks, sssp/sssp_bench_test.go (source)

**Contributing**: README.md § Contributing, SUMMARY.md § Contributing

**Citations**: README.md § References, SUMMARY.md § License & Citation

---

## 💡 Tips

1. **Use Ctrl+F** to search within documents
2. **Start with QUICKSTART.md** if you learn by doing
3. **Start with README.md** if you learn by reading
4. **Bookmark this page** for quick navigation
5. **Check SUMMARY.md** for quick stats
6. **Read ALGORITHM.md** before contributing code
7. **Consult BENCHMARKS.md** for performance questions

---

## 📝 Documentation Maintenance

### Adding New Documentation

When adding new documentation:
1. Add entry to this index
2. Update README.md if relevant
3. Cross-reference from related documents
4. Update SUMMARY.md statistics

### Documentation Standards

- Use Markdown formatting
- Include code examples
- Provide complexity analysis
- Add cross-references
- Keep examples runnable
- Update benchmarks

---

**Last Updated**: January 2026  
**Documentation Version**: 1.0  
**Maintainer**: phr3nzy
