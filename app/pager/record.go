package pager

import (
	"encoding/binary"
	"fmt"

	"github.com/samber/lo"
)

type Records []Record

type Record struct {
	Entries []RecordEntry
}

// 0 -> null
// 1 -> number
// 2 -> text
type RecordEntry struct {
	RecordEntryType uint8
	Null            *bool
	Number          *uint64
	Text            *string
}

func ReadRecord(rawRecord []byte) Record {
	// We first have to read the entire raw record header, which is at the
	// beginning of the raw record and is made up of one or more varints.
	// The first varint gives the size of the header, and we should keep
	// reading varints from the remainder until we reach this size.
	recordHeaderSize, offset := decodeVarInt(rawRecord)
	typeCodes := readTypeCodes(rawRecord[offset:recordHeaderSize])

	// Keep a pointer to the start of the actual record data, which we will move
	// whenever we read a record entry that takes up space in the raw record
	pointer := uint16(recordHeaderSize)

	// The remainder of the raw record is the actual record data, which we can
	// now read because we know the type codes
	var entries []RecordEntry
	for _, typeCode := range typeCodes {
		entry, offset := readRecordEntry(typeCode, rawRecord[pointer:])
		entries = append(entries, entry)
		pointer += offset
	}

	return Record{Entries: entries}
}

func readTypeCodes(header []byte) []uint64 {
	var typeCodes []uint64
	pointer := uint16(0)

	for {
		typeCode, offset := decodeVarInt(header[pointer:])
		typeCodes = append(typeCodes, typeCode)

		pointer += offset
		if int(pointer) >= len(header) {
			break
		}
	}

	return typeCodes
}

func readRecordEntry(typeCode uint64, data []byte) (RecordEntry, uint16) {
	intTypeCodeToSize := map[uint64]uint16{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 6,
		6: 8,
	}

	if typeCode == 0 {
		return RecordEntry{RecordEntryType: 0, Null: lo.ToPtr(true)}, 0
	} else if typeCode >= 1 && typeCode <= 6 {
		// Make an 8 byte empty slice
		result := make([]byte, 8)

		// Extract the correct number of bytes from the raw record
		size := intTypeCodeToSize[uint64(typeCode)]
		value := data[:size]

		// Copy the bytes into the result slice so that we can decode them
		// as a uint64
		for i := 7; 7-i < len(value); i-- {
			result[i] = value[7-i]
		}

		return RecordEntry{RecordEntryType: 1, Number: lo.ToPtr(binary.BigEndian.Uint64(result))}, size
	} else if typeCode == 8 {
		return RecordEntry{RecordEntryType: 1, Number: lo.ToPtr(uint64(0))}, 0
	} else if typeCode == 9 {
		return RecordEntry{RecordEntryType: 1, Number: lo.ToPtr(uint64(1))}, 0
	} else if typeCode >= 12 && typeCode%2 == 1 {
		length := uint16((typeCode - 13) / 2)
		return RecordEntry{RecordEntryType: 2, Text: lo.ToPtr(string(data[:length]))}, length
	} else {
		panic(fmt.Sprintf("Unimplemented: type code decoder %v", typeCode))
	}
}
