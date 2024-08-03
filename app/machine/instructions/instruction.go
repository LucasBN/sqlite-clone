package instructions

import "github/com/lucasbn/sqlite-clone/app/machine/state"

type Instruction interface {
	Execute(s state.MachineState) state.MachineState
}
