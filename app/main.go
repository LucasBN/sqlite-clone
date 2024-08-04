package main

import (
	"encoding/binary"
	"fmt"
	"github/com/lucasbn/sqlite-clone/app/machine"
	"github/com/lucasbn/sqlite-clone/app/machine/instructions"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/xwb1989/sqlparser"
)

// Usage: sqlite3.sh sample.db .dbinfo
func main() {

	instructions := []instructions.Instruction{
		instructions.Integer{Register: 1, Value: 1},
		instructions.Integer{Register: 2, Value: 2},
		instructions.ResultRow{FromRegister: 1, ToRegister: 2},
		instructions.Halt{},
	}
	spew.Dump(machine.Init(instructions).Run())

	return

	databaseFilePath := os.Args[1]
	command := os.Args[2]

	// Open the database file and defer its closing
	databaseFile, err := os.Open(databaseFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer databaseFile.Close()

	// Read the database header
	var header DatabaseHeader
	if err := binary.Read(databaseFile, binary.BigEndian, &header); err != nil {
		fmt.Println("Failed to read integer:", err)
		return
	}

	switch command {
	case ".dbinfo":
		dbinfo(databaseFile, header)
	case ".tables":
		tables(databaseFile, header)
	default:
		stmt, err := sqlparser.Parse(command)
		if err != nil {
			fmt.Println("Unknown command", command)
			os.Exit(1)
		}

		Execute(stmt, databaseFile, header)
	}

}

func dbinfo(databaseFile *os.File, header DatabaseHeader) {
	schema := readSQLiteSchema(databaseFile, header)

	fmt.Printf("database page size: %v\n", header.PageSize)
	fmt.Printf("number of tables: %v\n", schema.TableCount())
}

func tables(databaseFile *os.File, header DatabaseHeader) {
	schema := readSQLiteSchema(databaseFile, header)

	for _, row := range schema.Rows {
		if row.Type == "table" && row.Name != "sqlite_sequence" {
			fmt.Println(row.Name)
		}
	}
}
