package btree

import "github/com/lucasbn/sqlite-clone/app/pager"

type Cursor struct {
	ID          uint32
	CurrentPage uint32
	Position    uint32
	Pager       *pager.Pager
}

// Initializes the cursor so that it points at the first row in the table on the
// root page.
func InitCursor(id uint32, rootPage uint32, pager *pager.Pager) *Cursor {
	return &Cursor{
		ID:          id,
		CurrentPage: rootPage,
		Position:    0,
		Pager:       pager,
	}
}

// Next advances the cursor to the next row in the table.
func (c *Cursor) Next() bool {
	if c.Position >= 2 {
		return false
	}

	c.Position++
	return true
}

// ReadColumn reads the value of the column at the current row.
func (c *Cursor) ReadColumn(column int) int {

	page := c.Pager.GetPage(c.CurrentPage)

	entries := pager.ReadRecord(page.LeafTableCells[c.Position].Payload).Entries
	value := entries[column].Number

	return int(*value)
}
