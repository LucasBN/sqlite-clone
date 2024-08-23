package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type ResultRow struct {
	FromRegister int
	ToRegister   int
}

var _ Instruction = ResultRow{}

func (resultRow ResultRow) Execute(s *state.MachineState, p *btree.BTreeProcessor) [][]int {
	s.CurrentAddress++

	var result []int
	for i := resultRow.FromRegister; i <= resultRow.ToRegister; i++ {
		result = append(result, s.Registers.Get(i))
	}
	return [][]int{result}
}
