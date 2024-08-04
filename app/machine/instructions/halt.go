package instructions

import "github/com/lucasbn/sqlite-clone/app/machine/state"

type Halt struct{}

func (Halt) Execute(s *state.MachineState) [][]int {
	s.Halted = true
	return [][]int{}
}
