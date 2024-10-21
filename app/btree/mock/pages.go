package mock

import (
	"encoding/binary"
	"github/com/lucasbn/sqlite-clone/app/types"

	"github.com/samber/lo"
)

type MockDatabaseHeader struct {
	HeaderString     [16]byte
	PageSize         uint16
	FileWriteVersion uint8
	FileReadVersion  uint8
	ReservedSpace    uint8
	Middle           [38]byte
	TextEncoding     uint32
	End              [40]byte
}

type MockPageHeader struct {
	PageType               uint8
	FirstFreeBlock         uint16
	NumCells               uint16
	CellContentOffset      uint16
	NumFragmentedFreeBytes uint8
	RightMostPointer       *uint32
}

func (h MockPageHeader) Serialize() []byte {
	buf := make([]byte, 12)
	buf[0] = h.PageType
	binary.BigEndian.PutUint16(buf[1:3], h.FirstFreeBlock)
	binary.BigEndian.PutUint16(buf[3:5], h.NumCells)
	binary.BigEndian.PutUint16(buf[5:7], h.CellContentOffset)
	buf[7] = h.NumFragmentedFreeBytes
	if h.RightMostPointer != nil {
		binary.BigEndian.PutUint32(buf[8:12], *h.RightMostPointer)
		return buf
	}
	return buf[:8]
}

type MockInteriorTableCell struct {
	LeftChildPageNumber uint32
	Key                 uint64
}

func (c MockInteriorTableCell) Serialize() []byte {
	buf := make([]byte, 13)
	binary.BigEndian.PutUint32(buf[0:4], c.LeftChildPageNumber)
	keySize := binary.PutUvarint(buf[4:13], c.Key)
	return buf[:4+keySize]
}

type MockInteriorTablePage struct {
	Header MockPageHeader
	Cells  []MockInteriorTableCell
}

func (p MockInteriorTablePage) Serialize() []byte {
	// The header always goes first
	buf := make([]byte, pageSize)

	pointer := 0

	// Store the header in the buffer and increment the pointer
	header := p.Header.Serialize()
	pointer += copy(buf[pointer:pointer+len(header)], header)

	// First store all of the cell bytes
	cells := [][]byte{}
	for _, cell := range p.Cells {
		cells = append(cells, cell.Serialize())
	}

	cellPointer := pageSize - len(lo.Flatten(cells))
	for _, cell := range cells {
		// Add a pointer to the cell
		binary.BigEndian.PutUint16(buf[pointer:pointer+2], uint16(cellPointer))
		pointer += 2

		// Write the contents of the cell to the page
		cellPointer += copy(buf[cellPointer:cellPointer+len(cell)], cell)
	}

	return buf
}

type MockLeafTableCell struct {
	Key     uint64
	Entries []types.Entry
}

func (c MockLeafTableCell) Serialize() []byte {
	// Key (varint)
	key := make([]byte, binary.MaxVarintLen64)
	keyVarintSize := binary.PutUvarint(key, c.Key)
	key = key[:keyVarintSize]

	recordHeaderSize := 0

	serialTypes := []byte{}
	values := []byte{}

	for _, entry := range c.Entries {
		switch entry := entry.(type) {
		case types.NullEntry:
			serialTypes = append(serialTypes, 0)
			recordHeaderSize += 1
		case types.NumberEntry:
			// Store all numbers as big-endian 64-bit twos-complement integer
			value := make([]byte, 8)
			binary.BigEndian.PutUint64(value, uint64(entry.Value))

			values = append(values, value...)
			serialTypes = append(serialTypes, 6)
			recordHeaderSize += 1
		case types.TextEntry:
			value := []byte(entry.Value)
			size := uint64((len(value) * 2) + 13)

			serialType := make([]byte, binary.MaxVarintLen64)
			serialTypeSize := binary.PutUvarint(serialType, size)
			serialType = serialType[:serialTypeSize]

			values = append(values, value...)
			serialTypes = append(serialTypes, serialType...)
			recordHeaderSize += serialTypeSize
		}
	}

	// recordHeaderSize now includes the size of the serial types, but we also
	// need to include the number of bytes that storing the header size as a
	// varint will take up
	sizeVarintBuf := make([]byte, binary.MaxVarintLen64)
	sizeVarintSize := binary.PutUvarint(sizeVarintBuf, uint64(recordHeaderSize))
	sizeVarintBuf = sizeVarintBuf[:sizeVarintSize]

	// Recalculate total header size by including the size of the size varint itself
	finalRecordHeaderSize := recordHeaderSize + sizeVarintSize

	// If the header size crosses into a new varint size range, recalculate
	for sizeVarintSize != binary.PutUvarint(sizeVarintBuf, uint64(finalRecordHeaderSize)) {
		sizeVarintSize = binary.PutUvarint(sizeVarintBuf, uint64(finalRecordHeaderSize))
		sizeVarintBuf = sizeVarintBuf[:sizeVarintSize]
		finalRecordHeaderSize = recordHeaderSize + sizeVarintSize
	}

	payloadSizeBuf := make([]byte, binary.MaxVarintLen64)
	payloadSizeVarintSize := binary.PutUvarint(payloadSizeBuf, uint64(finalRecordHeaderSize+len(values)))
	payloadSizeBuf = payloadSizeBuf[:payloadSizeVarintSize]

	result := []byte{}
	result = append(result, payloadSizeBuf...)
	result = append(result, key...)
	result = append(result, sizeVarintBuf...)
	result = append(result, serialTypes...)
	result = append(result, values...)

	return result
}

type MockLeafTablePage struct {
	Header MockPageHeader
	Cells  []MockLeafTableCell
}

func (p MockLeafTablePage) Serialize() []byte {
	// The header always goes first
	buf := make([]byte, pageSize)

	pointer := 0

	// Store the header in the buffer and increment the pointer
	header := p.Header.Serialize()
	pointer += copy(buf[pointer:pointer+len(header)], header)

	// First store all of the cell bytes
	cells := [][]byte{}
	for _, cell := range p.Cells {
		cells = append(cells, cell.Serialize())
	}

	cellPointer := pageSize - len(lo.Flatten(cells))
	for _, cell := range cells {
		// Add a pointer to the cell
		binary.BigEndian.PutUint16(buf[pointer:pointer+2], uint16(cellPointer))
		pointer += 2

		// Write the contents of the cell to the page
		cellPointer += copy(buf[cellPointer:cellPointer+len(cell)], cell)
	}

	return buf
}
