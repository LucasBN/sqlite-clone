package state

import "github/com/lucasbn/sqlite-clone/app/machine/registers"

type MachineState[T any] struct {
	CurrentAddress int
	Registers      registers.RegisterFile[T]
	Halted         bool
}

func Init[T any]() *MachineState[T] {
	return &MachineState[T]{
		CurrentAddress: 0,
		Registers:      registers.Init[T](),
		Halted:         false,
	}
}
