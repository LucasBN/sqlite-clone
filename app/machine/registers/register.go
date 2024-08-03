package registers

type Registers struct {
	Value map[int]int
}

func Init() Registers {
	return Registers{Value: make(map[int]int)}
}

func (r Registers) Get(register int) int {
	return r.Value[register]
}

func (r Registers) Set(register int, value int) Registers {
	r.Value[register] = value
	return r
}
