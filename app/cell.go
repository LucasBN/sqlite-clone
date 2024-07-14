package main

type Cell struct {
	CellType          uint8
	TableLeafCell     *TableLeafCell
	TableInteriorCell *TableInteriorCell
	IndexLeafCell     *IndexLeafCell
	IndexInteriorCell *IndexInteriorCell
}

// Idea: a cell struct shouldn't even be aware about the existance of pages, so
// the payload should be the entire payload and not just the payload that fits
// on the page where the cell begins
type TableLeafCell struct {
	PayloadSize  uint64
	RowID        uint64
	Payload      []byte
	OverflowPage *uint32
}

type TableInteriorCell struct {
	LeftPage uint32
	Key      uint64
}

type IndexLeafCell struct {
	KeyPayloadSize uint64
	Payload        []byte
	OverflowPage   *uint32
}

type IndexInteriorCell struct {
	LeftPage       uint32
	KeyPayloadSize uint64
	Payload        []byte
	Key            uint64
}

func readCell(header DatabaseHeader, data []byte, pageType uint8) Cell {
	switch pageType {
	case LEAF_TAB_PAGE:
		return Cell{CellType: pageType, TableLeafCell: readTableLeafCell(header, data)}
	case INT_TAB_PAGE:
		return Cell{CellType: pageType, TableInteriorCell: readTableInteriorCell(data)}
	case LEAF_IDX_PAGE:
		return Cell{CellType: pageType, IndexLeafCell: readIndexLeafCell(data)}
	case INT_IDX_PAGE:
		return Cell{CellType: pageType, IndexInteriorCell: readIndexInteriorCell(data)}
	default:
		panic("Unknown page type")
	}
}

func readTableLeafCell(header DatabaseHeader, data []byte) *TableLeafCell {
	var cell TableLeafCell

	var offsetPayloadSize uint16
	cell.PayloadSize, offsetPayloadSize = decodeVarInt(data[:9])

	var offsetRowID uint16
	cell.RowID, offsetRowID = decodeVarInt(data[offsetPayloadSize : offsetPayloadSize+9])

	// This doesn't take into account the case where the cell overflows onto the
	// next page (which means that the last four bytes of the cell are a pointer
	// to the overflow page) but if we encounter this case we'll end up panicing
	// just after this
	cell.Payload = data[offsetPayloadSize+offsetRowID:]

	// There's some somewhat complicated logic to deal with the payload
	// overflowing onto another page. I'm not going to implement it here and
	// instead panic if we encounter this case
	if cell.PayloadSize > uint64(header.PageSize-uint16(header.ReservedSpace)) {
		panic("Payload overflow")
	}

	return &cell
}

func readTableInteriorCell(data []byte) *TableInteriorCell {
	panic("not implemented")
}

func readIndexLeafCell(data []byte) *IndexLeafCell {
	panic("not implemented")
}

func readIndexInteriorCell(data []byte) *IndexInteriorCell {
	panic("not implemented")
}
