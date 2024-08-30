package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Halt[T any] struct{}

var _ Instruction[any] = Halt[any]{}

func (Halt[T]) Execute(s *state.MachineState[T], b common.BTreeEngine[T]) [][]T {
	s.Halted = true
	return [][]T{}
}
