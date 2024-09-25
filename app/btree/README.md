## Example Usage

```go
b.NewCursor(1, 5) // Cursor ID and root page
b.RewindCursor(1) // Set the cursor to the first entry on the B-Tree
for {
	v, _ := b.ReadColumn(1, 0) // Read the value from the first column

	// Do something with the value

	didAdvance, _ := b.AdvanceCursor(1) // Go to the next entry in the B-Tree
	if !didAdvance {
		break
	}
}
```

This tells us which details of the B-Tree to expose to the caller: they want to
traverse it as if it were a linear list and don't care that it is in a tree
structure.

This means that we (as in, this package) _do_ need to care about the fact that
our data is in a B-Tree format distributed among multiple pages.

As an initial attempt at an algorithm to sequentially scan a B-Tree:

- Start at the root node of the B-Tree
- Read the B-Tree page header
- Determine the page type
- Interior page:
	- Go to the first cell
	- Somehow store a record of the position on this page
	- Go the the left child pointer
	- Repeat from step 2 
- Leaf page:
	- Go to the first cell
	- Read the value in the first cell (may or may not overflow onto another page)
	- Move to the next cell
	- Repeat until all values in cells of leaf page have been read
	- Go back to the page you came from (if any) and move to the next cell on that

```go
type BTreeEngine[T any] interface {
	NewCursor(id uint64, rootPageNum uint64) (bool, error)
	RewindCursor(id uint64) (bool, error)
	AdvanceCursor(id uint64) (bool, error)
	ReadColumn(id uint64, column uint64) (T, error)
}
```

## Decomposing an SQLite B-Tree

A B-Tree page can either be a

- Leaf Table Page
- Interior Table Page
- Leaf Index Page
- Interior Index Page

Every B-Tree page has a header, and so the first primitive that might make sense
to build could be:

```go
func ReadBTreePageHeader(page []byte) (BTreeHeader, uint64, error)
```

where the second return argument is the number of bytes in the header (8 for
leaf, 12 for interior) - although this might be fine to just be derived from
the page type which is specified in the header.

Each page then has a cell pointer array and a cell content area. Therefore I
think we need another primitive that

- Read cell pointer array
- Get the pointer to cell N
- Get the raw bytes of cell N
