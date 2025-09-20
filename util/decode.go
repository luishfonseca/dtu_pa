package util

import (
	"encoding/binary"
	"fmt"
)

func Decode(buf []byte, n any) error {
	if len(buf) != binary.Size(n) {
		return fmt.Errorf("buffer size %d does not match type size %d", len(buf), binary.Size(n))
	}

	if nrd, err := binary.Decode(buf, binary.BigEndian, n); err != nil {
		return err
	} else if nrd != binary.Size(n) {
		return fmt.Errorf("binary decode read %d bytes, expected %d", nrd, binary.Size(n))
	}

	return nil
}

func Encode(n any) ([]byte, error) {
	buf := make([]byte, binary.Size(n))

	if nwr, err := binary.Encode(buf, binary.BigEndian, n); err != nil {
		return nil, err
	} else if nwr != binary.Size(n) {
		return nil, fmt.Errorf("binary encode wrote %d bytes, expected %d", nwr, binary.Size(n))
	}

	return buf, nil
}
