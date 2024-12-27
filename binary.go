package pcktool

import (
	"encoding/binary"
	"io"
)

func ReadByte(r io.Reader) (byte, error) {
	var b int8
	if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
		return 0, err
	}
	return byte(b), nil
}

func ReadInt(r io.Reader) (int, error) {
	var b int32
	if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
		return 0, err
	}
	return int(b), nil
}

func ReadN(r io.Reader, n int) ([]byte, error) {
	buf := make([]byte, n)
	if err := binary.Read(r, binary.LittleEndian, &buf); err != nil {
		return nil, err
	}
	return buf, nil
}
