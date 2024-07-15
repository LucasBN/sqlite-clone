package main

import (
	"fmt"
	"os"

	"github.com/xwb1989/sqlparser"
)

func Execute(stmt sqlparser.Statement, databaseFile *os.File, header DatabaseHeader) {
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		// spew.Dump(stmt)

		// Get all of the records from the table
		var records []Record

		switch from := stmt.From[0].(type) {
		case *sqlparser.AliasedTableExpr:

			switch expr := from.Expr.(type) {
			case sqlparser.TableName:
				schema := readSQLiteSchema(databaseFile, header)
				for _, row := range schema.Rows {
					if row.Type == "table" && row.Name == expr.Name.String() {
						records = append(records, readTable(databaseFile, header, row.RootPage)...)
					}
				}
			default:
				panic("Unimplemented: SQL FROM TABLE NAME type")
			}

		default:
			panic("Unimplemented: SQL FROM statement type")
		}

		// We support returning the entire table, or returing the number of rows
		// on the entire table
		switch selectExpr := stmt.SelectExprs[0].(type) {
		case *sqlparser.AliasedExpr:
			switch expr := selectExpr.Expr.(type) {
			case *sqlparser.FuncExpr:
				if expr.Name.String() == "COUNT" {
					fmt.Println(len(records))
				}
			default:
				panic("Unimplemented: SQL SELECT EXPR statement type")
			}
		case *sqlparser.StarExpr:
			for _, record := range records {
				for _, column := range record.Entries {

					switch column.RecordEntryType {
					case 0:
						fmt.Print("NULL | ")
					case 1:
						fmt.Printf("%v | ", *column.Number)
					case 2:
						fmt.Printf("%s | ", *column.Text)
					}
				}
				fmt.Println("")
			}
		default:
			panic("Unimplemented: SQL SELECT EXPR statement type")
		}
	default:
		panic("Unimplemented: SQL statement type")
	}
}
