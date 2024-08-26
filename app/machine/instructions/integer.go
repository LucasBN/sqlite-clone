package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Integer[T any] struct {
	Register int
	Value    T
}

var _ Instruction[int] = Integer[int]{}

func (integer Integer[T]) Execute(s *state.MachineState[T], b common.BTreeEngine[T]) [][]T {
	s.CurrentAddress++

	s.Registers = s.Registers.Set(integer.Register, integer.Value)
	return [][]T{}
}
