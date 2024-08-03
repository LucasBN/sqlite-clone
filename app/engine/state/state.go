package state

import "github/com/lucasbn/sqlite-clone/app/engine/registers"

type MachineState struct {
	CurrentAddress int
	Registers      registers.Registers
	Halted         bool
}
