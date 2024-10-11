package test_btree

import (
	"encoding/binary"
	"github/com/lucasbn/sqlite-clone/app/types"

	"github.com/samber/lo"
)

func FirstPage() []byte {
	header := MockDatabaseHeader{
		HeaderString: [16]byte{'S', 'Q', 'L', 'i', 't', 'e', ' ', 'f', 'o', 'r', 'm', 'a', 't', ' ', '3', '\x00'},
		PageSize: pageSize,
		FileWriteVersion: 1,
		FileReadVersion: 1,
		ReservedSpace: 0,
		Middle: [38]byte{},
		TextEncoding: 1,
		End: [40]byte{},
	}

	buf := make([]byte, pageSize)
	copy(buf[:100], header.HeaderString[:])
	binary.BigEndian.PutUint16(buf[16:18], header.PageSize)
	buf[18] = header.FileWriteVersion
	buf[19] = header.FileReadVersion
	buf[20] = header.ReservedSpace
	copy(buf[21:59], header.Middle[:])
	binary.BigEndian.PutUint32(buf[59:63], header.TextEncoding)
	copy(buf[63:103], header.End[:])

	return buf
}

func TableWithSingleLeafPage() map[int][]byte {

	pageTwo := MockLeafTablePage{
		Header: MockPageHeader{
			PageType: 13,
			FirstFreeBlock: 0,
			NumCells: 3,
			CellContentOffset: 0,
			NumFragmentedFreeBytes: 0,
			RightMostPointer: nil,
		},
		Cells: []MockLeafTableCell{
			{
				Key: 1,
				Entries: []types.Entry{
					types.NumberEntry{Value: 1},
					types.NumberEntry{Value: 2},
					types.NumberEntry{Value: 3},
				},
			},
			{
				Key: 2,
				Entries: []types.Entry{
					types.NumberEntry{Value: 4},
					types.NumberEntry{Value: 5},
					types.NumberEntry{Value: 6},
				},
			},
			{
				Key: 3,
				Entries: []types.Entry{
					types.NumberEntry{Value: 7},
					types.NumberEntry{Value: 8},
					types.NumberEntry{Value: 9},
				},
			},
		},
	}

	return map[int][]byte{
		2: pageTwo.Serialize(),
	}
}

func TableWithInteriorPage() map[int][]byte {

	pageTwo := MockInteriorTablePage{
		Header: MockPageHeader{
			PageType: 5,
			FirstFreeBlock: 0,
			NumCells: 1,
			CellContentOffset: 0,
			NumFragmentedFreeBytes: 0,
			RightMostPointer: lo.ToPtr(uint32(4)),
		},
		Cells: []MockInteriorTableCell{
			{
				LeftChildPageNumber: 3,
				Key: 1,
			},
		},
	}

	pageThree := MockLeafTablePage{
		Header: MockPageHeader{
			PageType: 13,
			FirstFreeBlock: 0,
			NumCells: 1,
			CellContentOffset: 0,
			NumFragmentedFreeBytes: 0,
			RightMostPointer: nil,
		},
		Cells: []MockLeafTableCell{
			{
				Key: 1,
				Entries: []types.Entry{
					types.NumberEntry{Value: 1},
				},
			},
		},
	}

	pageFour := MockLeafTablePage{
		Header: MockPageHeader{
			PageType: 13,
			FirstFreeBlock: 0,
			NumCells: 1,
			CellContentOffset: 0,
			NumFragmentedFreeBytes: 0,
			RightMostPointer: nil,
		},
		Cells: []MockLeafTableCell{
			{
				Key: 2,
				Entries: []types.Entry{
					types.NumberEntry{Value: 2},
				},
			},
		},
	}

	return map[int][]byte{
		2: pageTwo.Serialize(),
		3: pageThree.Serialize(),
		4: pageFour.Serialize(),
	}
}