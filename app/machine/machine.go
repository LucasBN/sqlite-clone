package machine

import (
	"github/com/lucasbn/sqlite-clone/app/machine/instructions"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Machine struct {
	State   state.MachineState
	Program []instructions.Instruction
	Output  [][]int
}

func Init(instructions []instructions.Instruction) *Machine {
	return &Machine{
		State:   state.Init(),
		Program: instructions,
	}
}

func (m *Machine) Run() [][]int {
	for {
		if len(m.Program) <= m.State.CurrentAddress {
			panic("Unreachable: attemping to run instruction at invalid address")
		}

		// Fetch the instruction
		instruction := m.Program[m.State.CurrentAddress]

		// Execute the instruction and update the machine state
		m.State = instruction.Execute(m.State)

		// Stop the execution of the machine if the instruction resulted in a
		// halt
		if m.State.Halted {
			break
		}
	}
	return m.Output
}
