package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/samber/lo"
)

const INT_IDX_PAGE = 2
const INT_TAB_PAGE = 5
const LEAF_IDX_PAGE = 10
const LEAF_TAB_PAGE = 13

type BTreeHeader struct {
	PageType                uint8
	FirstFreeBlock          uint16
	NumCells                uint16
	CellContentOffset       uint16
	NumFragmenttedFreeBytes uint8
	RightMostPointer        uint32
}

type BTreePage struct {
	Header         BTreeHeader
	IntIdxCells    []IntIdxCell
	IntTableCells  []IntTableCell
	LeafIdxCells   []LeafIdxCell
	LeafTableCells []LeafTableCell
}

// pageNum is zero indexed which may be different to the SQLite standard
func readBTreePage(databaseFile *os.File, dbHeader DatabaseHeader, pageNum uint32) BTreePage {
	// Calculate the byte number in the file at which this page starts. The
	// first page in the database contains the database header, which is 100
	// bytes long
	pageStart := int64((pageNum - 1) * uint32(dbHeader.PageSize))

	// Seek to the beginning of the page
	databaseFile.Seek(lo.Ternary(pageNum == 1, pageStart+100, pageStart), io.SeekStart)

	// Read the first 12 bytes of the page into the header
	var header BTreeHeader
	if err := binary.Read(databaseFile, binary.BigEndian, &header); err != nil {
		fmt.Println("Failed to read integer:", err)
	}

	page := BTreePage{Header: header}

	// If the page isn't an interior b-tree page, we move back four bytes in the
	// file because the right most pointer isn't actually included in the header
	// for leaf pages.
	if header.PageType != INT_IDX_PAGE && header.PageType != INT_TAB_PAGE {
		databaseFile.Seek(-4, io.SeekCurrent)
	}

	// The cell pointer array immediately follows the header
	var cellPointerArray []uint16
	for i := 0; i < int(header.NumCells); i++ {
		var cellPointer uint16
		if err := binary.Read(databaseFile, binary.BigEndian, &cellPointer); err != nil {
			fmt.Println("Failed to read integer:", err)
		}
		cellPointerArray = append(cellPointerArray, cellPointer)
	}

	// Read the cells
	for idx, cellStart := range cellPointerArray {
		// The cellEnd of the cell is either the start of the next cell, or the
		// end of the page minus however many bytes are reserved at the end of
		// each page
		var cellEnd uint16
		if idx == 0 {
			cellEnd = dbHeader.PageSize - uint16(dbHeader.ReservedSpace)
		} else {
			cellEnd = cellPointerArray[idx-1]
		}

		// Seek to the beginning of the cell
		databaseFile.Seek(pageStart+int64(cellStart), io.SeekStart)

		cellBytes := make([]byte, cellEnd-cellStart)
		if err := binary.Read(databaseFile, binary.BigEndian, &cellBytes); err != nil {
			fmt.Println("Failed to read integer:", err)
		}

		// Appends the cell to the correct field in the page, depending on the
		// page type
		readCell(&page, dbHeader, cellBytes)
	}

	return page
}
