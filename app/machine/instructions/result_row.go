package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type ResultRow[T any] struct {
	FromRegister int
	ToRegister   int
}

var _ Instruction[any] = ResultRow[any]{}

func (resultRow ResultRow[T]) Execute(s *state.MachineState[T], b common.BTreeEngine[T]) [][]T {
	s.CurrentAddress++

	var result []T
	for i := resultRow.FromRegister; i <= resultRow.ToRegister; i++ {
		result = append(result, s.Registers.Get(i))
	}
	return [][]T{result}
}
