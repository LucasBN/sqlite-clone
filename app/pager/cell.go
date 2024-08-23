package pager

import "github.com/samber/lo"

type LeafTableCell struct {
	PayloadSize  uint64
	RowID        uint64
	Payload      []byte
	OverflowPage *uint32
}

type LeafIdxCell struct {
	LeftPage uint32
	Key      uint64
}

type IntTableCell struct {
	KeyPayloadSize uint64
	Payload        []byte
	OverflowPage   *uint32
}

type IntIdxCell struct {
	LeftPage       uint32
	KeyPayloadSize uint64
	Payload        []byte
	Key            uint64
}

func readCell(page *BTreePage, dbHeader DatabaseHeader, data []byte) {
	switch page.Header.PageType {
	case LEAF_TAB_PAGE:
		page.LeafTableCells = append(page.LeafTableCells, readLeafTableCell(dbHeader, data))
	case LEAF_IDX_PAGE:
		page.LeafIdxCells = append(page.LeafIdxCells, readLeafIdxCell(dbHeader, data))
	case INT_TAB_PAGE:
		page.IntTableCells = append(page.IntTableCells, readIntTableCell(dbHeader, data))
	case INT_IDX_PAGE:
		page.IntIdxCells = append(page.IntIdxCells, readIntIdxCell(dbHeader, data))
	default:
		panic("Unknown page type")
	}
}

func readLeafTableCell(header DatabaseHeader, data []byte) LeafTableCell {
	var cell LeafTableCell

	end := lo.Ternary(9 > len(data)-1, len(data)-1, 9)

	var offsetPayloadSize uint16
	cell.PayloadSize, offsetPayloadSize = decodeVarInt(data[:end])

	var offsetRowID uint16
	cell.RowID, offsetRowID = decodeVarInt(data[offsetPayloadSize : offsetPayloadSize+uint16(end)])

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

	return cell
}

func readLeafIdxCell(header DatabaseHeader, data []byte) LeafIdxCell {
	panic("not implemented")
}

func readIntTableCell(header DatabaseHeader, data []byte) IntTableCell {
	panic("not implemented")
}

func readIntIdxCell(header DatabaseHeader, data []byte) IntIdxCell {
	panic("not implemented")
}
