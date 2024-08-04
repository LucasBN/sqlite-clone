package registers

type RegisterFile struct {
	Value map[int]int
}

func Init() RegisterFile {
	return RegisterFile{Value: make(map[int]int)}
}

func (r RegisterFile) Get(register int) int {
	return r.Value[register]
}

func (r RegisterFile) Set(register int, value int) RegisterFile {
	r.Value[register] = value
	return r
}
