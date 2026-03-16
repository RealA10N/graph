package graph

type dfsBuilder struct {
	*Graph
	Dfs

	Visited      []bool
	NextPreTime  uint
	NextPostTime uint
}

func newDfsBuilder(g *Graph) dfsBuilder {
	n := g.Size()
	return dfsBuilder{
		Graph: g,
		Dfs: Dfs{
			PreOrder:    make([]uint, n),
			PostOrder:   make([]uint, n),
			Timeline:    make([]uint, 0, 2*n),
			Parent:      make([]uint, n),
			Depth:       make([]uint, n),
			SubtreeSize: make([]uint, n),
		},
		Visited:      make([]bool, n),
		NextPreTime:  0,
		NextPostTime: 0,
	}
}

func (g *dfsBuilder) recordPre(node uint) {
	g.Timeline = append(g.Timeline, node)
	g.PreOrder[node] = g.NextPreTime
	g.NextPreTime++
}

func (g *dfsBuilder) recordPost(node uint) {
	n := uint(len(g.Nodes))
	g.Timeline = append(g.Timeline, n+node)
	g.PostOrder[node] = g.NextPostTime
	g.NextPostTime++
}

func (g *dfsBuilder) dfs(node, from, depth uint) bool {
	if g.Visited[node] {
		return false
	}

	g.Visited[node] = true
	g.Parent[node] = from
	g.Depth[node] = depth
	g.SubtreeSize[node] = 1

	g.recordPre(node)

	for _, next := range g.Nodes[node].ForwardEdges {
		if g.dfs(next, node, depth+1) {
			g.SubtreeSize[node] += g.SubtreeSize[next]
		}
	}

	g.recordPost(node)
	return true
}

func (g *dfsBuilder) toDfs() Dfs {
	g.Dfs.PreOrderReversed = reversePermutation(g.Dfs.PreOrder)
	return g.Dfs
}

// dfsComponentBuilder builds a single DfsComponent. It shares a Visited slice
// with its parent dfsForestBuilder so that separate component runs do not
// revisit already-visited nodes.
type dfsComponentBuilder struct {
	*Graph
	Visited      []bool // shared with the parent forest builder
	root         uint
	PreOrder     []uint
	PostOrder    []uint
	Parent       []uint
	Depth        []uint
	SubtreeSize  []uint
	Timeline     []uint
	NextPreTime  uint
	NextPostTime uint
}

func newDfsComponentBuilder(g *Graph, visited []bool, root uint) dfsComponentBuilder {
	n := g.Size()
	preOrder := make([]uint, n)
	postOrder := make([]uint, n)
	parent := make([]uint, n)
	for i := range preOrder {
		preOrder[i] = DfsUnvisited
		postOrder[i] = DfsUnvisited
		parent[i] = DfsUnvisited
	}
	return dfsComponentBuilder{
		Graph:    g,
		Visited:  visited,
		root:     root,
		PreOrder: preOrder,
		PostOrder: postOrder,
		Parent:   parent,
		Depth:    make([]uint, n),
		SubtreeSize: make([]uint, n),
	}
}

func (b *dfsComponentBuilder) recordPre(node uint) {
	b.Timeline = append(b.Timeline, node)
	b.PreOrder[node] = b.NextPreTime
	b.NextPreTime++
}

func (b *dfsComponentBuilder) recordPost(node uint) {
	n := uint(len(b.Graph.Nodes))
	b.Timeline = append(b.Timeline, n+node)
	b.PostOrder[node] = b.NextPostTime
	b.NextPostTime++
}

func (b *dfsComponentBuilder) dfs(node, from, depth uint) {
	b.Visited[node] = true
	b.Parent[node] = from
	b.Depth[node] = depth
	b.SubtreeSize[node] = 1

	b.recordPre(node)

	for _, next := range b.Graph.Nodes[node].ForwardEdges {
		if !b.Visited[next] {
			b.dfs(next, node, depth+1)
			b.SubtreeSize[node] += b.SubtreeSize[next]
		}
	}

	b.recordPost(node)
}

func (b *dfsComponentBuilder) toComponent() DfsComponent {
	componentSize := b.NextPreTime
	preOrderReversed := make([]uint, componentSize)
	for node, preIdx := range b.PreOrder {
		if preIdx != DfsUnvisited {
			preOrderReversed[preIdx] = uint(node)
		}
	}
	return DfsComponent{
		Root:             b.root,
		PreOrder:         b.PreOrder,
		PostOrder:        b.PostOrder,
		PreOrderReversed: preOrderReversed,
		Timeline:         b.Timeline,
		Parent:           b.Parent,
		Depth:            b.Depth,
		SubtreeSize:      b.SubtreeSize,
	}
}

// dfsForestBuilder orchestrates a full DFS forest traversal over all nodes.
type dfsForestBuilder struct {
	*Graph
	Visited []bool
}

func newDfsForestBuilder(g *Graph) dfsForestBuilder {
	return dfsForestBuilder{
		Graph:   g,
		Visited: make([]bool, g.Size()),
	}
}

func (fb *dfsForestBuilder) toDfsForest() DfsForest {
	n := fb.Size()
	var components []DfsComponent
	for i := uint(0); i < n; i++ {
		if !fb.Visited[i] {
			cb := newDfsComponentBuilder(fb.Graph, fb.Visited, i)
			cb.dfs(i, i, 0)
			components = append(components, cb.toComponent())
		}
	}
	return DfsForest(components)
}
