package main

import (
	"github/com/lucasbn/sqlite-clone/app/generator"
	"github/com/lucasbn/sqlite-clone/app/machine"
	"github/com/lucasbn/sqlite-clone/app/parser"
	"os"

	"github.com/davecgh/go-spew/spew"
)

// Usage:
// - sqlite3.sh sample.db .dbinfo
// - sqlite3.sh sample.db "SELECT * FROM users;"
//
// This is really just a temporary entry point into the system. In the future we
// could add support for some sort of REPL... but that's not really the point of
// doing this project so I'll leave that for a rainy day.
func main() {
	dbFilePath := os.Args[1]
	command := os.Args[2]

	// 1. Parse the SQL string
	stmt := parser.MustParse(command)

	// 2. Generate the bytecode
	instructions := generator.Generate(stmt)

	// 3. Configure the virtual machine
	m := machine.Init(machine.MachineConfig{
		Instructions: instructions,
		DBFilePath:   dbFilePath,
	})

	// 4. Execute the program
	result := m.Run()

	// 5. Pretty print the result
	spew.Dump(result)
}

// Old code that I may use someday:
// ----------------------------------------------------------
// databaseFilePath := os.Args[1]
// command := os.Args[2]

// // Open the database file and defer its closing
// databaseFile, err := os.Open(databaseFilePath)
// if err != nil {
// 	log.Fatal(err)
// }
// defer databaseFile.Close()

// // Read the database header
// var header DatabaseHeader
// if err := binary.Read(databaseFile, binary.BigEndian, &header); err != nil {
// 	fmt.Println("Failed to read integer:", err)
// 	return
// }

// switch command {
// case ".dbinfo":
// 	dbinfo(databaseFile, header)
// case ".tables":
// 	tables(databaseFile, header)
// default:
// 	break
// }
