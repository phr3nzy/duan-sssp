package graph

// Edge represents a weighted directed connection.
type Edge struct {
	To     int
	Weight float64
}

// Graph uses an adjacency list.
type Graph struct {
	V   int
	Adj [][]Edge
}

func NewGraph(v int) *Graph {
	return &Graph{
		V:   v,
		Adj: make([][]Edge, v),
	}
}

func (g *Graph) AddEdge(u, v int, w float64) {
	g.Adj[u] = append(g.Adj[u], Edge{To: v, Weight: w})
}

// TransformedGraph holds the new graph and mapping data.
type TransformedGraph struct {
	G           *Graph
	OriginalTo  []int // Map original ID -> Start node in cycle
	NewToOrigin []int // Map new ID -> Original ID
}

// ToConstantDegree implements the transformation described in the paper.
// Each vertex v is replaced by a cycle of nodes, one for each edge.
func (g *Graph) ToConstantDegree() *TransformedGraph {
	// 1. Calculate size of new graph
	// Each original vertex v needs degree(v) + constant auxiliary nodes.
	// Simple strategy:
	// For vertex v with In-degree I and Out-degree O:
	// Create a "hub" cycle of size max(1, I+O).

	// Faster Approach:
	// Simply split node u into u_in and u_out is not enough for constant degree.
	// We must chain edges.

	// Implementation of [Fre83] style transformation:
	// For each node u, create a chain/cycle of auxiliary nodes.
	// Total new nodes approx = 2*|E| + |V|.

	// We need to group edges.
	// Construct simplified expansion:
	// For each node u:
	//   Create a "chain" of nodes u_0, u_1, ... u_k where k = out_degree.
	//   Edge (u,v) becomes (u_i, v_0) with weight w.
	//   Chain edges (u_i, u_{i+1}) have weight 0.
	//   Incoming edges?
	//   The paper says: "Substitute each vertex v with a cycle... For every neighbor w there is a vertex x_vw".

	// Let's use a simpler gadget:
	// Every original node u becomes a cycle of k nodes, where k = InDegree(u) + OutDegree(u).
	// If k=0, just 1 node.

	inDegree := make([]int, g.V)
	for u := 0; u < g.V; u++ {
		for _, e := range g.Adj[u] {
			inDegree[e.To]++
		}
	}

	starts := make([]int, g.V)
	sizes := make([]int, g.V)

	currentID := 0
	for u := 0; u < g.V; u++ {
		starts[u] = currentID
		sz := len(g.Adj[u]) + inDegree[u]
		if sz == 0 {
			sz = 1
		}
		sizes[u] = sz
		currentID += sz
	}

	newG := NewGraph(currentID)
	newToOrigin := make([]int, currentID)

	// Build Cycles and internal mappings
	// Map (u, v) edge to specific index in u's cycle (outgoing) and v's cycle (incoming)

	// We need to assign specific "slots" in the cycle for each edge.
	// slots[u] tracks next available slot for node u.
	slots := make([]int, g.V)

	// Create zero-weight cycles/chains
	for u := 0; u < g.V; u++ {
		start := starts[u]
		sz := sizes[u]
		for i := 0; i < sz; i++ {
			curr := start + i
			next := start + (i+1)%sz
			newG.AddEdge(curr, next, 0)
			newToOrigin[curr] = u
		}
	}

	// Add real edges
	for u := 0; u < g.V; u++ {
		for _, e := range g.Adj[u] {
			v := e.To
			w := e.Weight

			// u's slot for this outgoing edge
			uSlot := slots[u]
			slots[u]++
			uNode := starts[u] + uSlot

			// v's slot for this incoming edge
			vSlot := slots[v]
			slots[v]++
			vNode := starts[v] + vSlot

			newG.AddEdge(uNode, vNode, w)
		}
	}

	return &TransformedGraph{
		G:           newG,
		OriginalTo:  starts,
		NewToOrigin: newToOrigin,
	}
}

// MapDistances converts distances from the transformed graph back to the original.
// If target is provided with enough capacity, it will be reused to avoid allocation.
func (tg *TransformedGraph) MapDistances(dist []float64, target ...[]float64) []float64 {
	var res []float64
	if len(target) > 0 && cap(target[0]) >= len(tg.OriginalTo) {
		res = target[0][:len(tg.OriginalTo)]
	} else {
		res = make([]float64, len(tg.OriginalTo))
	}

	for i, startNode := range tg.OriginalTo {
		// The distance to original node i is the min distance to any node in its cycle
		// Or simply the distance to the "start" node of the cycle (since internal weights are 0)
		res[i] = dist[startNode]
	}
	return res
}
