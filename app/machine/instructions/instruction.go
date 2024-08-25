package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Instruction interface {
	Execute(s *state.MachineState, p common.BTreeEngine) [][]int
}
