package btree

import (
	"encoding/binary"
	"errors"
)

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

func (b *BTreeEngine[T]) constructLeafTableCellRecord(cell leafTableCell) (record[T], error) {
	cellPayload, err := cell.Payload()
	if err != nil {
		return record[T]{}, err
	}

	record := record[T]{
		ResultConstructor: b.resultConstructor,
		Data:              cellPayload,
	}

	return record, nil
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

func decodeUvarint(buf []byte) (uint64, uint64, error) {
	value, size := binary.Uvarint(buf)

	if size == 0 {
		return 0, 0, errors.New("uvarint: buffer too small")
	} else if size < 0 {
		return 0, 0, errors.New("uvarint: value does not fit in uint64")
	}

	return value, uint64(size), nil
}
