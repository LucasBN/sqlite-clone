package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Rewind struct {
	Cursor uint64
}

var _ Instruction = Next{}

func (rewind Rewind) Execute(s *state.MachineState, p btree.BTreeEngine) [][]int {
	s.CurrentAddress++

	_, err := p.RewindCursor(rewind.Cursor)
	if err != nil {
		panic(err)
	}

	return [][]int{}
}
