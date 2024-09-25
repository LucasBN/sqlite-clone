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

func getPageHeader(page []byte) (bTreeHeader, error) {
	var header bTreeHeader
	err := binary.Read(bytes.NewBuffer(page[0:12]), binary.BigEndian, &header)
	if err != nil {
		return bTreeHeader{}, err
	}
	return header, nil
}
