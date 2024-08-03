package instructions

import "github/com/lucasbn/sqlite-clone/app/engine/state"

type Halt struct{}

func (Halt) Execute(s state.MachineState) state.MachineState {
	return state.MachineState{
		CurrentAddress: s.CurrentAddress,
		Registers:      s.Registers,
		Halted:         true,
	}
}
