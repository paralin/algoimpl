package graph

const (
	unseen = iota
	seen
)

// O(V + E)
func (g *Graph) dfs(node *Node, finishList *[]*Node) {
	node.state = seen
	for _, adjacentNode := range g.adjacents[node] {
		if adjacentNode.state == unseen {
			adjacentNode.parent = node
			g.dfs(adjacentNode, finishList)
		}
	}
	*finishList = append(*finishList, node)
}

// Topologically sorts a directed acyclic graph.
// If the graph is cyclic, the sort order will change
// based on which node the sort starts on. O(V+E) complexity.
func TopologicalSort(g *Graph) []*Node {
	g.lazyInit()
	sorted := make([]*Node, 0, len(g.nodes))
	// sort preorder (first jacket, then shirt)
	for _, node := range g.nodes {
		if node.state == unseen {
			g.dfs(node, &sorted)
		}
	}
	// now make post order for correct sort (jacket follows shirt). O(V)
	length := len(sorted)
	for i := 0; i < length/2; i++ {
		sorted[i], sorted[length-i-1] = sorted[length-i-1], sorted[i]
	}
	return sorted
}

// Returns reversed copy of the directed graph g. O(V+E) complexity.
func Reverse(g *Graph) *Graph {
	reversed, _ := New("directed")
	// O(V)
	for _ = range g.nodes {
		reversed.MakeNode()
	}
	// O(V + E)
	for _, node := range g.nodes {
		for _, adjacent := range g.adjacents[node] {
			reversed.Connect(reversed.nodes[adjacent.Index], reversed.nodes[node.Index])
		}
	}
	return reversed
}

// Returns a slice of strongly connected nodes on a directed graph.
// If passed an undirected graph, returns nil.
// The returned components have reversed, nonexclusive edges.
// For example, if this is passed the graph
//     a->b, c
//     b->a, c
//     c
// will return components
//     [[c->a, b], [b->a], [a->b]]
// where -> represents the edges that the node contains.
// O(V+E) complexity.
func StronglyConnectedComponents(g *Graph) [][]*Node {
	if g.kind == 0 {
		return nil
	}
	g.lazyInit()
	components := make([][]*Node, 0)
	finishOrder := TopologicalSort(g)
	for i := range finishOrder {
		finishOrder[i].state = unseen
	}
	// creates a reversed graph with empty parents
	reversed := Reverse(g)
	for _, sink := range finishOrder {
		if reversed.nodes[sink.Index].state == unseen {
			component := make([]*Node, 0)
			reversed.dfs(reversed.nodes[sink.Index], &component)
			components = append(components, component)
		}
	}
	return components
}
