package btree

import (
	"errors"

	"github.com/samber/lo"
)

// NewCursor creates a new cursor with the given ID that points to the table or
// index at the given root page number. The cursor is initialized to point at
// the very beginning of the root page.
func (b *BTreeEngine[T]) NewCursor(id uint64, rootPageNum uint64) (bool, error) {
	// Check if a cursor with the given ID already exists
	if _, ok := b.cursors[id]; ok {
		return false, errors.New("cursor with the given ID already exists")
	}

	// Create a new cursor and add it to the map of cursors
	b.cursors[id] = &cursor{
		ID: id,
		PagePositionStack: []pagePosition{
			{
				ByteOffset: 0,
				PageNumber: rootPageNum,
			},
		},
		RootPage: rootPageNum,
	}

	return true, nil
}

// RewindCursor moves the cursor to the first leaf cell in the table by
// traversing the B-Tree (always taking the left most branch at each interior
// page, and setting the cursor on every visited page to point at the first cell
// on that page).
func (b *BTreeEngine[T]) RewindCursor(id uint64) (bool, error) {
	// Get the cursor with the given ID
	cursor, err := b.getCursor(id)
	if err != nil {
		return false, err
	}

	// Get the root page
	page, err := b.getPage(cursor.RootPage)
	if err != nil {
		return false, err
	}

	// If the root page is not a leaf page, we need to traverse the tree by
	// taking the left most branch at each interior page until we reach a leaf
	// page
	err = b.moveCursorToLeftMostLeafPage(page, cursor)
	if err != nil {
		return false, err
	}

	// Get the page left most leaf page
	page, err = b.getPage(cursor.CurrentPage())
	if err != nil {
		return false, err
	}

	return cursor.moveToCell(page, 0)
}

// AdvanceCursor moves the cursor to the next leaf cell in table. If there is
// a next cell on the current page (which should be a leaf page), then the
// cursor is simply moved to that page. However, if we are at the last cell on a
// leaf page, then we must use the cursor position stack to find the next leaf
// page to visit. If the cursor already points to the very last cell in the
// B-Tree, then this function will return false.
func (b *BTreeEngine[T]) AdvanceCursor(id uint64) (bool, error) {
	// Get the cursor with the given ID
	cursor, err := b.getCursor(id)
	if err != nil {
		return false, err
	}

	// If there are no pages on the position stack, then we've reached the end
	// of the tree
	if len(cursor.PagePositionStack) == 0 {
		return false, nil
	}

	page, err := b.getPage(cursor.CurrentPage())
	if err != nil {
		return false, err
	}

	// AdvanceCursor should only be called when the cursor is on a leaf page, so
	// that we can use this assumption to calculate how to get to the next cell
	if page.PageType() != leafTabPage {
		return false, errors.New("expected page to be a leaf table page")
	}

	// If there's another cell on this leaf page, just go to that
	if numCells, err := page.NumCells(); err != nil {
		return false, err
	} else if uint64(numCells) > *cursor.CurrentCell()+1 {
		return cursor.moveToCell(page, *cursor.CurrentCell()+1)
	}

	// There are no more cells on this page, so we need to remove it from the
	// page position stack (so that we can then find the next cell)
	cursor.PagePositionStack = cursor.PagePositionStack[:len(cursor.PagePositionStack)-1]

	// There are no parent pages and we have processed the last cell
	if len(cursor.PagePositionStack) == 0 {
		return false, nil
	}

	// The position stack now consists of at least one interior table page. It
	// could be the case that we've already read the leaf page pointed to by the
	// right most pointer of the current interior page (the page on the top of
	// the position stack) which means we need to continuously pop pages off the
	// stack until we find a page that has a cell we haven't visited yet.
	for {
		// Fetch the interior page at the top of the position stack
		page, err = b.getPage(cursor.CurrentPage())
		if err != nil {
			return false, err
		}

		// Since the right most pointer isn't actually stored in a cell, we'll
		// compare the cell we're currently at with the number of cells on the
		// page and if we're the last page we'll (i) mark the page pointed to by
		// the right most pointer as the next page to visit and (ii) pop the
		// current page off the position stack stack so that we don't visit it
		// again.
		numCells, err := page.NumCells()
		if err != nil {
			return false, err
		}

		var nextPage *uint64
		if *cursor.CurrentCell()+1 < uint64(numCells) {
			// If we're in this branch, then there's another cell on the current
			// page that we need to visit, so let's move to it and then mark the
			// left child as the next page to visit
			didAdvance, err := cursor.moveToCell(page, *cursor.CurrentCell()+1)
			if err != nil {
				return false, err
			}
			if !didAdvance {
				return false, errors.New("expected cursor to be able to advance")
			}

			cell, err := page.ReadInteriorTableCell(cursor.Position())
			if err != nil {
				return false, err
			}

			leftChild, err := cell.LeftChild()
			if err != nil {
				return false, err
			}

			nextPage = &leftChild
		} else if *cursor.CurrentCell()+1 == uint64(numCells) {
			// If we're in this branch, then we've visited every child page of
			// the current interior page, so we should drop it from the position
			// stack and mark the right most pointer as the next page to visit

			// Pop the current page off the stack
			cursor.PagePositionStack = cursor.PagePositionStack[:len(cursor.PagePositionStack)-1]

			// Move to the right pointer
			rightPointer, err := page.RightMostPointer()
			if err != nil {
				return false, err
			}

			nextPage = lo.ToPtr(uint64(*rightPointer))
		}

		// If we've found a next page to visit, add it to the position stack and
		// break out of the loop
		if nextPage != nil {
			cursor.PagePositionStack = append(cursor.PagePositionStack, pagePosition{
				ByteOffset: 0,
				PageNumber: *nextPage,
			})

			page, err = b.getPage(*nextPage)
			if err != nil {
				return false, err
			}

			break
		}

		// Otherwise, we've already visited every child page of this interior
		// page, so we should drop it from the position stack and continue
		cursor.PagePositionStack = cursor.PagePositionStack[:len(cursor.PagePositionStack)-1]

		// If the stack is now empty, we've reached the very last cell of the
		// tree
		if len(cursor.PagePositionStack) == 0 {
			return false, nil
		}
	}

	err = b.moveCursorToLeftMostLeafPage(page, cursor)
	if err != nil {
		return false, err
	}

	page, err = b.getPage(cursor.CurrentPage())
	if err != nil {
		return false, err
	}

	return cursor.moveToCell(page, 0)
}

// ReadColumn reads the column at the given index from the current cell that the
// cursor is pointing to.
func (b *BTreeEngine[T]) ReadColumn(id uint64, column uint64) (T, error) {
	// Get the cursor with the given ID
	cursor, err := b.getCursor(id)
	if err != nil {
		return b.resultConstructor.Null(), err
	}

	// Get the page of the table or index
	page, err := b.getPage(cursor.CurrentPage())
	if err != nil {
		return b.resultConstructor.Null(), err
	}

	// Read the cell data that the cursor is pointing at
	cell, err := page.ReadLeafTableCell(cursor.Position())
	if err != nil {
		return b.resultConstructor.Null(), err
	}

	// Construct a record from the cell
	record, err := b.constructLeafTableCellRecord(cell)
	if err != nil {
		return b.resultConstructor.Null(), err
	}

	return record.ReadColumn(column)
}
