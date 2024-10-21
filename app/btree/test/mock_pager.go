package test_btree

const pageSize = 4096

type MockPager struct {
	pages map[int][]byte
}

func NewMockPager() *MockPager {
	return &MockPager{
		pages: map[int][]byte{
			1: FirstPage(),
		},
	}
}

func (p *MockPager) WritePage(pageNum int, data []byte) error {
	p.pages[pageNum] = data
	return nil
}

func (p *MockPager) WritePages(pageMap map[uint32]LeafOrInteriorPage) error {
	for pageNum, data := range pageMap {
		switch data.PageType {
		case 5:
			p.pages[int(pageNum)] = data.Interior.Serialize()
		case 13:
			p.pages[int(pageNum)] = data.Leaf.Serialize()
		default:
			return nil
		}
	}
	return nil
}

func (p *MockPager) Close() error {
	return nil
}

func (p *MockPager) PageSize() uint64 {
	return pageSize
}

func (p *MockPager) ReservedSpace() uint64 {
	return 0
}

func (p *MockPager) GetPage(pageNum uint64) ([]byte, error) {
	if data, ok := p.pages[int(pageNum)]; ok {
		return data, nil
	}
	return nil, nil
}
