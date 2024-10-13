package btree

import "errors"

func (b *BTreeEngine[T]) getPage(pageNum uint64) (btreePage, error) {
	p, err := b.pager.GetPage(pageNum)
	if err != nil {
		return btreePage{}, err
	}

	return btreePage{
		PageNumber:    pageNum,
		ReservedSpace: b.pager.ReservedSpace(),
		PageSize:      b.pager.PageSize(),
		Data:          p,
	}, nil
}

func (b *BTreeEngine[T]) getCursor(id uint64) (*cursor, error) {
	cursor, ok := b.cursors[id]
	if !ok {
		return nil, errors.New("cursor with the given ID does not exist")
	}

	return cursor, nil
}

func (b *BTreeEngine[T]) moveCursorToLeftMostLeafPage(p btreePage, c *cursor) error {
	if p.PageType() == leafTabPage {
		return nil
	}

	ok, err := c.moveToCell(p, 0)
	if err != nil || !ok {
		return err
	}

	cell, err := p.ReadInteriorTableCell(c.Position())
	if err != nil {
		return err
	}

	leftChild, err := cell.LeftChild()
	if err != nil {
		return err
	}

	// Update the cursor position stack
	c.PagePositionStack = append(c.PagePositionStack, pagePosition{
		ByteOffset: 0,
		PageNumber: leftChild,
	})

	// Get the left page
	nextPage, err := b.getPage(leftChild)
	if err != nil {
		return err
	}

	return b.moveCursorToLeftMostLeafPage(nextPage, c)
}
