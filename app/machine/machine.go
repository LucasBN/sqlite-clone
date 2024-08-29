package machine

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/instructions"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Machine[T any] struct {
	BTreeEngine common.BTreeEngine[T]
	State       *state.MachineState[T]
	Program     []instructions.Instruction[T]
	Output      [][]T
}

type MachineConfig[T any] struct {
	Instructions []instructions.Instruction[T]
	BTreeEngine  common.BTreeEngine[T]
}

func NewMachine[T any](config MachineConfig[T]) *Machine[T] {
	return &Machine[T]{
		BTreeEngine: config.BTreeEngine,
		Program:     config.Instructions,
		State:       state.Init[T](),
	}
}

func (m *Machine[T]) Run() [][]T {
	for {
		if len(m.Program) <= m.State.CurrentAddress {
			panic("Unreachable: attemping to run instruction at invalid address")
		}

		// Fetch the instruction
		instruction := m.Program[m.State.CurrentAddress]

		// Execute the instruction and update the machine state
		out := instruction.Execute(m.State, m.BTreeEngine)

		// Append the output of the instruction to the machine output
		m.Output = append(m.Output, out...)

		// Stop the execution of the machine if the instruction resulted in a
		// halt
		if m.State.Halted {
			break
		}
	}
	return m.Output
}
