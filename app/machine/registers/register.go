package registers

type RegisterFile[T any] struct {
	Value map[int]T
}

func Init[T any]() RegisterFile[T] {
	return RegisterFile[T]{Value: make(map[int]T)}
}

func (r RegisterFile[T]) Get(register int) T {
	return r.Value[register]
}

func (r RegisterFile[T]) Set(register int, value T) RegisterFile[T] {
	r.Value[register] = value
	return r
}
