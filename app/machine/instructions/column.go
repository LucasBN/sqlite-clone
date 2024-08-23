package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Column struct {
	Cursor   uint32
	Column   int
	Register int
}

var _ Instruction = Column{}

func (column Column) Execute(s *state.MachineState, p *btree.BTreeProcessor) [][]int {
	s.CurrentAddress++

	cursor := p.GetCursor(column.Cursor)

	value := cursor.ReadColumn(column.Column)

	s.Registers.Set(column.Register, value)

	return [][]int{}
}
