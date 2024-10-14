package btree

import (
	"bytes"
	"encoding/binary"
)

// -------------------------------------
// Interior Table Cells
// -------------------------------------
type interiorTableCell []byte

func (c interiorTableCell) LeftChild() (uint64, error) {
	var leftChild uint32
	err := binary.Read(bytes.NewBuffer(c[:4]), binary.BigEndian, &leftChild)
	if err != nil {
		return 0, err
	}
	return uint64(leftChild), nil
}

// -------------------------------------
// Leaf Table Cells
// -------------------------------------
type leafTableCell []byte

func (c leafTableCell) Payload() ([]byte, error) {
	pointer := uint64(0)

	// Detemine the payload size and the size of the varint that stores it
	_, payloadSizeVarintSize, err := decodeUvarint(c[pointer:])
	if err != nil {
		return nil, err
	}
	pointer += payloadSizeVarintSize

	// Determine the size of the rowid varint
	_, rowidVarintSize, err := decodeUvarint(c[pointer:])
	if err != nil {
		return nil, err
	}
	pointer += rowidVarintSize

	return c[pointer:], nil
}
