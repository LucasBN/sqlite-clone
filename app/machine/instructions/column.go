package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
	"github/com/lucasbn/sqlite-clone/app/types"
)

type Column struct {
	Cursor   uint64
	Column   uint64
	Register int
}

var _ Instruction = Column{}

func (column Column) Execute(s *state.MachineState, b common.BTreeEngine) [][]int {
	s.CurrentAddress++

	entry, err := b.ReadColumn(column.Cursor, column.Column)
	if err != nil {
		panic(err)
	}

	switch e := entry.(type) {
	case types.NumberEntry:
		s.Registers = s.Registers.Set(column.Register, int(e.Value))
	default:
		panic("Unknown entry type")
	}

	return [][]int{}
}
