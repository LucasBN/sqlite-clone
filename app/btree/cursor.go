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
type cursor struct {
	// The ID of the cursor
	ID uint64

	// The byte offset of the cursor on the current page
	Position uint64

	// The cell number of the cell that the cursor is currently pointing to
	CurrentCell *uint64

	// The page number of the page that the cursor is currently pointing to
	CurrentPage uint64

	// The page number of the root page of the B-Tree
	RootPage uint64
}

const INT_IDX_PAGE = 2
const INT_TAB_PAGE = 5
const LEAF_IDX_PAGE = 10
const LEAF_TAB_PAGE = 13

type bTreeHeader struct {
	PageType                uint8
	FirstFreeBlock          uint16
	NumCells                uint16
	CellContentOffset       uint16
	NumFragmenttedFreeBytes uint8
	RightMostPointer        uint32
}
