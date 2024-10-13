package btree

// A cursor points to a specific entry in a b-tree, which means that it points
// to a specific byte offset in the database.
//
// Currently, a cursor makes a very incorrect assumption that every page is a
// leaf table page (no indexes, no interior pages). This means that we only need
// to store the absolute byte offset within a database file that the cursor is
// pointing to.
//
// Adding support for interior pages might require us to store more information,
// as we'll probably need a way to jump from one page to another.
//
// Cursors also assume that the caller 'knows' what they're doing, and therefore
// do not try to protect against 'invalid' operations. For example, if the
// caller attempts to call ReadColumn on a cursor that isn't actually pointing
// to a valid record, the cursor will read the bytes at the current position and
// interpret them as a record (and get the column data from it). However, errors
// may still occur if the cursor attempts, for example, to read past the end of
// the page.

type pagePosition struct {
	ByteOffset uint64
	CellNumber *uint64
	PageNumber uint64
}

type cursor struct {
	// The ID of the cursor
	ID uint64

	PagePositionStack []pagePosition

	// The page number of the root page of the B-Tree
	RootPage uint64
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
