package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Halt struct{}

var _ Instruction = Halt{}

func (Halt) Execute(s *state.MachineState, b common.BTreeEngine) [][]int {
	s.Halted = true
	return [][]int{}
}
