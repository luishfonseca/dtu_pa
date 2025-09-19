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
