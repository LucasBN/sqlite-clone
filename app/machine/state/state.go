package state

import "github/com/lucasbn/sqlite-clone/app/machine/registers"

type MachineState struct {
	CurrentAddress int
	Registers      registers.RegisterFile
	Halted         bool
}

func Init() *MachineState {
	return &MachineState{
		CurrentAddress: 0,
		Registers:      registers.Init(),
		Halted:         false,
	}
}
