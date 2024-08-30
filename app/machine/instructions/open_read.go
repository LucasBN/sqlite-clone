package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type OpenRead[T any] struct {
	RootPage uint64
	CursorID uint64
}

var _ Instruction[any] = OpenRead[any]{}

func (openRead OpenRead[T]) Execute(s *state.MachineState[T], b common.BTreeEngine[T]) [][]T {
	s.CurrentAddress++

	b.NewCursor(openRead.CursorID, openRead.RootPage)

	return [][]T{}
}
