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


// Read the page