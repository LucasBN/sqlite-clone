package common

type BTreeEngine[T any] interface {
	NewCursor(id uint64, rootPageNum uint64) (bool, error)
	RewindCursor(id uint64) (bool, error)
	AdvanceCursor(id uint64) (bool, error)
	ReadColumn(id uint64, column uint64) (T, error)
}
