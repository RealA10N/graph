package graph

import "alon.kr/x/set"

type Liveness struct {
	*Graph

	// LiveIn[v] contains the set of variables that are live at the entry of
	// the basic block represented by node v.
	LiveIn []set.Set[uint]

	// LiveOut[v] contains the set of variables that are live at the exit of
	// the basic block represented by node v.
	LiveOut []set.Set[uint]
}
