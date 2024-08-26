package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Next[T any] struct {
	Cursor      uint64
	FromAddress uint64
}

var _ Instruction[any] = Next[any]{}

func (next Next[T]) Execute(s *state.MachineState[T], b common.BTreeEngine[T]) [][]T {
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

	return [][]T{}
}
