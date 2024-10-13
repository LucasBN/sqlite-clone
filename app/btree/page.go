package btree

import (
	"bytes"
	"encoding/binary"
)

const intIdxPage = 2
const intTabPage = 5
const leafIdxPage = 10
const leafTabPage = 13

type bTreeHeader struct {
	PageType                uint8
	FirstFreeBlock          uint16
	NumCells                uint16
	CellContentOffset       uint16
	NumFragmenttedFreeBytes uint8
	RightMostPointer        uint32
}

func getPageHeader(page []byte, pageNum uint64) (bTreeHeader, error) {
	pointer := 0
	if pageNum == 1 {
		pointer += 100
	}

	var header bTreeHeader
	err := binary.Read(bytes.NewBuffer(page[pointer:pointer+12]), binary.BigEndian, &header)
	if err != nil {
		return bTreeHeader{}, err
	}
	return header, nil
}

// --------------------------------------------

type page struct {
	PageNumber uint64
	Data       []byte
}

func (p page) PageType() uint8 {
	if p.PageNumber == 1 {
		return p.Data[100]
	}
	return p.Data[0]
}

func (p page) NumCells() (uint16, error) {
	pointer := 0
	if p.PageNumber == 1 {
		pointer += 100
	}

	var numCells uint16
	err := binary.Read(bytes.NewBuffer(p.Data[pointer+3:pointer+4]), binary.BigEndian, &numCells)
	if err != nil {
		return 0, err
	}
	return numCells, nil
}

func (p page) RightMostPointer() (*uint32, error) {
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

// ReadInteriorTableCell interprets the bytes starting at the given offset as an
// interior table cell.
func (p page) ReadInteriorTableCell(offset uint64) (interiorTableCell, error) {
	return p.Data[offset:], nil
}
