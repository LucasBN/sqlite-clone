package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Integer struct {
	Register int
	Value    int
}

var _ Instruction = Integer{}

func (integer Integer) Execute(s *state.MachineState, b common.BTreeEngine) [][]int {
	s.CurrentAddress++

	s.Registers = s.Registers.Set(integer.Register, integer.Value)
	return [][]int{}
}
