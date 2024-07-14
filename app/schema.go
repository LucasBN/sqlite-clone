package main

import (
	"os"
)

type SQLiteSchemaTuple struct {
	Type      string
	Name      string
	TableName string
	RootPage  uint32
	SQL       string
}

type SQLiteSchema struct {
	Rows []SQLiteSchemaTuple
}

// readSQLiteSchema does the following:
// - Reads the first page of the database file
// - If the page is a leaf page:
//   - Iterates over each cell and collects the entire cell payload, which may
//     involve reading overflow pages (currently unsupported)
//   - Decodes the cell payload into a record format
//
// - If the page is an interior page:
//   - Currently unsupported
func readSQLiteSchema(databaseFile *os.File, dbHeader DatabaseHeader) SQLiteSchema {
	// Read the first page, which is always a table b-tree page (but can either
	// be a leaf page or an interior page)
	page := readBTreePage(databaseFile, dbHeader, 0)

	var rows []SQLiteSchemaTuple
	switch page.Header.PageType {
	case LEAF_TAB_PAGE:
		for _, cell := range page.LeafTableCells {
			// We don't yet support the cell overflowing
			if cell.OverflowPage != nil {
				panic("Unimplemented: cell has overflow page")
			}

			// This is separate here to serve as a reminder that we would in
			// theory need to do some work to combine the cell payload with the
			// overflow page(s)
			cellPayload := cell.Payload

			record := readRecord(cellPayload)

			row := SQLiteSchemaTuple{
				Type:      *record.Entries[0].Text,
				Name:      *record.Entries[1].Text,
				TableName: *record.Entries[2].Text,
				RootPage:  uint32(*record.Entries[3].Number),
				SQL:       *record.Entries[4].Text,
			}
			rows = append(rows, row)
		}
	case INT_TAB_PAGE:
		panic("Unimplemented: SQLite schema root page is an interior page")
	default:
		panic("SQLite schema root page is neither a leaf page nor an interior page")
	}

	return SQLiteSchema{Rows: rows}
}

func (schema SQLiteSchema) TableCount() int {
	count := 0
	for _, row := range schema.Rows {
		if row.Type == "table" {
			count += 1
		}
	}
	return count
}
