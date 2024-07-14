package main

import (
	"encoding/binary"
	"fmt"

	"github.com/samber/lo"
)

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

func readRecord(payload []byte) Record {
	// Since the payload consists of a series of varints, we keep a pointer so
	// that we know where we are in the payload.
	pointer := uint16(0)

	// We first have to read the entire payload header, which is a the
	// beginning of the payload and is made up of one or more varints.
	// The first varint gives the size of the header, and we should keep
	// reading varints from the remainder until we reach this size.
	payloadHeaderSize, offset := decodeVarInt(payload)
	pointer += offset

	// Read all of the type codes from the payload header
	var typeCodes []uint64
	for {
		typeCode, offset := decodeVarInt(payload[pointer:])
		pointer += offset

		typeCodes = append(typeCodes, typeCode)

		if uint64(pointer) >= payloadHeaderSize {
			break
		}
	}

	// The remainder of the payload is the actual record data, which we can now
	// read because we know the type codes
	var entries []RecordEntry
	for _, typeCode := range typeCodes {
		if typeCode == 0 {
			entries = append(entries, RecordEntry{RecordEntryType: 0, Null: lo.ToPtr(true)})
			pointer += 0
		} else if typeCode >= 1 && typeCode <= 6 {
			// Make an 8 byte empty slice
			result := make([]byte, 8)

			// Extract the correct number of bytes from the payload
			size := intTypeCodeByteSize(typeCode)
			value := payload[pointer : pointer+size]

			// Copy the bytes into the result slice so that we can decode them
			// as a uint64
			for i := 7; 7-i < len(value); i-- {
				result[i] = value[7-i]
			}

			entries = append(entries, RecordEntry{RecordEntryType: 1, Number: lo.ToPtr(binary.BigEndian.Uint64(result))})
			pointer += size
		} else if typeCode == 8 {
			entries = append(entries, RecordEntry{RecordEntryType: 1, Number: lo.ToPtr(uint64(0))})
		} else if typeCode == 9 {
			entries = append(entries, RecordEntry{RecordEntryType: 1, Number: lo.ToPtr(uint64(1))})
		} else if typeCode >= 12 && typeCode%2 == 1 {
			length := uint16((typeCode - 13) / 2)
			entries = append(entries, RecordEntry{RecordEntryType: 2, Text: lo.ToPtr(string(payload[pointer : pointer+length]))})
			pointer += length
		} else {
			panic(fmt.Sprintf("Unimplemented: type code decoder %v", typeCode))
		}
	}

	return Record{Entries: entries}
}

func intTypeCodeByteSize(typeCode uint64) uint16 {
	switch typeCode {
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 3
	case 4:
		return 4
	case 5:
		return 6
	case 6:
		return 8
	default:
		panic("Unreachable: invalid type code")
	}
}
