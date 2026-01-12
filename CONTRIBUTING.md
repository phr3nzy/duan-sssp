# Contributing to duan-sssp

Thank you for your interest in contributing to this implementation of the Duan et al. (2025) SSSP algorithm!

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/duan-sssp.git
   cd duan-sssp
   ```

3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/phr3nzy/duan-sssp.git
   ```

4. **Install dependencies**:
   ```bash
   go mod download
   ```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

Use descriptive branch names:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation improvements
- `perf/` - Performance improvements
- `test/` - Test additions/improvements

### 2. Make Your Changes

- Write clear, commented code
- Follow Go conventions and idioms
- Update documentation as needed
- Add tests for new functionality

### 3. Test Your Changes

```bash
# Run tests
go test ./...

# Run tests with race detector
go test -race ./...

# Run benchmarks
go test -bench=. ./sssp/

# Check test coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 4. Format and Lint

```bash
# Format code
go fmt ./...

# Run linter (if installed)
golangci-lint run
```

### 5. Commit Your Changes

Write clear, descriptive commit messages:

```bash
git add .
git commit -m "feat: add support for weighted graphs"
```

Commit message format:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `perf:` - Performance improvements
- `test:` - Test changes
- `refactor:` - Code refactoring
- `chore:` - Maintenance tasks

### 6. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## Pull Request Guidelines

### PR Description

Include in your PR description:

1. **What** - What does this PR do?
2. **Why** - Why is this change necessary?
3. **How** - How does it work?
4. **Testing** - How was it tested?
5. **Benchmarks** - Performance impact (if applicable)

Example:
```markdown
## What
Fixes stack overflow in FindPivots function

## Why
The tree size calculation could recurse infinitely on graphs with cycles

## How
Added cycle detection using three-state memoization

## Testing
- Added test cases for cyclic graphs
- All existing tests pass
- No performance regression

## Benchmarks
No significant performance impact (< 2% overhead)
```

### Checklist

Before submitting, ensure:

- [ ] Code compiles without errors
- [ ] All tests pass
- [ ] New tests added for new functionality
- [ ] Documentation updated
- [ ] Benchmarks run successfully
- [ ] No linter warnings
- [ ] Commit messages are clear

## Areas for Contribution

### High Priority

1. **Fix reachability bug** - Not all vertices are discovered in some graphs
2. **Performance optimization** - Reduce constant factors
3. **Memory optimization** - Reduce transformation overhead
4. **Test coverage** - Add more edge case tests

### Medium Priority

5. **Parallel implementation** - Parallelize FindPivots
6. **Additional graph formats** - Edge list, matrix support
7. **Visualization tools** - Graph and performance visualization
8. **Incremental SSSP** - Support for dynamic graphs

### Documentation

9. **Tutorial blog posts** - Explain the algorithm
10. **Video walkthrough** - Visual explanation
11. **Jupyter notebooks** - Interactive examples
12. **API documentation** - Godoc improvements

## Code Style

### Go Conventions

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Use meaningful variable names
- Comment exported functions
- Keep functions focused (< 50 lines when possible)

### Algorithm-Specific

- Maintain O(m log^(2/3) n) complexity
- Document complexity in comments
- Explain non-obvious optimizations
- Reference paper sections when relevant

### Example

```go
// FindPivots identifies vertices with large shortest path trees.
// It implements Algorithm 1 from Section 3 of the paper.
// Time complexity: O(k·m') where m' is edges in subgraph.
func (s *Solver) FindPivots(B float64, S []int) ([]int, []int) {
    // Implementation with clear comments
    // explaining each step
}
```

## Testing Guidelines

### Unit Tests

- Test happy paths and edge cases
- Use table-driven tests for multiple scenarios
- Test error conditions
- Keep tests fast (< 1s per test)

### Benchmark Tests

- Use consistent input sizes
- Run multiple iterations (-benchtime)
- Compare against baseline
- Document performance characteristics

### Example Test

```go
func TestFindPivots(t *testing.T) {
    tests := []struct {
        name     string
        graph    *graph.Graph
        bound    float64
        expected int  // expected pivot count
    }{
        {
            name:     "Simple graph",
            graph:    makeSimpleGraph(),
            bound:    100.0,
            expected: 2,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            solver := NewSolver(tt.graph)
            pivots, _ := solver.FindPivots(tt.bound, []int{0})
            
            if len(pivots) != tt.expected {
                t.Errorf("got %d pivots, want %d", len(pivots), tt.expected)
            }
        })
    }
}
```

## Performance Considerations

When optimizing:

1. **Measure first** - Profile before optimizing
2. **Maintain correctness** - Tests must pass
3. **Document trade-offs** - Explain complexity vs simplicity
4. **Benchmark impact** - Show before/after numbers

## Documentation Standards

### Code Comments

- Explain **why**, not **what**
- Reference paper algorithms/lemmas
- Document complexity
- Describe assumptions

### README Updates

- Update features list
- Add examples for new functionality
- Update benchmark results
- Document breaking changes

### Algorithm Documentation

- Update ALGORITHM.md for implementation changes
- Update BENCHMARKS.md for performance changes
- Keep QUICKSTART.md beginner-friendly

## Questions?

- Open an issue for questions
- Check existing issues and PRs
- Review ALGORITHM.md for implementation details
- Read the original paper: [arXiv:2504.17033](https://arxiv.org/abs/2504.17033)

## Code of Conduct

- Be respectful and constructive
- Welcome newcomers
- Focus on the code, not the person
- Give credit where due
- Help others learn

## Recognition

Contributors will be acknowledged in:
- README.md Contributors section
- CHANGELOG.md for significant changes
- Release notes

Thank you for contributing to advancing SSSP research! 🚀
