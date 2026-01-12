# Algorithm Implementation Details

This document provides a detailed walkthrough of the implementation of the Duan et al. (2025) SSSP algorithm.

## Table of Contents

1. [Overview](#overview)
2. [Graph Transformation](#graph-transformation)
3. [Data Structures](#data-structures)
4. [Core Algorithms](#core-algorithms)
5. [Complexity Analysis](#complexity-analysis)
6. [Implementation Decisions](#implementation-decisions)

## Overview

The algorithm breaks Dijkstra's O(m + n log n) barrier by reducing the size of the "frontier" - the set of vertices being processed. The key insight is that instead of maintaining total order over all frontier vertices, we can work with a smaller set of "pivot" vertices.

### Main Idea

At any point during Dijkstra's algorithm, the priority queue maintains a frontier S such that if vertex u is "incomplete" (current distance estimate > true distance), the shortest s-u path must visit some complete vertex v ∈ S.

Traditional Dijkstra picks the closest vertex from S (requiring full ordering), but we can reduce |S| to |Ũ|/log^(Ω(1))(n) where Ũ is the set of vertices of interest.

## Graph Transformation

### Constant-Degree Transformation

**File**: `graph/graph.go`

The algorithm requires constant-degree graphs. We implement Frederickson's (1983) transformation:

```
Original vertex v with d_in incoming and d_out outgoing edges
    ↓
Cycle of k nodes (k = d_in + d_out, or 1 if isolated)
```

#### Implementation Details

```go
func (g *Graph) ToConstantDegree() *TransformedGraph
```

**Steps**:

1. **Calculate in-degrees**: O(m)
   ```go
   inDegree := make([]int, g.V)
   for u := 0; u < g.V; u++ {
       for _, e := range g.Adj[u] {
           inDegree[e.To]++
       }
   }
   ```

2. **Allocate node IDs**: Each vertex gets a contiguous range of IDs
   ```go
   starts[u] = currentID
   sizes[u] = len(g.Adj[u]) + inDegree[u]
   currentID += sizes[u]
   ```

3. **Build internal cycles**: Zero-weight edges within each cycle
   ```go
   for i := 0; i < sz; i++ {
       curr := start + i
       next := start + (i+1)%sz
       newG.AddEdge(curr, next, 0)
   }
   ```

4. **Map original edges**: Each edge (u,v,w) becomes (u_slot, v_slot, w)
   ```go
   uNode := starts[u] + slots[u]
   vNode := starts[v] + slots[v]
   newG.AddEdge(uNode, vNode, w)
   ```

**Complexity**: O(m + n)

**Properties**:
- New graph has O(m) vertices
- Maximum degree is 2 (one cycle edge + at most one real edge)
- Shortest paths preserved (internal weights are 0)

### Distance Mapping

To recover original distances:

```go
func (tg *TransformedGraph) MapDistances(dist []float64) []float64 {
    res := make([]float64, len(tg.OriginalTo))
    for i, startNode := range tg.OriginalTo {
        res[i] = dist[startNode]  // Distance to cycle start
    }
    return res
}
```

## Data Structures

### Block-Based Priority Queue

**File**: `ds/ds.go`

Implements Lemma 3.3 from the paper: a priority queue supporting:
- **Insert(key, value)**: Amortized O(max{1, log(N/M)})
- **BatchPrepend(items)**: O(|items|/M) amortized
- **Pull()**: Returns M smallest items in O(M)

#### Structure

```go
type DataStructure struct {
    M     int       // Block size = 2^((l-1)t)
    B     float64   // Upper bound
    Count int       // Total items
    d0    []*block  // Buffer for BatchPrepend
    d1    []*block  // Sorted sequence of blocks
}

type block struct {
    head       *Item   // Linked list of items
    tail       *Item
    size       int
    upperBound float64  // Max value in block
}
```

#### Insert Operation

```go
func (ds *DataStructure) Insert(key int, val float64)
```

1. **Binary search** for appropriate block in d1 by upperBound: O(log(N/M))
   ```go
   idx := sort.Search(len(ds.d1), func(i int) bool {
       return ds.d1[i].upperBound >= val
   })
   ```

2. **Append** to linked list: O(1)
   ```go
   item.next = targetBlock.head
   targetBlock.head = item
   targetBlock.size++
   ```

3. **Split** if size exceeds M: O(M log M) amortized
   ```go
   if targetBlock.size > ds.M {
       ds.split(idx)
   }
   ```

#### BatchPrepend Operation

```go
func (ds *DataStructure) BatchPrepend(items []Item)
```

For inserting items known to be smaller than current minimum:

1. **Sort items**: O(|items| log |items|)
2. **Chunk into blocks** of size M: O(|items|/M) blocks
3. **Prepend to d0**: O(1) per block

```go
for i := 0; i < len(items); i += ds.M {
    chunk := items[i:min(i+ds.M, len(items))]
    blk := createBlockFromSortedItems(chunk)
    ds.d0 = append([]*block{blk}, ds.d0...)
}
```

#### Pull Operation

```go
func (ds *DataStructure) Pull() ([]Item, float64)
```

Returns M smallest items:

1. **Drain d0 first** (contains prepended items): O(M)
2. **Drain d1 if needed**: O(M log M) to sort blocks
3. **Return next boundary** Bi: O(1)

### Priority Queue for Base Case

Standard binary heap for the BaseCase algorithm:

```go
type PriorityQueue []*PQItem

type PQItem struct {
    u        int      // Vertex
    priority float64  // Distance
    index    int      // Heap index
}
```

Uses Go's `container/heap` with standard operations: Push, Pop, Fix in O(log n).

## Core Algorithms

### Algorithm 1: FindPivots

**File**: `sssp/sssp.go:239`

Identifies "pivot" vertices - those with large shortest path trees.

```go
func (s *Solver) FindPivots(B float64, S []int) ([]int, []int)
```

**Parameters**:
- `B`: Upper bound on distances to consider
- `S`: Source set

**Returns**:
- `P`: Pivot set (vertices with tree size ≥ k)
- `W`: All vertices reachable in k steps from S

**Algorithm**:

1. **Relax k steps** from S (Bellman-Ford style):
   ```go
   for i := 1; i <= s.K; i++ {
       for each u in Wi_prev:
           for each edge (u,v):
               if s.Dist[u] + w < s.Dist[v]:
                   s.Dist[v] = s.Dist[u] + w
                   Add v to Wi
   ```

2. **Early return** if |W| > k·|S|:
   ```go
   if len(W_list) > s.K * len(S) {
       return S, W_list  // P = S
   }
   ```

3. **Build forest F** on W:
   - Edge (u,v) in F iff u,v ∈ W and d[v] = d[u] + weight(u,v)
   
4. **Compute subtree sizes** with memoization:
   ```go
   var calcSize func(u int) int
   calcSize = func(u int) int {
       if memoSize[u] != 0 {
           return memoSize[u]
       }
       count := 1
       for each edge (u,v) in F:
           count += calcSize(v)
       memoSize[u] = count
       return count
   }
   ```

5. **Select pivots**: vertices in S with subtree size ≥ k
   ```go
   for each u in S:
       if calcSize(u) >= s.K:
           P = append(P, u)
   ```

**Complexity**: O(k·(|W|+|E_W|)) where E_W are edges within W

**Intuition**: If |W| is small relative to k·|S|, then some vertices in S must have large trees (pigeonhole principle), becoming pivots. Otherwise, S itself is already a good frontier.

### Algorithm 2: BaseCase

**File**: `sssp/sssp.go:319`

Handles recursion base case (l=0) with limited Dijkstra exploration.

```go
func (s *Solver) BaseCase(B float64, S []int) (float64, []int)
```

**Parameters**:
- `B`: Upper bound
- `S`: Source set (treated as multi-source)

**Returns**:
- `B'`: New upper bound
- `U`: Set of processed vertices

**Algorithm**:

1. **Initialize multi-source** priority queue:
   ```go
   for each x in S:
       U0[x] = true
       heap.Push(pq, &PQItem{u: x, priority: s.Dist[x]})
   ```

2. **Limited Dijkstra** (at most k+1 vertices):
   ```go
   limit := s.K + 1
   while pq not empty and len(U0) < limit:
       u = heap.Pop(pq)
       for each edge (u,v):
           if s.Dist[u] + w < s.Dist[v] and s.Dist[u] + w < B:
               s.Dist[v] = s.Dist[u] + w
               heap.Push(pq, v)
   ```

3. **Determine return bound**:
   - If |U0| ≤ k: return B, U0 (success)
   - Otherwise: return max{d[u] : u ∈ U0}, {u ∈ U0 : d[u] < max}

**Complexity**: O(k log k + |edges touched|)

**Intuition**: Small instances can be solved with standard Dijkstra. We limit exploration to k+1 vertices to maintain efficiency.

### Algorithm 3: BMSSP (Main)

**File**: `sssp/sssp.go:93`

The main recursive Bounded Multi-Source Shortest Path algorithm.

```go
func (s *Solver) BMSSP(l int, B float64, S []int) (float64, []int)
```

**Parameters**:
- `l`: Recursion level (max = ⌈log(n)/t⌉)
- `B`: Upper bound on distances
- `S`: Source set

**Returns**:
- `B'`: Refined upper bound
- `U`: Set of completed vertices

**Algorithm**:

1. **Base case** (l=0):
   ```go
   if l == 0 {
       return s.BaseCase(B, S)
   }
   ```

2. **Find pivots**:
   ```go
   P, W := s.FindPivots(B, S)
   ```

3. **Initialize data structure**:
   ```go
   M := 2^((l-1)·t)
   D := ds.NewDataStructure(M)
   for each x in P:
       D.Insert(x, s.Dist[x])
   ```

4. **Main loop** (pull-recurse-relax):
   ```go
   U := empty set
   limit := k · 2^(l·t)
   
   while len(U) < limit and D not empty:
       Si, Bi := D.Pull()  // Get M smallest
       
       Bi', Ui := s.BMSSP(l-1, Bi, Si)  // Recurse
       
       U = U ∪ Ui
       
       K := empty set
       for each u in Ui:
           for each edge (u,v):
               newDist := s.Dist[u] + weight
               if newDist < s.Dist[v]:
                   s.Dist[v] = newDist
                   if Bi ≤ newDist < B:
                       D.Insert(v, newDist)
                   else if Bi' ≤ newDist < Bi:
                       K = K ∪ {v}
       
       // Batch prepend K and updated Si
       D.BatchPrepend(K ∪ {x ∈ Si : Bi' ≤ d[x] < Bi})
       
       if len(U) > limit:
           return Bi', U ∪ {w ∈ W : d[w] < Bi'}
   ```

5. **Success**: Return B, U ∪ {w ∈ W : d[w] < B}

**Complexity**: O(t + l·log k) per vertex in frontier, amortized

**Key Invariant**: At each iteration, vertices in U have correct distances, and any incomplete vertex's shortest path must go through a vertex in D or W.

## Complexity Analysis

### Parameter Selection

For graph with n vertices:

```go
k := ⌊log^(1/3)(n)⌋
t := ⌊log^(2/3)(n)⌋
l := ⌈log(n)/t⌉ = ⌈log(n)/log^(2/3)(n)⌉ = ⌈log^(1/3)(n)⌉
```

Example for n = 1,000,000:
- log(n) ≈ 13.8
- k ≈ 2.4 → 2
- t ≈ 5.7 → 5
- l ≈ 2.8 → 3 levels

### Time Complexity Breakdown

1. **Transformation**: O(m)

2. **FindPivots** per call: O(k·m') where m' = edges in subgraph
   - Called once per BMSSP recursion
   - Total across all calls: O(m·l) = O(m·log^(1/3) n)

3. **Data structure operations**:
   - Insert: O(log(N/M)) amortized
   - N = frontier size ≤ n
   - M = 2^((l-1)t) ≈ n^(1 - 1/l) = n^(1 - 1/log^(1/3) n)
   - log(N/M) ≈ log(n)/log^(1/3)(n) = log^(2/3)(n)
   - Per vertex: O(log^(2/3) n)
   - Total: O(m·log^(2/3) n)

4. **Recursion overhead**: O(l) = O(log^(1/3) n) per vertex

**Total**: O(m·log^(2/3) n)

### Space Complexity

- Transformed graph: O(m)
- Distance array: O(n)
- Data structure: O(m) in worst case
- Recursion stack: O(l·n) = O(n·log^(1/3) n)

**Total**: O(m + n·log^(1/3) n) = O(m) for sparse graphs

## Implementation Decisions

### 1. Slice vs Map for Sets

**Decision**: Use `map[int]bool` for dynamic sets (U), `[]bool` for static membership (W).

**Rationale**:
- Maps: O(1) insert/lookup, sparse representation
- Boolean arrays: O(1) lookup, dense representation, better cache locality
- W needs fast membership testing during relaxation
- U grows dynamically and may be sparse

### 0. Cycle Detection in Tree Size Calculation

**Decision**: Use three-state memoization (-1 = processing, 0 = unvisited, positive = computed).

**Rationale**:
- Prevents stack overflow from cycles in the graph
- The shortest path forest should be acyclic, but floating-point errors or graph structure can create apparent cycles
- Graceful degradation: treat cyclic nodes as size 1
- Essential for stability on large/complex graphs

### 2. Linked Lists in Blocks

**Decision**: Use singly-linked lists within blocks.

**Rationale**:
- O(1) insertion at head
- No need for random access within blocks
- Memory efficient for variable-size blocks
- Simplifies splitting without copying

### 3. Sorting Strategy

**Decision**: Sort blocks on-demand during Pull.

**Rationale**:
- Amortize sorting cost across multiple operations
- Many blocks may never need full sorting
- Pulling M items at a time amortizes O(M log M) cost

### 4. Early Termination

**Decision**: Multiple early return paths in BMSSP.

**Rationale**:
- If len(P) = 0, no work needed
- If len(U) > limit during loop, frontier too large - return early
- Reduces practical running time on easy instances

### 5. Epsilon Comparison for Floats

**Decision**: Use `math.Abs(diff) < 1e-9` for float equality.

**Rationale**:
- Floating-point arithmetic introduces small errors
- Forest construction needs to detect shortest path edges
- 1e-9 chosen to balance precision and robustness

### 6. Multi-Source Base Case

**Decision**: BaseCase handles sets |S| > 1.

**Rationale**:
- Paper implies S is singleton at l=0, but recursion structure may produce small sets
- Multi-source Dijkstra is natural generalization
- Initialize all sources in priority queue simultaneously

### 7. Distance Initialization

**Decision**: Initialize all distances to Infinity at start, update incrementally.

**Rationale**:
- Standard SSSP convention
- Allows easy detection of unreachable vertices
- Simplifies relaxation condition (newDist < current)

## Testing Strategy

### Correctness Tests

Compare against naive Dijkstra on random graphs:

```go
func TestCorrectness(t *testing.T) {
    for various graph sizes:
        expected := naiveDijkstra(g, source)
        actual := duanAlgorithm(g, source)
        assert distances match within epsilon
}
```

### Performance Tests

1. **Scalability**: Measure time vs n for sparse graphs
2. **Density**: Vary edge-to-vertex ratio
3. **Component analysis**: Time each sub-algorithm
4. **Memory**: Track allocations and peak usage

### Edge Cases

- Empty graph (n=0, m=0)
- Single vertex (n=1, m=0)
- Disconnected components
- Self-loops and multi-edges (handled by transformation)
- Zero-weight edges
- Very large weights (near float64 max)

## Future Optimizations

### Potential Improvements

1. **Parallel Processing**:
   - FindPivots relaxation steps
   - Batch operations in data structure
   - Independent BMSSP recursive calls

2. **Cache Optimization**:
   - Adjacency list layout (CSR format)
   - Block memory layout
   - Prefetching in relaxation loops

3. **Adaptive Parameters**:
   - Choose k, t based on graph properties
   - Different strategies for dense vs sparse regions

4. **Lazy Evaluation**:
   - Defer transformation until needed
   - On-demand block sorting
   - Incremental distance updates

5. **Specialized Cases**:
   - Optimized path for already constant-degree graphs
   - Special handling for DAGs (topological order)
   - Integer weights (word-level tricks)

### Research Directions

- Can we achieve O(m) for directed graphs in comparison-addition model?
- Better practical constants through engineering?
- Parallel/distributed variants?
- Dynamic SSSP with edge updates?

## References

See README.md for complete references to the paper and related work.

---

**Document Version**: 1.0  
**Last Updated**: 2025  
**Implementation**: github.com/phr3nzy/duan-sssp
