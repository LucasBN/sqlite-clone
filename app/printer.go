package main

import (
	"fmt"
)

func (r Records) print() {

	var rows [][]string
	for _, record := range r {
		var row []string
		for _, column := range record.Entries {
			switch column.RecordEntryType {
			case 0:
				row = append(row, "NULL")
			case 1:
				row = append(row, fmt.Sprintf("%v", *column.Number))
			case 2:
				row = append(row, *column.Text)
			}
		}
		rows = append(rows, row)
	}

	var columnWidths []int
	for _, row := range rows {
		for i, column := range row {
			if i >= len(columnWidths) {
				columnWidths = append(columnWidths, len(column))
			} else if len(column) > columnWidths[i] {
				columnWidths[i] = len(column)
			}
		}
	}

	for _, row := range rows {
		for i, column := range row {
			fmt.Printf("%-*s", columnWidths[i], column)
			if i != len(row)-1 {
				fmt.Print(" | ")
			}
		}
		fmt.Println()
	}

}
