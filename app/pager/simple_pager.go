package pager

import (
	"io"
	"os"
)

// SimplePager uses a basic caching mechanism to store pages in memory once
// they've been accessed. It makes no attempt at any other optimizations.
type SimplePager struct {
	File   *os.File
	Cache  map[uint64]*Page
	Config simplePagerConfig
}

type simplePagerConfig struct {
	PageSize      uint64
	ReservedSpace uint64
}

var _ Pager = &SimplePager{}

func NewSimplePager(filepath string, pageSize uint64, reservedSpace uint64) (*SimplePager, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	return &SimplePager{
		File:   file,
		Config: simplePagerConfig{PageSize: pageSize, ReservedSpace: reservedSpace},
		Cache:  make(map[uint64]*Page),
	}, nil
}

func (p *SimplePager) Close() error {
	return p.File.Close()
}

func (p *SimplePager) PageSize() uint64 {
	return p.Config.PageSize
}

func (p *SimplePager) ReservedSpace() uint64 {
	return p.Config.ReservedSpace
}

func (p *SimplePager) GetPage(pageNum uint64) (*Page, error) {
	// If the page isn't already in the cache, we should read it directly from
	// the file
	if _, ok := p.Cache[pageNum]; !ok {
		// Calculate the byte number at which this page starts
		pageStart := int64((pageNum - 1) * uint64(p.PageSize()))

		// Seek to the beginning of the page
		p.File.Seek(pageStart, io.SeekStart)

		// Read the entire page into a byte slice
		pageBytes := make([]byte, p.PageSize())
		if _, err := p.File.Read(pageBytes); err != nil {
			return nil, err
		}

		// Insert the page into cache
		p.Cache[pageNum] = &Page{
			Bytes:  pageBytes,
			Offset: uint64(pageStart),
		}
	}

	return p.Cache[pageNum], nil
}
