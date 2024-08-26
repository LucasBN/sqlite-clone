package btree

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/samber/lo"
)

type Pager interface {
	Close() error
	PageSize() uint64
	ReservedSpace() uint64
	GetPage(pageNum uint64) ([]byte, error)
}

type ValueType[T any] interface {
	Number(uint64) T
	Text(string) T
	Null() T
}

type BTreeEngine[T any] struct {
	Pager             Pager
	Cursors           map[uint64]*cursor
	ResultConstructor ValueType[T]
}

func NewBTreeEngine[T any](pager Pager, resultConstructor ValueType[T]) (*BTreeEngine[T], error) {
	return &BTreeEngine[T]{
		Pager:             pager,
		Cursors:           make(map[uint64]*cursor),
		ResultConstructor: resultConstructor,
	}, nil
}

// NewCursor creates a new cursor with the given ID that points to the table or
// index at the given root page number. The cursor is initialized to point at
// the very beginning of the root page.
func (b *BTreeEngine[T]) NewCursor(id uint64, rootPageNum uint64) (bool, error) {
	// Check if a cursor with the given ID already exists
	if _, ok := b.Cursors[id]; ok {
		return false, errors.New("cursor with the given ID already exists")
	}

	// Create a new cursor and add it to the map of cursors
	b.Cursors[id] = &cursor{
		ID:          id,
		Position:    0,
		CurrentPage: rootPageNum,
		RootPage:    rootPageNum,
	}

	return true, nil
}

// RewindCusor moves the cursor to the first entry in the database table or
// index.
func (b *BTreeEngine[T]) RewindCursor(id uint64) (bool, error) {
	// Get the cursor with the given ID
	cursor, err := b.getCursor(id)
	if err != nil {
		return false, err
	}

	// The first cell in the page is cell zero
	var firstCellNumber uint64 = 0

	// Get the position of the first cell in the root page
	newPosition, ok, err := b.getCellPointer(firstCellNumber, cursor.RootPage)
	if err != nil || !ok {
		return false, err
	}

	// Set the position of the cursor to be the byte offset of the first cell
	cursor.Position = *newPosition

	// Set the current cell of the cursor to be the first cell
	cursor.CurrentCell = &firstCellNumber

	return true, nil
}

// AdvanceCursor moves the cursor to the next leaf cell in the database table or
// index.
func (b *BTreeEngine[T]) AdvanceCursor(id uint64) (bool, error) {
	// Get the cursor with the given ID
	cursor, err := b.getCursor(id)
	if err != nil {
		return false, err
	}

	var newCellNumber uint64 = *cursor.CurrentCell + 1

	// Get the position of the next cell in the cyrrent page
	newPosition, didAdvance, err := b.getCellPointer(newCellNumber, cursor.CurrentPage)
	if err != nil || !didAdvance {
		return false, err
	}

	// Set the position of the cursor to be the byte offset of the first cell
	cursor.Position = *newPosition

	// Set the current cell of the cursor to be the first cell
	cursor.CurrentCell = &newCellNumber

	return true, nil
}

