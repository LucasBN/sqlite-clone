package instructions

import "github/com/lucasbn/sqlite-clone/app/engine/state"

type Instruction interface {
	Execute(s state.MachineState) state.MachineState
}
