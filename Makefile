.PHONY: all build test bench clean visualbench install

# Build configuration
BINARY_NAME=duan-sssp
VISUALBENCH_NAME=visualbench
GO=go
GOFLAGS=-v

all: build test

build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

visualbench:
	$(GO) build $(GOFLAGS) -o $(VISUALBENCH_NAME) ./cmd/visualbench

install:
	$(GO) install ./...
	$(GO) install ./cmd/visualbench

test:
	$(GO) test -v ./...

bench:
	$(GO) test -bench=. -benchmem ./sssp/

bench-comparison:
	$(GO) test -bench=BenchmarkComparison -benchtime=10x ./sssp/

bench-all:
	$(GO) test -bench=. -benchtime=10x -benchmem ./sssp/

# Run visual benchmark with all cores
visual: visualbench
	./$(VISUALBENCH_NAME) \
		-vertices=5000 \
		-edge-factor=3 \
		-iterations=10 \
		-parallel=true \
		-show-graph=true

# Run visual benchmark with web interface
visual-web: visualbench
	./$(VISUALBENCH_NAME) \
		-vertices=2000 \
		-edge-factor=3 \
		-iterations=5 \
		-parallel=true \
		-web=true

# Run large benchmark (uses all cores)
visual-large: visualbench
	./$(VISUALBENCH_NAME) \
		-vertices=50000 \
		-edge-factor=3 \
		-iterations=5 \
		-parallel=true \
		-show-graph=false

# Run with custom parameters
visual-custom: visualbench
	@echo "Usage: make visual-custom VERTICES=10000 EDGES=30000 ITER=10"
	./$(VISUALBENCH_NAME) \
		-vertices=$(or $(VERTICES),10000) \
		-edge-factor=$(or $(EDGEFACTOR),3) \
		-iterations=$(or $(ITER),10) \
		-parallel=true

# Profile CPU
profile-cpu:
	$(GO) test -bench=BenchmarkSSSP -cpuprofile=cpu.prof ./sssp/
	$(GO) tool pprof -http=:8081 cpu.prof

# Profile memory
profile-mem:
	$(GO) test -bench=BenchmarkSSSP -memprofile=mem.prof -benchmem ./sssp/
	$(GO) tool pprof -http=:8081 mem.prof

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME) $(VISUALBENCH_NAME)
	rm -f cpu.prof mem.prof
	rm -f benchmark_viz.html
	rm -f *.test
	$(GO) clean

# Format code
fmt:
	$(GO) fmt ./...

# Lint
lint:
	golangci-lint run

# Full CI simulation
ci: fmt lint test bench

help:
	@echo "Duan SSSP Makefile Commands:"
	@echo ""
	@echo "  make build          - Build main binary"
	@echo "  make test           - Run tests"
	@echo "  make bench          - Run benchmarks"
	@echo "  make visual         - Run visual benchmark (terminal)"
	@echo "  make visual-web     - Run visual benchmark (browser)"
	@echo "  make visual-large   - Run large-scale benchmark"
	@echo "  make profile-cpu    - CPU profiling"
	@echo "  make profile-mem    - Memory profiling"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make ci             - Run full CI suite"
	@echo ""
	@echo "Visual benchmark options:"
	@echo "  make visual VERTICES=10000 EDGEFACTOR=5 ITER=20"