func (b *BTreeEngine[T]) ReadColumn(id uint64, column uint64) (T, error) {
	// Get the cursor with the given ID
	cursor, err := b.getCursor(id)
	if err != nil {
		return b.ResultConstructor.Null(), err
	}

	// Get the page of the table or index
	page, err := b.Pager.GetPage(cursor.CurrentPage)
	if err != nil {
		return b.ResultConstructor.Null(), err
	}

	// Get the byte offset at which the cell ends
	var cellEnd uint64
	if *cursor.CurrentCell == 0 {
		// If the current cell is the first cell in the page, the cell end is
		// the end of the page minus the reserved space
		cellEnd = b.Pager.PageSize() - b.Pager.ReservedSpace()
	} else {
		// Otherwise, the cell end is the start of the previous cell
		nextCell, ok, err := b.getCellPointer(*cursor.CurrentCell-1, cursor.CurrentPage)
		if err != nil || !ok {
			return b.ResultConstructor.Null(), err
		}
		cellEnd = *nextCell
	}

	// Read the cell data from the page
	cell := page[cursor.Position:cellEnd]

	// We'll start at the beginning of the cell and keep a pointer to keep track
	// of our position within the cell
	var pointer uint64 = 0

	// Read payload size from the cell - a varint could be up to 9 bytes long,
	// we'll need to read at most 9 bytes
	payloadSize, offset := decodeVarInt(cell[pointer:lo.Min([]int{len(cell) - 1, 9})])
	pointer += uint64(offset)

	// There's some somewhat complicated logic to deal with the payload
	// overflowing onto another page. I'm not going to implement it here and
	// instead panic if we encounter this case
	if payloadSize > uint64(b.Pager.PageSize()-b.Pager.ReservedSpace()) {
		return b.ResultConstructor.Null(), errors.New("unsupported: payload overflows onto another page")
	}

	// Read the row ID from the cell
	_, offset = decodeVarInt(cell[pointer:lo.Min([]int{len(cell) - int(offset) - 1, 9})])
	pointer += uint64(offset)

	// Read the record header size
	recordHeaderSize, offset := decodeVarInt(cell[pointer:])
	pointer += uint64(offset)

	// Read the type codes from the record header
	var typeCodes []uint64
	{
		// Read the record header
		recordHeaderEnd := pointer + recordHeaderSize - uint64(offset)

		for {
			typeCode, offset := decodeVarInt(cell[pointer:])
			typeCodes = append(typeCodes, typeCode)

			pointer += uint64(offset)
			if pointer >= recordHeaderEnd {
				break
			}
		}
	}

	// Read the column from the record data
	intTypeCodeToSize := map[uint64]uint64{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 6,
		6: 8,
	}

	for idx, typeCode := range typeCodes {
		var entry T
		if typeCode == 0 {
			entry = b.ResultConstructor.Null()
		} else if typeCode >= 1 && typeCode <= 6 {
			// Make an 8 byte empty slice
			result := make([]byte, 8)

			// Extract the correct number of bytes from the raw record
			size := intTypeCodeToSize[typeCode]
			value := cell[pointer : pointer+size]
			pointer += size

			// Copy the bytes into the result slice so that we can decode them
			// as a uint64
			for i := 7; 7-i < len(value); i-- {
				result[i] = value[7-i]
			}
			entry = b.ResultConstructor.Number(binary.BigEndian.Uint64(result))
		} else if typeCode == 8 {
			entry = b.ResultConstructor.Number(0)
		} else if typeCode == 9 {
			entry = b.ResultConstructor.Number(1)
		} else if typeCode >= 12 && typeCode%2 == 1 {
			length := (typeCode - 12) / 2
			entry = b.ResultConstructor.Text(string(cell[pointer : pointer+length]))
			pointer += length
		} else {
			return b.ResultConstructor.Null(), errors.New("unsupported type code: not implemented")
		}

		if uint64(idx) == column {
			// This is the entry we want, so read
			return entry, nil
		}
	}

	// If the column wasn't found, return a null entry
	return b.ResultConstructor.Null(), nil
}

func (b *BTreeEngine[T]) getCursor(id uint64) (*cursor, error) {
	cursor, ok := b.Cursors[id]
	if !ok {
		return nil, errors.New("cursor with the given ID does not exist")
	}

	return cursor, nil
}

func (b *BTreeEngine[T]) getCellPointer(cellNum uint64, pageNum uint64) (*uint64, bool, error) {
	// Get the page of the table or index
	page, err := b.Pager.GetPage(pageNum)
	if err != nil {
		return nil, false, err
	}

	// We'll start at the beginning of the page and read the BTreeHeader to find
	// the position of the first leaf cell
	var pointer uint64 = 0

	// If we're on the very first page, we need to skip the database header
	if pageNum == 1 {
		pointer += 100
	}

	// Read the BTreeHeader
	var header bTreeHeader
	err = binary.Read(bytes.NewBuffer(page[pointer:pointer+12]), binary.BigEndian, &header)
	if err != nil {
		return nil, false, err
	}

	// Check if the cell number is out of bounds and return false as the 'ok' parameter
	if uint64(header.NumCells) <= cellNum {
		return nil, false, nil
	}

	// Advance the pointer 12 bytes to skip over the BTreeHeader
	pointer += 12

	// We only support leaf table pages for now, so return an error if the page
	// is anything else
	if header.PageType != LEAF_TAB_PAGE {
		return nil, false, errors.New("cursor only supports leaf table pages")
	}

	// Move back 4 bytes if we're not in an interior page because we don't store
	// the right most pointer in the header for leaf pages
	if header.PageType != INT_TAB_PAGE && header.PageType != INT_IDX_PAGE {
		pointer -= 4
	}

	// Get the first entry of the cell pointer array, which starts immediately
	// after the B-Tree header and consists of 2-byte unsigned integers
	var cellPointer uint16
	cellPointerStart := pointer + (2 * cellNum)
	err = binary.Read(bytes.NewBuffer(page[cellPointerStart:cellPointerStart+2]), binary.BigEndian, &cellPointer)
	if err != nil {
		return nil, false, err
	}

	// Cast the cell pointer to a uint64
	var cellPointer64 uint64 = uint64(cellPointer)

	return &cellPointer64, true, nil
}

// A varint consists of either zero or more bytes which have the high-order
// bit set followed by a single byte with the high-order bit clear, or nine
// bytes, whichever is shorter.
func decodeVarInt(data []byte) (uint64, uint16) {
	var value uint64
	for i := 0; i < 8; i++ {
		value = (value << 7) | uint64(data[i]&0x7F)
		if data[i]&0x80 == 0 {
			return value, uint16(i + 1)
		}
	}
	return value<<8 | uint64(data[8]), 9
}
