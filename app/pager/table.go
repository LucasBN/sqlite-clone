package pager

import "os"

func readTable(databaseFile *os.File, dbHeader DatabaseHeader, rootPage uint32) []Record {
	// Read the first page, which is always a table b-tree page (but can either
	// be a leaf page or an interior page)
	page := readBTreePage(databaseFile, dbHeader, rootPage)

	// In the simplest case, each cell on the first page corresponds to an
	// entire record. However, if the first page is an interior page or if a
	// cell overflows onto another page, the records for the table will be
	// spread across multiple pages. When decoding records, we're not really
	// interested in whether or not they span multiple pages, so we're first
	// going to do some work to collect all of the "raw records" together (which
	// will potentially involve reading more pages), and then we can decode
	// these records in their entirety.
	var rawRecords [][]byte

	switch page.Header.PageType {
	case LEAF_TAB_PAGE:
		for _, cell := range page.LeafTableCells {
			// TODO: support the cell overflowing
			if cell.OverflowPage != nil {
				panic("Unimplemented: cell has overflow page")
			}

			// This is separate here to serve as a reminder that we would in
			// theory need to do some work to combine the cell payload with the
			// overflow page(s)
			cellPayload := cell.Payload
			rawRecords = append(rawRecords, cellPayload)
		}
	case INT_TAB_PAGE:
		panic("Unimplemented: SQLite schema root page is an interior page")
	default:
		panic("SQLite schema root page is neither a leaf page nor an interior page")
	}

	// Decode the raw records into actual records
	var records []Record
	for _, rawRecord := range rawRecords {
		records = append(records, ReadRecord(rawRecord))
	}

	return records
}
