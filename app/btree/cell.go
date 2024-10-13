package btree

import (
	"bytes"
	"encoding/binary"
)

type interiorTableCell []byte

func (c interiorTableCell) LeftChild() (uint64, error) {
	var leftChild uint32
	err := binary.Read(bytes.NewBuffer(c[:4]), binary.BigEndian, &leftChild)
	if err != nil {
		return 0, err
	}
	return uint64(leftChild), nil
}
