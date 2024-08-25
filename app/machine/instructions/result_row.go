package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type ResultRow struct {
	FromRegister int
	ToRegister   int
}

var _ Instruction = ResultRow{}

func (resultRow ResultRow) Execute(s *state.MachineState, b common.BTreeEngine) [][]int {
	s.CurrentAddress++

	var result []int
	for i := resultRow.FromRegister; i <= resultRow.ToRegister; i++ {
		result = append(result, s.Registers.Get(i))
	}
	return [][]int{result}
}
