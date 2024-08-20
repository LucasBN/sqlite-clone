package pager

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

func readSQLiteSchema(databaseFile *os.File, dbHeader DatabaseHeader) SQLiteSchema {
	records := readTable(databaseFile, dbHeader, 1)

	var rows []SQLiteSchemaTuple
	for _, record := range records {
		row := SQLiteSchemaTuple{
			Type:      *record.Entries[0].Text,
			Name:      *record.Entries[1].Text,
			TableName: *record.Entries[2].Text,
			RootPage:  uint32(*record.Entries[3].Number),
			SQL:       *record.Entries[4].Text,
		}
		rows = append(rows, row)
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
