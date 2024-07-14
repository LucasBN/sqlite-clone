# Architecture & Design Decisions

The entry point is `main.go`, which is responsible for opening the database
file, reading the database header and handling the execution of the command by
handing the responsibility to another function.

## Reading database tables

I've made a few design decisions related to how pages, cells and records are
stored, as well as how tables are read.

A cell is treated as something that exists entirely within the context of a
b-tree page, and the type of the cell is entirely dependent on the type of
b-tree page. A b-tree page looks like:

```
type BTreePage struct {
	Header         BTreeHeader
	IntIdxCells    []IntIdxCell
	IntTableCells  []IntTableCell
	LeafIdxCells   []LeafIdxCell
	LeafTableCells []LeafTableCell
}
```

The header contains information like the page type, number of cells and cell
content offset (among others, as defined by the SQLite database file format).

I briefly considered separating out cells and then attempting to traverse all
overflow pages to bring the entire payload together, but I decided this was a
really bad idea because:

i. I'm doing this inside a function called `readBTreePage` which takes a page
number and it feels weird for this to make subsequent calls to read other pages
and to return data that is actually stored on other pages.

ii. I'm not sure how this would end up looking for interior pages (which I
haven't yet implemented). My suspicion is that it would be quite horrible.

I decided to go for an approach which seems much more sensible: the
`readBTreePage` function would only read and return the data on the page that it
was asked about, and the caller would be responsible for interpreting that data.

Right now, the only caller is `readTable`, which takes a `rootPage`. This
function first calls `readBTreePage` on the root page itself, but of course the
data of the table may be split across multiple pages - which we only find out
after we have read the root page. After reading the root page, `readTable` will
read all other necessary pages to create a list of 'raw records' which are
essentially byte slices that contain *all* of the data for a record and have no
dependency or reference to pages. This allows us to pass the entire raw record
to a function that is responsible for decoding the record into a useable format
(see `record.go`).

So far, I think this is a good structure and will allow me to (fingers crossed)
easily add support for things like overflown cells and interior pages