package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Instruction[T any] interface {
	Execute(s *state.MachineState[T], p common.BTreeEngine[T]) [][]T
}
