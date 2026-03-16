package graph_test

import (
	"testing"

	"alon.kr/x/graph"
	"github.com/stretchr/testify/assert"
)

func TestDfsForestTwoDisconnectedChains(t *testing.T) {
	// 0 -> 1    2 -> 3  (two disconnected chains)

	n := uint(4)
	g := graph.NewGraph([][]uint{{1}, {}, {3}, {}})
	dfs := g.DfsForest()

	assert.EqualValues(t, []uint{0, 1, 2, 3}, dfs.PreOrder)
	assert.EqualValues(t, []uint{1, 0, 3, 2}, dfs.PostOrder)
	assert.EqualValues(t, []uint{0, 1, n + 1, n + 0, 2, 3, n + 3, n + 2}, dfs.Timeline)
	assert.EqualValues(t, []uint{0, 0, 2, 2}, dfs.Parent)
	assert.EqualValues(t, []uint{0, 1, 0, 1}, dfs.Depth)
	assert.EqualValues(t, []uint{2, 1, 2, 1}, dfs.SubtreeSize)
}

func TestDfsForestIsolatedNodes(t *testing.T) {
	// 0    1    2  (three isolated nodes)

	n := uint(3)
	g := graph.NewGraph([][]uint{{}, {}, {}})
	dfs := g.DfsForest()

	assert.EqualValues(t, []uint{0, 1, 2}, dfs.PreOrder)
	assert.EqualValues(t, []uint{0, 1, 2}, dfs.PostOrder)
	assert.EqualValues(t, []uint{0, n + 0, 1, n + 1, 2, n + 2}, dfs.Timeline)
	assert.EqualValues(t, []uint{0, 1, 2}, dfs.Parent)
	assert.EqualValues(t, []uint{0, 0, 0}, dfs.Depth)
	assert.EqualValues(t, []uint{1, 1, 1}, dfs.SubtreeSize)
}

func TestDfsForestSingleComponent(t *testing.T) {
	// 0 -> 1 -> 2  (single chain, result should match Dfs(0))

	g := graph.NewGraph([][]uint{{1}, {2}, {}})
	dfsForest := g.DfsForest()
	dfs := g.Dfs(0)

	assert.EqualValues(t, dfs.PreOrder, dfsForest.PreOrder)
	assert.EqualValues(t, dfs.PostOrder, dfsForest.PostOrder)
	assert.EqualValues(t, dfs.Timeline, dfsForest.Timeline)
	assert.EqualValues(t, dfs.Parent, dfsForest.Parent)
	assert.EqualValues(t, dfs.Depth, dfsForest.Depth)
	assert.EqualValues(t, dfs.SubtreeSize, dfsForest.SubtreeSize)
	assert.EqualValues(t, dfs.PreOrderReversed, dfsForest.PreOrderReversed)
}

func TestDfsForestReverseChain(t *testing.T) {
	// 2 -> 1 -> 0  (edges point toward lower indices)
	// DfsForest visits 0 first, so 1 and 2 each become their own roots
	// when reached (their forward edges point to already-visited nodes).

	n := uint(3)
	g := graph.NewGraph([][]uint{{}, {0}, {1}})
	dfs := g.DfsForest()

	assert.EqualValues(t, []uint{0, 1, 2}, dfs.PreOrder)
	assert.EqualValues(t, []uint{0, 1, 2}, dfs.PostOrder)
	assert.EqualValues(t, []uint{0, n + 0, 1, n + 1, 2, n + 2}, dfs.Timeline)
	assert.EqualValues(t, []uint{0, 1, 2}, dfs.Parent)
	assert.EqualValues(t, []uint{0, 0, 0}, dfs.Depth)
	assert.EqualValues(t, []uint{1, 1, 1}, dfs.SubtreeSize)
}

func TestDfsForestEmpty(t *testing.T) {
	g := graph.NewGraph([][]uint{})
	dfs := g.DfsForest()

	assert.Empty(t, dfs.PreOrder)
	assert.Empty(t, dfs.PostOrder)
	assert.Empty(t, dfs.Timeline)
	assert.Empty(t, dfs.Parent)
	assert.Empty(t, dfs.Depth)
	assert.Empty(t, dfs.SubtreeSize)
	assert.Empty(t, dfs.PreOrderReversed)
}
