package graph_test

import (
	"testing"

	"alon.kr/x/graph"
	"github.com/stretchr/testify/assert"
)

func TestLineDfs(t *testing.T) {
	// 0 -> 1 -> 2

	n := uint(3)
	g := graph.NewGraph([][]uint{{1}, {2}, {}})
	dfs := g.Dfs(0)

	assert.EqualValues(t, []uint{0, 1, 2}, dfs.PreOrder)
	assert.EqualValues(t, []uint{2, 1, 0}, dfs.PostOrder)
	assert.EqualValues(t, []uint{0, 1, 2, n + 2, n + 1, n + 0}, dfs.Timeline)
	assert.EqualValues(t, []uint{0, 0, 1}, dfs.Parent)
	assert.EqualValues(t, []uint{0, 1, 2}, dfs.Depth)
}

func TestBinaryTree(t *testing.T) {
	//    0
	//  1   2
	// 3 4 5 6

	n := uint(7)
	g := graph.NewGraph([][]uint{{1, 2}, {3, 4}, {5, 6}, {}, {}, {}, {}})
	dfs := g.Dfs(0)

	//                    index: 0  1  2  3  4  5  6
	assert.EqualValues(t, []uint{0, 1, 4, 2, 3, 5, 6}, dfs.PreOrder)
	assert.EqualValues(t, []uint{6, 2, 5, 0, 1, 3, 4}, dfs.PostOrder)
	assert.EqualValues(t, []uint{0, 1, 3, n + 3, 4, n + 4, n + 1, 2, 5, n + 5, 6, n + 6, n + 2, n + 0}, dfs.Timeline)
	assert.EqualValues(t, []uint{0, 0, 0, 1, 1, 2, 2}, dfs.Parent)
	assert.EqualValues(t, []uint{0, 1, 1, 2, 2, 2, 2}, dfs.Depth)
}

func TestDfsForestSingleComponent(t *testing.T) {
	// 0 -> 1 -> 2 (connected graph, should produce one component identical to Dfs)

	n := uint(3)
	g := graph.NewGraph([][]uint{{1}, {2}, {}})
	forest := g.DfsForest()

	assert.Len(t, forest, 1)

	c := forest[0]
	assert.EqualValues(t, 0, c.Root)
	assert.EqualValues(t, []uint{0, 1, 2}, c.PreOrder)
	assert.EqualValues(t, []uint{2, 1, 0}, c.PostOrder)
	assert.EqualValues(t, []uint{0, 1, 2, n+2, n+1, n+0}, c.Timeline)
	assert.EqualValues(t, []uint{0, 1, 2}, c.PreOrderReversed)
	assert.EqualValues(t, []uint{0, 0, 1}, c.Parent)
	assert.EqualValues(t, []uint{0, 1, 2}, c.Depth)
}

func TestDfsForestDisconnected(t *testing.T) {
	// Two disconnected components: 0 -> 1 and 2 -> 3
	//
	// Component 0: root=0, nodes={0,1}
	// Component 1: root=2, nodes={2,3}

	n := uint(4)
	g := graph.NewGraph([][]uint{{1}, {}, {3}, {}})
	forest := g.DfsForest()

	assert.Len(t, forest, 2)

	// Component rooted at 0
	c0 := forest[0]
	assert.EqualValues(t, 0, c0.Root)
	assert.EqualValues(t, []uint{0, 1, graph.DfsUnvisited, graph.DfsUnvisited}, c0.PreOrder)
	assert.EqualValues(t, []uint{1, 0, graph.DfsUnvisited, graph.DfsUnvisited}, c0.PostOrder)
	assert.EqualValues(t, []uint{0, 1, n+1, n+0}, c0.Timeline)
	assert.EqualValues(t, []uint{0, 1}, c0.PreOrderReversed)
	assert.EqualValues(t, []uint{0, 0, graph.DfsUnvisited, graph.DfsUnvisited}, c0.Parent)
	assert.EqualValues(t, []uint{0, 1, 0, 0}, c0.Depth)

	// Component rooted at 2
	c2 := forest[1]
	assert.EqualValues(t, 2, c2.Root)
	assert.EqualValues(t, []uint{graph.DfsUnvisited, graph.DfsUnvisited, 0, 1}, c2.PreOrder)
	assert.EqualValues(t, []uint{graph.DfsUnvisited, graph.DfsUnvisited, 1, 0}, c2.PostOrder)
	assert.EqualValues(t, []uint{2, 3, n+3, n+2}, c2.Timeline)
	assert.EqualValues(t, []uint{2, 3}, c2.PreOrderReversed)
	assert.EqualValues(t, []uint{graph.DfsUnvisited, graph.DfsUnvisited, 2, 2}, c2.Parent)
	assert.EqualValues(t, []uint{0, 0, 0, 1}, c2.Depth)
}

func TestDfsForestThreeIsolatedNodes(t *testing.T) {
	// Three isolated nodes: 0, 1, 2 — each is its own component

	n := uint(3)
	g := graph.NewGraph([][]uint{{}, {}, {}})
	forest := g.DfsForest()

	assert.Len(t, forest, 3)

	for i, c := range forest {
		root := uint(i)
		assert.EqualValues(t, root, c.Root)
		assert.EqualValues(t, []uint{root}, c.PreOrderReversed)

		expected := make([]uint, n)
		for j := range expected {
			expected[j] = graph.DfsUnvisited
		}
		expected[root] = 0

		assert.EqualValues(t, expected, c.PreOrder)
		assert.EqualValues(t, expected, c.PostOrder)
		assert.EqualValues(t, []uint{root, n + root}, c.Timeline)
	}
}

func TestSpecialEdges(t *testing.T) {
	//    0      with an additional cross edge 5 -> 1
	//  1   2    a back edge 4 -> 0
	// 3 4 5 6   and a forward edge 0 -> 6

	n := uint(7)
	g := graph.NewGraph([][]uint{
		{1, 2, 6}, {3, 4}, {5, 6}, {}, {0}, {1}, {},
	})
	dfs := g.Dfs(0)

	//                    index: 0  1  2  3  4  5  6
	assert.EqualValues(t, []uint{0, 1, 4, 2, 3, 5, 6}, dfs.PreOrder)
	assert.EqualValues(t, []uint{6, 2, 5, 0, 1, 3, 4}, dfs.PostOrder)
	assert.EqualValues(t, []uint{0, 1, 3, n + 3, 4, n + 4, n + 1, 2, 5, n + 5, 6, n + 6, n + 2, n + 0}, dfs.Timeline)
	assert.EqualValues(t, []uint{0, 0, 0, 1, 1, 2, 2}, dfs.Parent)
	assert.EqualValues(t, []uint{0, 1, 1, 2, 2, 2, 2}, dfs.Depth)
}
