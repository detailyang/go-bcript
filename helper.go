package bscript

import (
	"errors"
)

var (
	ErrHelperReadOverflow = errors.New("helper: read nbytes overflow")
)

func readNBytes(data []byte, nbytes int) (int, error) {
	if len(data) < nbytes {
		return 0, ErrHelperReadOverflow
	}

	n := 0
	for i := 0; i < nbytes; i++ {
		n |= int(data[i]) << uint8(8*i)
	}

	return n, nil
}
