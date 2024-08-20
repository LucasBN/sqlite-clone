package btree

import "github/com/lucasbn/sqlite-clone/app/pager"

type BTreeProcessor struct {
	Pager   *pager.Pager
	Cursors map[uint32]*Cursor
}

func Init(_ string) *BTreeProcessor {
	return &BTreeProcessor{
		Pager:   pager.Init(),
		Cursors: make(map[uint32]*Cursor),
	}
}

func (b *BTreeProcessor) GetReadOnlyCursor(_ uint32, _ uint32) *Cursor {
	return &Cursor{}
}
