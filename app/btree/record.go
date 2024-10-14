package btree

import (
	"bytes"
	"encoding/binary"
)

type record[T any] struct {
	ResultConstructor resultConstructor[T]
	Data              []byte
}

func sizeForTypeCode(typeCode uint64) uint64 {
	typeCodeToSizeForInt := [7]uint64{1, 2, 3, 4, 6, 8, 8}

	if typeCode == 0 || typeCode == 8 || typeCode == 9 {
		return 0
	} else if typeCode >= 1 && typeCode <= 7 {
		return typeCodeToSizeForInt[typeCode-1]
	} else if typeCode >= 12 && typeCode%2 == 1 {
		return (typeCode - 12) / 2
	} else {
		panic("unsupported type code: not implemented")
	}
}

func (r record[T]) ReadColumn(column uint64) (T, error) {
	// If the column is an integer, then we need to use the type code
	pointer := uint64(0)

	// Read the header size varint
	headerSize, headerSizeVarintSize, err := decodeUvarint(r.Data)
	if err != nil {
		return r.ResultConstructor.Null(), err
	}
	pointer += headerSizeVarintSize

	// We need to loop over every column that comes before the one we want to
	// read so that we can determine both the offset in the body at which the
	// column we want is, as well as the type code of the column (which tell us
	// how to interpret it). columnOffset is the offset from the start of the
	// body (i.e the end of the header).
	typeCode, columnOffset, columnSize := uint64(0), uint64(0), uint64(0)
	for i := uint64(0); i <= column; i++ {
		code, varintSize, err := decodeUvarint(r.Data[pointer:])
		if err != nil {
			return r.ResultConstructor.Null(), err
		}
		pointer += varintSize

		typeCode = code
		columnSize = sizeForTypeCode(code)
		columnOffset += columnSize
	}

	// Extract the exact column data from the body
	columnData := r.Data[headerSize+columnOffset-columnSize : headerSize+columnOffset]

	// Interpret the column data based on the type code
	if typeCode == 0 {
		return r.ResultConstructor.Null(), nil
	} else if typeCode >= 1 && typeCode <= 6 {
		// We're going to read the column data as a big-endian 8 byte unsigned
		// integer, so we need to pad the data with 0s if it's not already 8
		// bytes
		padding := bytes.Repeat([]byte{0}, 8-int(columnSize))
		columnData = append(padding, columnData...)

		return r.ResultConstructor.Number(int64(binary.BigEndian.Uint64(columnData))), nil
	} else if typeCode == 8 {
		return r.ResultConstructor.Number(0), nil
	} else if typeCode == 9 {
		return r.ResultConstructor.Number(1), nil
	} else if typeCode >= 12 && typeCode%2 == 1 {
		return r.ResultConstructor.Text(string(columnData)), nil
	} else {
		panic("unsupported type code: not implemented")
	}
}
