package btree

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const intIdxPage = 2
const intTabPage = 5

// const leafIdxPage = 10
const leafTabPage = 13

type btreePage struct {
	PageNumber    uint64
	ReservedSpace uint64
	PageSize      uint64
	Data          []byte
}

func (p btreePage) PageType() uint8 {
	if p.PageNumber == 1 {
		return p.Data[100]
	}
	return p.Data[0]
}

func (p btreePage) NumCells() (uint16, error) {
	pointer := 0
	if p.PageNumber == 1 {
		pointer += 100
	}

	var numCells uint16
	err := binary.Read(bytes.NewBuffer(p.Data[pointer+3:pointer+5]), binary.BigEndian, &numCells)
	if err != nil {
		return 0, err
	}
	return numCells, nil
}

func (p btreePage) RightMostPointer() (*uint32, error) {
	// We only have pointers to other pages on interior pages
	if p.PageType() != intIdxPage && p.PageType() != intTabPage {
		return nil, nil
	}

	pointer := 0
	if p.PageNumber == 1 {
		pointer += 100
	}

	var rightMostPointer uint32
	err := binary.Read(bytes.NewBuffer(p.Data[pointer+8:pointer+12]), binary.BigEndian, &rightMostPointer)
	if err != nil {
		return nil, err
	}
	return &rightMostPointer, nil
}

func (p btreePage) CellPointer(cellNum uint64) (uint64, error) {
	// We'll start at the beginning of the page and calculate the offset to the
	// cell pointer
	pointer := uint64(0)

	// If we're on the very first page, we need to skip the database header
	if p.PageNumber == 1 {
		pointer += 100
	}

	// Validate that the cellNum refers to an actual cell
	if cellCount, err := p.NumCells(); err != nil {
		return 0, err
	} else if uint64(cellCount) <= cellNum {
		return 0, errors.New("cellNum is greater than the number of cells on the page")
	}

	// Advance the pointer 12 bytes to skip over the BTreeHeader
	pointer += 12

	// Move back 4 bytes if we're not in an interior page because we don't store
	// the right most pointer in the header for leaf pages
	if p.PageType() != intTabPage && p.PageType() != intIdxPage {
		pointer -= 4
	}

	// The first cell pointer is at offset 12
	offset := pointer + cellNum*2
	var cellPointer uint16
	err := binary.Read(bytes.NewBuffer(p.Data[offset:offset+2]), binary.BigEndian, &cellPointer)
	if err != nil {
		return 0, err
	}
	return uint64(cellPointer), nil
}

// ReadInteriorTableCell interprets the bytes starting at the given offset as an
// interior table cell.
func (p btreePage) ReadInteriorTableCell(offset uint64) (interiorTableCell, error) {
	// Read the varint after the left child pointer so that we can determine the
	// total size of the cell
	_, bytesRead := binary.Uvarint(p.Data[offset+4:])

	return p.Data[offset : offset+uint64(bytesRead)+4], nil
}

func (p btreePage) ReadLeafTableCell(offset uint64) (leafTableCell, error) {
	cellEnd := offset

	// Detemine the payload size and the size of the varint that stores it
	payloadSize, payloadSizeVarintSize := binary.Uvarint(p.Data[offset:])
	cellEnd += payloadSize + uint64(payloadSizeVarintSize)

	// Determine the size of the rowid varint
	_, rowidVarintSize := binary.Uvarint(p.Data[offset+uint64(payloadSizeVarintSize):])
	cellEnd += uint64(rowidVarintSize)

	// For now, we're only going to support cells that have zero overflow pages
	if p.PageSize-p.ReservedSpace < cellEnd {
		return nil, errors.New("unsupported: cell has overflow pages")
	}

	return p.Data[offset:cellEnd], nil
}
