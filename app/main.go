package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

// Usage: sqlite3.sh sample.db .dbinfo
func main() {
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
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}

func dbinfo(databaseFile *os.File, header DatabaseHeader) {
	schema := readSQLiteSchema(databaseFile, header)

	fmt.Printf("database page size: %v\n", header.PageSize)
	fmt.Printf("number of tables: %v\n", schema.TableCount())
}
