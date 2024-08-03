package instructions

import "github/com/lucasbn/sqlite-clone/app/machine/state"

type ResultRow struct {
	FromRegister int
	ToRegister   int
}

func (resultRow ResultRow) Execute(s *state.MachineState) [][]int {
	s.CurrentAddress++

	var result []int
	for i := resultRow.FromRegister; i <= resultRow.ToRegister; i++ {
		result = append(result, s.Registers.Get(i))
	}
	return [][]int{result}
}
