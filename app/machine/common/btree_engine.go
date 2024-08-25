package common

import "github/com/lucasbn/sqlite-clone/app/types"

type BTreeEngine interface {
	NewCursor(id uint64, rootPageNum uint64) (bool, error)
	RewindCursor(id uint64) (bool, error)
	AdvanceCursor(id uint64) (bool, error)
	ReadColumn(id uint64, column uint64) (types.Entry, error)
}
