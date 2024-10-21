package test_btree

import (
	"fmt"
	"github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/types"
	"math/rand"
	"testing"
	"time"
)

func TestRandomTable(t *testing.T) {
	heights := []uint64{1, 2, 3}

	for _, height := range heights {
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
		randomPages := GenerateRandomTable(2, height, columnTypes, r)

		// Create the mock pager, which initially contains a single page with the
		// database header and is otherwise empty.
		pager := NewMockPager()

		// Write the table data to the pager
		pager.WritePages(randomPages)

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

		for k := 2; k < len(randomPages)+2; k++ {
			// If the page is an interior page, skip it
			if randomPages[uint32(k)].PageType == 5 {
				continue
			}

			for i := 0; i < len(randomPages[uint32(k)].Leaf.Cells); i++ {
				for j := 0; j < numColumns; j++ {
					entries := randomPages[uint32(k)].Leaf.Cells[i].Entries

					column, err := engine.ReadColumn(0, uint64(j))
					if err != nil {
						t.Errorf("Failed to read column (%d, %d)", i, j)
						return
					}

					result := CompareValue(entries[j], column)
					if result != "" {
						t.Errorf(result+" (%d, %d)", i, j)
						return
					}
				}

				_, err = engine.AdvanceCursor(0)
				if err != nil {
					t.Errorf("Failed to advance Cursor")
					return
				}
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
}

func CompareValue(expectedValue types.Entry, value types.Entry) string {
	switch expectedValue := expectedValue.(type) {
	case types.NumberEntry:
		switch value := value.(type) {
		case types.NumberEntry:
			if value.Value != expectedValue.Value {
				return fmt.Sprintf("Expected %d, got %d", expectedValue.Value, value.Value)
			}
		default:
			return fmt.Sprintf("Expected NumberEntry, got %T", value)
		}
	case types.TextEntry:
		switch value := value.(type) {
		case types.TextEntry:
			if value.Value != expectedValue.Value {
				return fmt.Sprintf("Expected %s, got %s", expectedValue.Value, value.Value)
			}
		default:
			return fmt.Sprintf("Expected TextEntry, got %T", value)
		}
	case types.NullEntry:
		switch value := value.(type) {
		case types.NullEntry:
			break
		default:
			return fmt.Sprintf("Expected NullEntry, got %T", value)
		}
	}

	return ""
}
