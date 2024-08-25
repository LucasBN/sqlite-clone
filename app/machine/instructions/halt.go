package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Halt struct{}

var _ Instruction = Halt{}

func (Halt) Execute(s *state.MachineState, p btree.BTreeEngine) [][]int {
	s.Halted = true
	return [][]int{}
}
