package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Column struct {
	Cursor   uint64
	Column   uint64
	Register int
}

var _ Instruction = Column{}

func (column Column) Execute(s *state.MachineState, p btree.BTreeEngine) [][]int {
	s.CurrentAddress++

	entry, err := p.ReadColumn(column.Cursor, column.Column)
	if err != nil {
		panic(err)
	}

	switch entry.(type) {
	case btree.BTreeNumberEntry:
		s.Registers = s.Registers.Set(column.Register, int(entry.(btree.BTreeNumberEntry).Value))
	default:
		panic("Unknown entry type")
	}

	return [][]int{}
}
