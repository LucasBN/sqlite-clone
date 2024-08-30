package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Rewind[T any] struct {
	Cursor uint64
}

var _ Instruction[any] = Rewind[any]{}

func (rewind Rewind[T]) Execute(s *state.MachineState[T], b common.BTreeEngine[T]) [][]T {
	s.CurrentAddress++

	_, err := b.RewindCursor(rewind.Cursor)
	if err != nil {
		panic(err)
	}

	return [][]T{}
}
