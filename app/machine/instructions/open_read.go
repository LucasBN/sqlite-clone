package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type OpenRead struct {
	RootPage uint64
	CursorID uint64
}

var _ Instruction = OpenRead{}

func (openRead OpenRead) Execute(s *state.MachineState, p btree.BTreeEngine) [][]int {
	s.CurrentAddress++

	p.NewCursor(openRead.CursorID, openRead.RootPage)

	return [][]int{}
}
