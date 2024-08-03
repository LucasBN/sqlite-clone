package instructions

import "github/com/lucasbn/sqlite-clone/app/machine/state"

type Integer struct {
	Register int
	Value    int
}

func (integer Integer) Execute(s state.MachineState) state.MachineState {
	return state.MachineState{
		CurrentAddress: s.CurrentAddress + 1,
		Registers:      s.Registers.Set(integer.Register, integer.Value),
		Halted:         s.Halted,
	}
}
