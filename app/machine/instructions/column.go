package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Column[T any] struct {
	Cursor   uint64
	Column   uint64
	Register int
}

var _ Instruction[any] = Column[any]{}

func (column Column[T]) Execute(s *state.MachineState[T], b common.BTreeEngine[T]) [][]T {
	s.CurrentAddress++

	entry, err := b.ReadColumn(column.Cursor, column.Column)
	if err != nil {
		panic(err)
	}

	s.Registers = s.Registers.Set(column.Register, entry)

	return [][]T{}
}
