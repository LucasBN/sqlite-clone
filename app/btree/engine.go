package btree

type BTreeEngine[T any] struct {
	pager             pager
	cursors           map[uint64]*cursor
	resultConstructor resultConstructor[T]
}

func NewBTreeEngine[T any](pager pager, resultConstructor resultConstructor[T]) (*BTreeEngine[T], error) {
	return &BTreeEngine[T]{
		pager:             pager,
		cursors:           make(map[uint64]*cursor),
		resultConstructor: resultConstructor,
	}, nil
}

type pager interface {
	Close() error
	PageSize() uint64
	ReservedSpace() uint64
	GetPage(pageNum uint64) ([]byte, error)
}

type resultConstructor[T any] interface {
	Number(int64) T
	Text(string) T
	Null() T
}
