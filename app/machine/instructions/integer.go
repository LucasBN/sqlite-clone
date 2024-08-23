package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Integer struct {
	Register int
	Value    int
}

var _ Instruction = Integer{}

func (integer Integer) Execute(s *state.MachineState, p *btree.BTreeProcessor) [][]int {
	s.CurrentAddress++

	s.Registers = s.Registers.Set(integer.Register, integer.Value)
	return [][]int{}
}
