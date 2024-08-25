package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type OpenRead struct {
	RootPage uint64
	CursorID uint64
}

var _ Instruction = OpenRead{}

func (openRead OpenRead) Execute(s *state.MachineState, b common.BTreeEngine) [][]int {
	s.CurrentAddress++

	b.NewCursor(openRead.CursorID, openRead.RootPage)

	return [][]int{}
}
