package graph_test

import (
	"testing"

	"alon.kr/x/graph"
	"alon.kr/x/set"
	"github.com/stretchr/testify/assert"
)

func TestLivenessSimple(t *testing.T) {
	// Simple linear graph: 0 -> 1
	// Variable 0 is defined in node 0 and used in node 1

	g := graph.NewGraph([][]uint{{1}, {}})

	def := []set.Set[uint]{
		set.FromSlice([]uint{0}), // Node 0: var 0 is defined
		set.FromSlice([]uint{}),  // Node 1: nothing defined
	}

	use := []set.Set[uint]{
		set.FromSlice([]uint{}),  // Node 0: nothing used
		set.FromSlice([]uint{0}), // Node 1: var 0 is used
	}

	liveness := g.LivenessAnalysis(use, def)

	assert.ElementsMatch(t, []uint{}, liveness.LiveIn[0].ToSlice())
	assert.ElementsMatch(t, []uint{0}, liveness.LiveIn[1].ToSlice())
	assert.ElementsMatch(t, []uint{0}, liveness.LiveOut[0].ToSlice())
	assert.ElementsMatch(t, []uint{}, liveness.LiveOut[1].ToSlice())
}

func TestLivenessLoop(t *testing.T) {
	// Loop graph: 0 -> 1 -> 2
	//             ^---------+
	// Variables:
	// - var 0 is defined in node 0
	// - var 1 is defined in node 1 and used in node 2
	// - var 2 is defined in node 2 and used in node 1

	g := graph.NewGraph([][]uint{{1}, {2}, {1}})

	def := []set.Set[uint]{
		set.FromSlice([]uint{0}), // Node 0: var 0 is defined
		set.FromSlice([]uint{1}), // Node 1: var 1 is defined
		set.FromSlice([]uint{2}), // Node 2: var 2 is defined
	}

	use := []set.Set[uint]{
		set.FromSlice([]uint{}),  // Node 0: nothing used
		set.FromSlice([]uint{2}), // Node 1: var 2 is used
		set.FromSlice([]uint{1}), // Node 2: var 1 is used
	}

	liveness := g.LivenessAnalysis(use, def)

	assert.ElementsMatch(t, []uint{2}, liveness.LiveIn[0].ToSlice())
	assert.ElementsMatch(t, []uint{2}, liveness.LiveIn[1].ToSlice())
	assert.ElementsMatch(t, []uint{1}, liveness.LiveIn[2].ToSlice())

	assert.ElementsMatch(t, []uint{2}, liveness.LiveOut[0].ToSlice())
	assert.ElementsMatch(t, []uint{1}, liveness.LiveOut[1].ToSlice())
	assert.ElementsMatch(t, []uint{2}, liveness.LiveOut[2].ToSlice())
}

func TestLivenessIfElse(t *testing.T) {
	// If-else graph: 0 -> 1 -> 3
	//                +--> 2 ---^
	// Variables:
	// - var 0 is defined in node 0 and used in nodes 1 and 2
	// - var 1 is defined in node 0 and used in node 1
	// - var 2 is defined in node 0 and used in node 2
	// - var 3 is defined in node 0 and used in node 3

	g := graph.NewGraph([][]uint{{1, 2}, {3}, {3}, {}})

	def := []set.Set[uint]{
		set.FromSlice([]uint{0, 1, 2, 3}), // Node 0: vars 0,1,2,3 are defined
		set.FromSlice([]uint{}),           // Node 1: nothing defined
		set.FromSlice([]uint{}),           // Node 2: nothing defined
		set.FromSlice([]uint{}),           // Node 3: nothing defined
	}

	use := []set.Set[uint]{
		set.FromSlice([]uint{}),     // Node 0: nothing used
		set.FromSlice([]uint{0, 1}), // Node 1: vars 0,1 are used
		set.FromSlice([]uint{0, 2}), // Node 2: vars 0,2 are used
		set.FromSlice([]uint{3}),    // Node 3: var 3 is used
	}

	liveness := g.LivenessAnalysis(use, def)

	assert.ElementsMatch(t, []uint{}, liveness.LiveIn[0].ToSlice())
	assert.ElementsMatch(t, []uint{0, 1, 3}, liveness.LiveIn[1].ToSlice())
	assert.ElementsMatch(t, []uint{0, 2, 3}, liveness.LiveIn[2].ToSlice())
	assert.ElementsMatch(t, []uint{3}, liveness.LiveIn[3].ToSlice())

	assert.ElementsMatch(t, []uint{0, 1, 2, 3}, liveness.LiveOut[0].ToSlice())
	assert.ElementsMatch(t, []uint{3}, liveness.LiveOut[1].ToSlice())
	assert.ElementsMatch(t, []uint{3}, liveness.LiveOut[2].ToSlice())
	assert.ElementsMatch(t, []uint{}, liveness.LiveOut[3].ToSlice())
}
