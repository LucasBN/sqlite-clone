package pager

// A `Pager` is responsible for reading and writing pages to and from the disk.
type Pager interface {
	Close() error
	PageSize() uint64
	ReservedSpace() uint64
	GetPage(pageNum uint64) (*Page, error)
}

type Page struct {
	// The raw bytes of the page
	Bytes []byte

	// The offset from the beginning of the file at which the page starts
	Offset uint64
}
