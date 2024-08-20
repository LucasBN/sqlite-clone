package machine

import (
	"github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/machine/instructions"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Machine struct {
	BTreeProcessor *btree.BTreeProcessor
	State          *state.MachineState
	Program        []instructions.Instruction
	Output         [][]int
}

type MachineConfig struct {
	DBFilePath   string
	Instructions []instructions.Instruction
}

func Init(config MachineConfig) *Machine {
	return &Machine{
		BTreeProcessor: btree.Init(config.DBFilePath),
		State:          state.Init(),
		Program:        config.Instructions,
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
		out := instruction.Execute(m.State)

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
