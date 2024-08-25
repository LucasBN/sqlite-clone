package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Next struct {
	Cursor      uint64
	FromAddress uint64
}

var _ Instruction = Next{}

func (next Next) Execute(s *state.MachineState, b common.BTreeEngine) [][]int {
	didAdvance, err := b.AdvanceCursor(next.Cursor)
	if err != nil {
		panic(err)
	}
	// If we couldn't advanced the cursor because it was already at the end, we
	// want to jump to the next instruction. Otherwise, we want to jump to the
	// address specified in the instruction.
	if !didAdvance {
		s.CurrentAddress++
	} else {
		s.CurrentAddress = int(next.FromAddress)
	}

	return [][]int{}
}
