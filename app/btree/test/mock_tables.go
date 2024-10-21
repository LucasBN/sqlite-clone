package test_btree

import (
	"encoding/binary"
	"github/com/lucasbn/sqlite-clone/app/types"

	"math/rand"

	"github.com/samber/lo"
)

func FirstPage() []byte {
	header := MockDatabaseHeader{
		HeaderString:     [16]byte{'S', 'Q', 'L', 'i', 't', 'e', ' ', 'f', 'o', 'r', 'm', 'a', 't', ' ', '3', '\x00'},
		PageSize:         pageSize,
		FileWriteVersion: 1,
		FileReadVersion:  1,
		ReservedSpace:    0,
		Middle:           [38]byte{},
		TextEncoding:     1,
		End:              [40]byte{},
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
			PageType:               13,
			FirstFreeBlock:         0,
			NumCells:               3,
			CellContentOffset:      0,
			NumFragmentedFreeBytes: 0,
			RightMostPointer:       nil,
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

func generateRandomString(length int, r *rand.Rand) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}

func GenerateRandomLeafTablePage(columnTypes []int, r *rand.Rand) MockLeafTablePage {
	// The BTreeHeader is 8 bytes long, so we set this as the initial amount of
	// used bytes in the page
	totalSize := 8

	var cells []MockLeafTableCell

	for i := 0; true; i++ {
		var entries []types.Entry
		for _, columnType := range columnTypes {
			switch columnType {
			case 0:
				// Add a null entry
				entries = append(entries, types.NullEntry{})
			case 1:
				// Add a random int entry
				randomValue := r.Int63()
				entries = append(entries, types.NumberEntry{Value: randomValue})
			case 2:
				// Add a random text entry
				length := r.Intn(20) + 1
				randomString := generateRandomString(length, r)
				entries = append(entries, types.TextEntry{Value: randomString})
			}
		}

		cell := MockLeafTableCell{
			Key:     uint64(i),
			Entries: entries,
		}

		// Check if this page will cause the page to overflow and break
		// if it does. Each cell we add will use 2 bytes for the cell pointer,
		// and then however many bytes the cell itself takes up
		totalSize += 2
		totalSize += len(cell.Serialize())
		if totalSize > pageSize {
			break
		}

		cells = append(cells, cell)
	}

	return MockLeafTablePage{
		Header: MockPageHeader{
			PageType:               13,
			FirstFreeBlock:         0,
			NumCells:               uint16(len(cells)),
			CellContentOffset:      0,
			NumFragmentedFreeBytes: 0,
			RightMostPointer:       nil,
		},
		Cells: cells,
	}
}

type LeafOrInteriorPage struct {
	PageType int
	Leaf     MockLeafTablePage
	Interior MockInteriorTablePage
}

func GenerateRandomTable(firstPageNum uint32, depth uint64, columnTypes []int, r *rand.Rand) map[uint32]LeafOrInteriorPage {
	// The base case is when the depth is 1, in which case we create a leaf page
	if depth == 1 {
		return map[uint32]LeafOrInteriorPage{
			firstPageNum: {
				PageType: 13,
				Leaf:     GenerateRandomLeafTablePage(columnTypes, r),
			},
		}
	}

	// Otherwise, we need to create an interior page and then recursively create
	// the children of the interior page

	// The BTreeHeader is 12 bytes long, so we set this as the initial amount of
	// used bytes in the page.
	totalSize := 12

	pages := make(map[uint32]LeafOrInteriorPage)

	var cells []MockInteriorTableCell

	pagePointer := uint32(firstPageNum + 1)
	rightMostPointer := uint32(0)

	for i := 0; true; i++ {
		cell := MockInteriorTableCell{
			LeftChildPageNumber: pagePointer,
			Key:                 uint64(i),
		}

		// Generate the child page
		childPages := GenerateRandomTable(pagePointer, depth-1, columnTypes, r)

		// Add all of the child pages to the map
		for k, v := range childPages {
			pages[k] = v
		}

		// Check if this page will cause the page to overflow and break
		// if it does. Each cell we add will use 2 bytes for the cell pointer,
		// and then however many bytes the cell itself takes up
		totalSize += 2
		totalSize += len(cell.Serialize())
		if totalSize > pageSize {
			// Instead of appending the cell, we fill out the right most pointer
			// and break
			rightMostPointer = pagePointer
			break
		}

		cells = append(cells, cell)

		// The recursive call could result in multiple child pages being
		// generated
		pagePointer += uint32(len(childPages))
	}

	interiorPage := MockInteriorTablePage{
		Header: MockPageHeader{
			PageType:               5,
			FirstFreeBlock:         0,
			NumCells:               uint16(len(cells)),
			CellContentOffset:      0,
			NumFragmentedFreeBytes: 0,
			RightMostPointer:       lo.ToPtr(rightMostPointer),
		},
		Cells: cells,
	}

	// Add the interior page to the page map
	pages[firstPageNum] = LeafOrInteriorPage{
		PageType: 5,
		Interior: interiorPage,
	}

	return pages
}

func TableWithInteriorPage() map[int][]byte {

	pageTwo := MockInteriorTablePage{
		Header: MockPageHeader{
			PageType:               5,
			FirstFreeBlock:         0,
			NumCells:               2,
			CellContentOffset:      0,
			NumFragmentedFreeBytes: 0,
			RightMostPointer:       lo.ToPtr(uint32(5)),
		},
		Cells: []MockInteriorTableCell{
			{
				LeftChildPageNumber: 3,
				Key:                 1,
			},
			{
				LeftChildPageNumber: 4,
				Key:                 2,
			},
		},
	}

	pageThree := MockLeafTablePage{
		Header: MockPageHeader{
			PageType:               13,
			FirstFreeBlock:         0,
			NumCells:               1,
			CellContentOffset:      0,
			NumFragmentedFreeBytes: 0,
			RightMostPointer:       nil,
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
			PageType:               13,
			FirstFreeBlock:         0,
			NumCells:               1,
			CellContentOffset:      0,
			NumFragmentedFreeBytes: 0,
			RightMostPointer:       nil,
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

	pageFive := MockLeafTablePage{
		Header: MockPageHeader{
			PageType:               13,
			FirstFreeBlock:         0,
			NumCells:               1,
			CellContentOffset:      0,
			NumFragmentedFreeBytes: 0,
			RightMostPointer:       nil,
		},
		Cells: []MockLeafTableCell{
			{
				Key: 3,
				Entries: []types.Entry{
					types.NumberEntry{Value: 3},
				},
			},
		},
	}

	return map[int][]byte{
		2: pageTwo.Serialize(),
		3: pageThree.Serialize(),
		4: pageFour.Serialize(),
		5: pageFive.Serialize(),
	}
}
