package graph

// DfsUnvisited is a sentinel value used in DfsComponent arrays to indicate
// that a node does not belong to the component.
const DfsUnvisited = ^uint(0)

// DfsComponent represents the DFS spanning tree rooted at a single node,
// covering one weakly-connected component of the graph.
type DfsComponent struct {
	// Root is the root node of this DFS spanning tree.
	Root uint

	// PreOrder[v] is the component-local index of v in the DFS pre-order
	// traversal. DfsUnvisited for nodes not in this component.
	PreOrder []uint

	// PreOrderReversed[i] is the node with the i-th index in the component's
	// pre-order. Only contains nodes in this component.
	PreOrderReversed []uint

	// PostOrder[v] is the component-local index of v in the DFS post-order
	// traversal. DfsUnvisited for nodes not in this component.
	PostOrder []uint

	// Timeline is a slice of length 2k (where k is the component size),
	// recording each push and pop during DFS. Uses the same encoding as
	// Dfs.Timeline, with pop offsets relative to the full graph size n.
	Timeline []uint

	// Parent[v] is the parent of v in the DFS spanning tree.
	// The root's parent is itself. DfsUnvisited for nodes not in this component.
	Parent []uint

	// Depth[v] is the depth of v in the DFS spanning tree.
	// 0 for nodes not in this component (also 0 for the root; use Parent to
	// distinguish).
	Depth []uint

	// SubtreeSize[v] is the size of the subtree of v, including v itself.
	// 0 for nodes not in this component.
	SubtreeSize []uint
}

// DfsForest is the result of a complete DFS traversal of all nodes in a graph,
// covering every weakly-connected component. Each element is the DFS spanning
// tree for one component.
type DfsForest []DfsComponent

type Dfs struct {
	// PreOrder[v] is the index of v in the DFS pre-order traversal.
	PreOrder []uint

	// PreOrderReversed[i] is the node with the i-th index in the preorder.
	PreOrderReversed []uint

	// PostOrder[v] is the index of v in the DFS post-order traversal.
	PostOrder []uint

	// Timeline is a slice of length 2n, where each push and pop operation
	// during the DFS traversal of the stack is recorded.
	//
	// Operations are ordered in chronological order. If Timeline[i] = k,
	// where k in [0, n), then the i-th operation in the traversal was a push
	// operation of the value k. Otherwise, if k in [n, 2n), then it represents
	// a pop of the value k-n.
	Timeline []uint

	// Parent[v] is the parent of v in the DFS spanning tree.
	// The parent of the root of the tree is itself.
	Parent []uint

	// Depth[v] is the depth of v in the DFS spanning tree.
	Depth []uint

	// SubtreeSize[v] is the size of the subtree of v in the DFS spanning tree,
	// including v itself.
	SubtreeSize []uint
}

func (d *Dfs) IsAncestor(ancestor uint, descendant uint) bool {
	return d.PreOrder[ancestor] <= d.PreOrder[descendant] &&
		d.PostOrder[ancestor] >= d.PostOrder[descendant]
}

func (d *Dfs) IsStrictAncestor(ancestor uint, descendant uint) bool {
	return d.PreOrder[ancestor] < d.PreOrder[descendant] &&
		d.PostOrder[ancestor] > d.PostOrder[descendant]
}

func (d *Dfs) IsDeeper(deeper uint, shallower uint) bool {
	return d.Depth[deeper] >= d.Depth[shallower]
}

func (d *Dfs) IsStrictlyDeeper(deeper uint, shallower uint) bool {
	return d.Depth[deeper] > d.Depth[shallower]
}

func (d *Dfs) Subtree(v uint) []uint {
	start := d.PreOrder[v]
	end := start + d.SubtreeSize[v]
	return d.PreOrderReversed[start:end]
}
