package test_btree

import (
	"github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/types"
	"math/rand"
	"testing"
	"time"
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
		t.Errorf("Failed to rewind Cursor: %s", err)
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
				t.Errorf("Failed to read column (%d, %d): %s", i, j, err)
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
			t.Errorf("Failed to advance Cursor: %s", err)
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
	// 3,
	for i := 0; i < 3; i++ {
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

func TestRandomSingleLeafPage(t *testing.T) {
	// Generate a random seed based on the current time
	seed := time.Now().UnixNano()
	t.Logf("Using seed: %d", seed)

	// Create a new random generator with the seed
	r := rand.New(rand.NewSource(seed))

	// Randomly choose the number of columns between 1 and 10
	numColumns := r.Intn(10) + 1
	t.Logf("Number of columns: %d", numColumns)

	// Create a slice to hold the column types
	columnTypes := make([]int, numColumns)

	// Randomly populate the columnTypes slice with values 0, 1, or 2
	for i := 0; i < numColumns; i++ {
		columnTypes[i] = r.Intn(3)
	}

	// Generate a random leaf table page with the generated column types
	randomPage := GenerateRandomLeafTablePage(columnTypes, r)

	// Create the mock pager, which initially contains a single page with the
	// database header and is otherwise empty.
	pager := NewMockPager()

	// Write the table data
	pager.WritePage(2, randomPage.Serialize())

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

	for i := 0; i < len(randomPage.Cells); i++ {
		for j := 0; j < numColumns; j++ {
			entries := randomPage.Cells[i].Entries

			column, err := engine.ReadColumn(0, uint64(j))
			if err != nil {
				t.Errorf("Failed to read column (%d, %d)", i, j)
				return
			}

			switch expectedValue := entries[j].(type) {
			case types.NumberEntry:
				switch value := column.(type) {
				case types.NumberEntry:
					if value.Value != expectedValue.Value {
						t.Errorf("Expected %d, got %d (%d, %d)", expectedValue.Value, value.Value, i, j)
						return
					}
				default:
					t.Errorf("Expected NumberEntry, got %T (%d, %d)", value, i, j)
					return
				}
			case types.TextEntry:
				switch value := column.(type) {
				case types.TextEntry:
					if value.Value != expectedValue.Value {
						t.Errorf("Expected %s, got %s (%d, %d)", expectedValue.Value, value.Value, i, j)
						return
					}
				default:
					t.Errorf("Expected NumberEntry, got %T (%d, %d)", value, i, j)
					return
				}
			case types.NullEntry:
				switch value := column.(type) {
				case types.NullEntry:
					break
				default:
					t.Errorf("Expected NumberEntry, got %T (%d, %d)", value, i, j)
					return
				}
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
