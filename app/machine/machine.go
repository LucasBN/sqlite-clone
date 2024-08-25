package machine

import (
	"github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/machine/instructions"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Machine struct {
	BTreeEngine btree.BTreeEngine
	State       *state.MachineState
	Program     []instructions.Instruction
	Output      [][]int
}

type MachineConfig struct {
	Instructions []instructions.Instruction
	BTreeEngine  btree.BTreeEngine
}

func Init(config MachineConfig) *Machine {
	return &Machine{
		BTreeEngine: config.BTreeEngine,
		Program:     config.Instructions,
		State:       state.Init(),
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
