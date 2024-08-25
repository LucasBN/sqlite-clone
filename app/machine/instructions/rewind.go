package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Rewind struct {
	Cursor uint64
}

var _ Instruction = Rewind{}

func (rewind Rewind) Execute(s *state.MachineState, b common.BTreeEngine) [][]int {
	s.CurrentAddress++

	_, err := b.RewindCursor(rewind.Cursor)
	if err != nil {
		panic(err)
	}

	return [][]int{}
}
