package pager

import (
	"encoding/binary"
	"log"
	"os"
)

type Pager struct {
	File     *os.File
	DBHeader DatabaseHeader
	Cache    map[uint32]BTreePage
}

func Init(filepath string) *Pager {
	databaseFile, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	// defer databaseFile.Close()

	// TODO: This is just not the place to do this really
	var header DatabaseHeader
	if err := binary.Read(databaseFile, binary.BigEndian, &header); err != nil {
		log.Fatal(err)
	}

	return &Pager{
		File:     databaseFile,
		DBHeader: header,
		Cache:    make(map[uint32]BTreePage),
	}
}

func (p *Pager) GetPage(pageNum uint32) BTreePage {
	if _, ok := p.Cache[pageNum]; !ok {
		p.Cache[pageNum] = readBTreePage(p.File, p.DBHeader, pageNum)
	}
	return p.Cache[pageNum]
}
