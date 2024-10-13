package test_btree

import (
	"github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/types"
	"testing"
)

func TestReadSingleLeafPageTable(t *testing.T) {
	// Create the mock pager, which initially contains a single page with the
	// database header and is otherwise empty.
	pager := NewMockPager()

	// Write the table data
	for pageNum, data := range TableWithSingleLeafPage() {
		pager.WritePage(pageNum, data)
	}

	engine, err := btree.NewBTreeEngine(pager, &types.EntryConstructor{})
	if err != nil {
		t.Errorf("Failed to create BTreeEngine")
		return
	}

	// Create a cursor that points to the second page (which is the leaf table
	// page).
	ok, err := engine.NewCursor(0, 2)
	if !ok || err != nil {
		t.Errorf("Failed to create new Cursor")
		return
	}

	// Move the cursor to the first entry in the table
	ok, err = engine.RewindCursor(0)
	if !ok || err != nil {
		t.Errorf("Failed to rewind Cursor")
		return
	}

	// We expect to read three rows:
	// 1, 2, 3
	// 4, 5, 6
	// 7, 8, 9
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			column, err := engine.ReadColumn(0, uint64(j))
			if err != nil {
				t.Errorf("Failed to read column (%d, %d)", i, j)
				return
			}
			switch value := column.(type) {
			case types.NumberEntry:
				expecting := int64((3 * i) + j + 1)
				if value.Value != expecting {
					t.Errorf("Expected %d, got %d (%d, %d)", expecting, value.Value, i, j)
					return
				}
			default:
				t.Errorf("Expected NumberEntry, got %T (%d, %d)", value, i, j)
				return
			}
		}

		_, err = engine.AdvanceCursor(0)
		if err != nil {
			t.Errorf("Failed to advance Cursor")
			return
		}
	}

	// Check that we have reached the end of the table
	didAdvance, err := engine.AdvanceCursor(0)
	if err != nil {
		t.Errorf("Failed to advance Cursor")
		return
	}
	if didAdvance {
		t.Errorf("Expected to have reached end of table")
		return
	}
}

func TestReadSingleInteriorPageTable(t *testing.T) {
	// Create the mock pager, which initially contains a single page with the
	// database header and is otherwise empty.
	pager := NewMockPager()

	// Write the table data
	for pageNum, data := range TableWithInteriorPage() {
		pager.WritePage(pageNum, data)
	}

	engine, err := btree.NewBTreeEngine(pager, &types.EntryConstructor{})
	if err != nil {
		t.Errorf("Failed to create BTreeEngine")
		return
	}

	// Create a cursor that points to the second page (which is the leaf table
	// page).
	ok, err := engine.NewCursor(0, 2)
	if !ok || err != nil {
		t.Errorf("Failed to create new Cursor")
		return
	}

	// Move the cursor to the first entry in the table
	ok, err = engine.RewindCursor(0)
	if !ok || err != nil {
		t.Errorf("Failed to rewind Cursor")
		return
	}

	// We expect to read two rows:
	// 1,
	// 2,
	for i := 0; i < 1; i++ {
		for j := 0; j < 1; j++ {
			column, err := engine.ReadColumn(0, uint64(j))
			if err != nil {
				t.Errorf("Failed to read column (%d, %d)", i, j)
				return
			}
			switch value := column.(type) {
			case types.NumberEntry:
				expecting := int64(i + j + 1)
				if value.Value != expecting {
					t.Errorf("Expected %d, got %d (%d, %d)", expecting, value.Value, i, j)
					return
				}
			default:
				t.Errorf("Expected NumberEntry, got %T (%d, %d)", value, i, j)
				return
			}
		}

		_, err = engine.AdvanceCursor(0)
		if err != nil {
			t.Errorf("Failed to advance Cursor")
			return
		}
	}

	// Check that we have reached the end of the table
	didAdvance, err := engine.AdvanceCursor(0)
	if err != nil {
		t.Errorf("Failed to advance Cursor")
		return
	}
	if didAdvance {
		t.Errorf("Expected to have reached end of table")
		return
	}
}
