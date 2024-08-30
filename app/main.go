package main

import (
	"encoding/binary"
	btreepkg "github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/generator"
	"github/com/lucasbn/sqlite-clone/app/machine"
	pagerpkg "github/com/lucasbn/sqlite-clone/app/pager"
	"github/com/lucasbn/sqlite-clone/app/parser"
	"github/com/lucasbn/sqlite-clone/app/types"
	"os"

	"github.com/davecgh/go-spew/spew"
)

type DatabaseHeader struct {
	HeaderString     [16]byte
	PageSize         uint16
	FileWriteVersion uint8
	FileReadVersion  uint8
	ReservedSpace    uint8
	Middle           [38]byte
	TextEncoding     uint32
	End              [40]byte
}

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

	// Get the database header
	header, err := getDatabaseHeader(dbFilePath)
	if err != nil {
		panic(err)
	}

	// Initialise a pager and defer closing it
	pager, err := pagerpkg.NewPager(
		dbFilePath,
		pagerpkg.PagerConfig{
			PageSize:      uint64(header.PageSize),
			ReservedSpace: uint64(header.ReservedSpace),
		},
	)
	if err != nil {
		panic(err)
	}
	defer pager.Close()

	// Initialise a BTreeEngine
	bTreeEngine, err := btreepkg.NewBTreeEngine(pager, &types.EntryConstructor{})
	if err != nil {
		panic(err)
	}

	// 1. Parse the SQL string
	stmt := parser.MustParse(command)

	// 2. Generate the bytecode
	instructions := generator.Generate[types.Entry](stmt)

	// 3. Configure the virtual machine
	m := machine.NewMachine(
		machine.MachineConfig[types.Entry]{
			Instructions: instructions,
			BTreeEngine:  bTreeEngine,
		},
	)

	// 4. Execute the program
	result := m.Run()

	// 5. Pretty print the result
	spew.Dump(result)
}

func getDatabaseHeader(filepath string) (*DatabaseHeader, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var header DatabaseHeader
	if err := binary.Read(file, binary.BigEndian, &header); err != nil {
		return nil, err
	}

	return &header, nil
}
