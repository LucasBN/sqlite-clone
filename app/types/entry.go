package types

// Entry allows us to represent the different types of entries that can be
// stored. We have three types of entries: null, number, and text.
//
// The caller can use the Type method to determine the type of the entry:
//
// var entry RecordEntry
// entry = NumberEntry{Value: 12345}
//
// switch v := entry.(type) {
// ...
// }
//
// It might be useful at some point in the future to add methods to this
// interface.
type Entry interface{}

type NullEntry struct{}

type NumberEntry struct {
	Value int64
}

type TextEntry struct {
	Value string
}

type EntryConstructor struct{}

func (r *EntryConstructor) Number(value int64) Entry {
	return NumberEntry{Value: value}
}

func (r *EntryConstructor) Text(value string) Entry {
	return TextEntry{Value: value}
}

func (r *EntryConstructor) Null() Entry {
	return NullEntry{}
}

var _ Entry = &NullEntry{}
var _ Entry = &NumberEntry{}
var _ Entry = &TextEntry{}
