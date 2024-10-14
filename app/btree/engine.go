package btree

// A BTreeEngine is a high-level interface to the underlying B-Tree implentation
// (raw bytes spread across many pages) that is used to store tables in SQLite.
// This package exports many "operations" that can be performed on a
// BTreeEngine, which are predominantly centered around controlling "cursors".
// The caller is able to create cursors, which point to specific entries in a
// B-Tree, and then able to perform operations on those cursors to move to
// different entries in the B-Tree.
type BTreeEngine[T any] struct {
	pager             pager
	cursors           map[uint64]*cursor
	resultConstructor resultConstructor[T]
}

// NewBTreeEngine creates a new BTreeEngine that is backed by the given pager.
// The resultConstructor tells the BTreeEngine how to convert the types that it
// reads from the database into the type T that the caller has provided.
func NewBTreeEngine[T any](pager pager, resultConstructor resultConstructor[T]) (*BTreeEngine[T], error) {
	return &BTreeEngine[T]{
		pager:             pager,
		cursors:           make(map[uint64]*cursor),
		resultConstructor: resultConstructor,
	}, nil
}

// The pager interface is used to abstract the underlying storage mechanism of
// the B-Tree. The creator of the BTreeEngine can provide any implementation
// that satisfies this interface (which are methods required by the BTreeEngine
// to function).
type pager interface {
	Close() error
	PageSize() uint64
	ReservedSpace() uint64
	GetPage(pageNum uint64) ([]byte, error)
}

// The resultConstructor is provided when creating a BTreeEngine and means that
// the caller can specify how to transform the Number/Text/Null results into
// the types that they want.
type resultConstructor[T any] interface {
	Number(int64) T
	Text(string) T
	Null() T
}
