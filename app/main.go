package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	// "github.com/xwb1989/sqlparser"
)

// Usage: sqlite3.sh sample.db .dbinfo
func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]

	databaseFile, err := os.Open(databaseFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer databaseFile.Close()

	switch command {
	case ".dbinfo":
		// Read the database header
		var header DatabaseHeader
		if err := binary.Read(databaseFile, binary.BigEndian, &header); err != nil {
			fmt.Println("Failed to read integer:", err)
			return
		}

		// Print the page size
		fmt.Printf("database page size: %v\n", header.PageSize)

		// Print the text encoding size
		fmt.Printf("database text encoding: %v\n", header.TextEncoding)

		schema := readSQLiteSchema(databaseFile, header)

		// Print the number of tables
		fmt.Printf("number of tables: %v\n", schema.TableCount())

	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
