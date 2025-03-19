package graph

import "alon.kr/x/set"

type livenessBuilder struct {
	*Graph

	// Use[v] contains the set of variables that are used before any assignment
	// in the basic block represented by node v.
	Use []set.Set[uint]

	// Def[v] contains the set of variables that are defined, or assigned a value,
	// in the basic block represented by node v.
	Def []set.Set[uint]

	// LiveIn[v] contains the set of variables that are live at the entry of
	// the basic block represented by node v.
	LiveIn []set.Set[uint]

	// LiveOut[v] contains the set of variables that are live at the exit of
	// the basic block represented by node v.
	LiveOut []set.Set[uint]
}

// iterate applies a single iteration of the liveness analysis.
// It returns true if the liveness analysis has converged, i.e., if the
// LiveIn and LiveOut sets have not changed between the current and previous
// iteration.
func (b *livenessBuilder) iterate() bool {
	changed := false

	// Iterate over all nodes in the graph
	for v := range b.Nodes {
		// Save old LiveIn and LiveOut for comparison
		oldLiveIn := b.LiveIn[v].Copy()
		oldLiveOut := b.LiveOut[v].Copy()

		// Compute LiveOut[v] = ⋃(LiveIn[s]) for all successors s of v
		newLiveOut := set.New[uint]()
		for _, s := range b.Nodes[v].ForwardEdges {
			newLiveOut = newLiveOut.Union(b.LiveIn[s])
		}
		b.LiveOut[v] = newLiveOut

		// Compute LiveIn[v] = Use[v] ∪ (LiveOut[v] - Def[v])
		diffSet := b.LiveOut[v].Difference(b.Def[v])
		b.LiveIn[v] = b.Use[v].Union(diffSet)

		// Check if anything changed
		if !oldLiveIn.Equals(b.LiveIn[v]) || !oldLiveOut.Equals(b.LiveOut[v]) {
			changed = true
		}
	}

	return changed
}

func newLivenessBuilder(g *Graph, use []set.Set[uint], def []set.Set[uint]) *livenessBuilder {
	n := g.Size()
	liveIn := make([]set.Set[uint], n)
	liveOut := make([]set.Set[uint], n)

	// Initialize LiveIn and LiveOut sets
	for i := uint(0); i < n; i++ {
		liveIn[i] = set.New[uint]()
		liveOut[i] = set.New[uint]()
	}

	return &livenessBuilder{
		Graph:   g,
		Use:     use,
		Def:     def,
		LiveIn:  liveIn,
		LiveOut: liveOut,
	}
}

func (b *livenessBuilder) iterateUntilFixedpoint() Liveness {
	for b.iterate() {
		// Continue until no more changes (fixed point reached)
	}

	return Liveness{
		Graph:   b.Graph,
		LiveIn:  b.LiveIn,
		LiveOut: b.LiveOut,
	}
}
