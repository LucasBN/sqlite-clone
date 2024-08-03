package main

type Machine struct {
	CurrentAddress int
	Program        []Instruction
	Registers      map[int]int
	Output         [][]int
	Halted         bool
}

type BetterInstruction interface {
	Execute()
}



// Maybe make this an interface with an execute function?
type Instruction struct {
	Opcode string
	P1     int
	P2     int
	P3     int
	P4     int
	P5     int
}

func NewMachine(instructions []Instruction) *Machine {
	return &Machine{
		CurrentAddress: 0,
		Program:        instructions,
		Registers:      make(map[int]int),
		Halted:         false,
	}
}

func (m *Machine) Run() [][]int {
	for {
		if len(m.Program) <= m.CurrentAddress {
			panic("Unreachable: attemping to run instruction at invalid address")
		}

		instruction := m.Program[m.CurrentAddress]
		m.Execute(instruction)
		if m.Halted {
			break
		}
	}
	return m.Output
}

func (m *Machine) Execute(instruction Instruction) {
	switch instruction.Opcode {
	case "Integer":
		m.ExecuteInteger(instruction)
	case "ResultRow":
		m.ExecuteResultRow(instruction)
	case "Halt":
		m.ExecuteHalt(instruction)
	default:
		break
	}
}

// Integer puts the value in p1 in the register specified by p2
func (m *Machine) ExecuteInteger(instruction Instruction) {
	m.Registers[instruction.P2] = instruction.P1
	m.CurrentAddress += 1
}

// Sets the output to be the value of registers P1 through to and including P2
func (m *Machine) ExecuteResultRow(instruction Instruction) {
	var result []int
	for i := instruction.P1; i <= instruction.P2; i++ {
		result = append(result, m.Registers[i])
	}
	m.Output = append(m.Output, result)
	m.CurrentAddress += 1
}

func (m *Machine) ExecuteHalt(_ Instruction) {
	m.Halted = true
	m.CurrentAddress += 1
}
