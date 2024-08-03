package instructions

import "github/com/lucasbn/sqlite-clone/app/machine/state"

type Integer struct {
	Register int
	Value    int
}

func (integer Integer) Execute(s *state.MachineState) [][]int {
	s.CurrentAddress++

	s.Registers = s.Registers.Set(integer.Register, integer.Value)
	return [][]int{}
}
