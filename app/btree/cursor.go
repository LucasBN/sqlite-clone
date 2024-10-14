package btree

// A cursor is an internal structure that points at a particular entry in a
// B-Tree. The PagePositionStack keeps track of how the cursor got to the
// current page, so that it can back track and (for example) move from the last
// entry on one page to the first entry on the other page. The last element in
// the stack is the current page and entry that the cursor is pointing at.
type cursor struct {
	// The ID of the cursor
	ID uint64

	PagePositionStack []pagePosition

	// The page number of the root page of the B-Tree
	RootPage uint64
}

// A pagePosition points to a particular byte offset in a page. The cell number
// is a pointer because it the cursor may not be pointing at a cell (but if it
// is, then the CellNumber *should* be set and if it isn't then it should be
// nil).
type pagePosition struct {
	ByteOffset uint64
	CellNumber *uint64
	PageNumber uint64
}

func (cursor *cursor) CurrentPage() uint64 {
	return cursor.PagePositionStack[len(cursor.PagePositionStack)-1].PageNumber
}

func (cursor *cursor) CurrentCell() *uint64 {
	return cursor.PagePositionStack[len(cursor.PagePositionStack)-1].CellNumber
}

func (cursor *cursor) Position() uint64 {
	return cursor.PagePositionStack[len(cursor.PagePositionStack)-1].ByteOffset
}

func (cursor *cursor) SetPosition(position uint64) {
	cursor.PagePositionStack[len(cursor.PagePositionStack)-1].ByteOffset = position
}

func (cursor *cursor) SetCellNumber(cellNumber uint64) {
	cursor.PagePositionStack[len(cursor.PagePositionStack)-1].CellNumber = &cellNumber
}

func (c *cursor) moveToCell(p btreePage, cellNumber uint64) (bool, error) {
	newPosition, err := p.CellPointer(cellNumber)
	if err != nil {
		return false, err
	}

	// Set the position of the cursor to be the byte offset of the first cell
	c.SetPosition(newPosition)

	// Set the current cell of the cursor to be the first cell
	c.SetCellNumber(cellNumber)

	return true, nil
}
