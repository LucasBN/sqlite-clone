package btree

import "github/com/lucasbn/sqlite-clone/app/pager"

type BTreeProcessor struct {
	Pager       *pager.Pager
	Cursors     map[uint32]*Cursor
	CursorCount uint32
}

func Init(dbFilePath string) *BTreeProcessor {
	return &BTreeProcessor{
		Pager:   pager.Init(dbFilePath),
		Cursors: make(map[uint32]*Cursor),
	}
}

func (b *BTreeProcessor) NewCursor(id uint32, rootPage uint32) *Cursor {
	// Check if a cursor with the given ID already exists
	if _, ok := b.Cursors[id]; ok {
		panic("Cursor with the given ID already exists")
	}

	// It doesn't exist, create a new cursor
	newCursor := InitCursor(b.CursorCount, rootPage, b.Pager)

	// Add the new cursor to the map of cursors
	b.Cursors[newCursor.ID] = newCursor

	// Increment the cursor count
	b.CursorCount++

	// Return the cursor
	return newCursor
}

func (b *BTreeProcessor) GetCursor(id uint32) *Cursor {
	// Check if a cursor with the given ID exists
	if cursor, ok := b.Cursors[id]; ok {
		return cursor
	}

	// It doesn't exist, panic
	panic("Cursor with the given ID doesn't exist")
}
