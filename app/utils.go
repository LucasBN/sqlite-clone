package main

type DatabaseHeader struct {
	HeaderString     [16]byte
	PageSize         uint16
	FileWriteVersion uint8
	FileReadVersion  uint8
	ReservedSpace    uint8
	Middle           [38]byte
	TextEncoding     uint32
	End              [40]byte
}

// A varint consists of either zero or more bytes which have the high-order
// bit set followed by a single byte with the high-order bit clear, or nine
// bytes, whichever is shorter.
func decodeVarInt(data []byte) (uint64, uint16) {
	var value uint64
	for i := 0; i < 8; i++ {
		value = (value << 7) | uint64(data[i]&0x7F)
		if data[i]&0x80 == 0 {
			return value, uint16(i + 1)
		}
	}
	return value<<8 | uint64(data[8]), 9
}
