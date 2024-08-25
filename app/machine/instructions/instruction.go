package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Instruction interface {
	Execute(s *state.MachineState, p btree.BTreeEngine) [][]int
}
