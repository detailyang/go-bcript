package bscript

import (
	"errors"
	"math"
)

var (
	ErrNumberNonMinimalEncode = errors.New("number: is non-minimally encoded")
	ErrNumberOverflow         = errors.New("number: overflow")
)

const (
	NumberDefaultElementSize = 4
)

// Number are interpreted as little-endian variable-length integers
// with the most significant bit determining the sign of the integer
// Numeric opcodes (OP_1ADD, etc) are restricted to operating on 4-byte
// integers. The semantics are subtle, though: operands must be in the range
// [-2^31 +1...2^31 -1], but results may overflow (and are valid as long as
// they are not used in a subsequent numeric operation). CScriptNum enforces
// those semantics by storing results as an int64 and allowing out-of-range
// values to be returned as a vector of bytes but throwing an exception if
// arithmetic is done or the result is interpreted as an integer.
type Number int64

func NewNumber(d []byte) Number {
	if len(d) == 0 {
		return 0
	}

	var rv int64
	for i, val := range d {
		rv |= int64(val) << uint8(8*i)
	}

	if d[len(d)-1]&0x80 != 0 {
		// The maximum length of v has already been determined to be 4
		// above, so uint8 is enough to cover the max possible shift
		// value of 24.
		rv &= ^(int64(0x80) << uint8(8*(len(d)-1)))
		return Number(-rv)
	}

	return Number(rv)
}

func NewNumberFromBytes(d []byte, required bool, limit int) (Number, error) {
	if len(d) > limit {
		return 0, ErrNumberOverflow
	}

	if required && !isMinimallyEncode(d) {
		return 0, ErrNumberNonMinimalEncode
	}

	return NewNumber(d), nil
}

func (n Number) IsNegative() bool {
	return n < 0
}

func (n Number) IsMinimallyEncode(limit int) error {
	b := n.Bytes()
	if len(b) > limit {
		return ErrNumberOverflow
	}

	if ok := isMinimallyEncode(b); !ok {
		return ErrNumberNonMinimalEncode
	}

	return nil
}

func isMinimallyEncode(d []byte) bool {
	if len(d) == 0 {
		return true
	}

	// Check that the number is encoded with the minimum possible number
	// of bytes.
	//
	// If the most-significant-byte - excluding the sign bit - is zero
	// then we're not minimal. Note how this test also rejects the
	// negative-zero encoding, 0x80.
	if d[len(d)-1]&0x7F == 0 {
		if len(d) == 1 {
			return false
		}

		if d[len(d)-2]&0x80 == 0 {
			return false
		}
	}

	return true
}

// Bytes -- returns the number serialized as a little endian with a sign bit
// Example encodings:
//       127 -> [0x7f]
//      -127 -> [0xff]
//       128 -> [0x80 0x00]
//      -128 -> [0x80 0x80]
//       129 -> [0x81 0x00]
//      -129 -> [0x81 0x80]
//       256 -> [0x00 0x01]
//      -256 -> [0x00 0x81]
//     32767 -> [0xff 0x7f]
//    -32767 -> [0xff 0xff]
//     32768 -> [0x00 0x80 0x00]
//    -32768 -> [0x00 0x80 0x80]
func (n Number) Bytes() []byte {
	if n == 0 {
		return []byte{0x00}
	}

	absn := n
	sigbyte := 0x00
	if n < 0 {
		absn = -n
		sigbyte = 0x80
	}

	rv := make([]byte, 0, 9)
	for absn != 0 {
		rv = append(rv, byte(absn&0x000000FF))
		absn >>= 8
	}

	// - If the most significant byte is >= 0x80 and the value is positive,
	// push a new zero-byte to make the significant byte < 0x80 again.
	// - If the most significant byte is >= 0x80 and the value is negative,
	// push a new 0x80 byte that will be popped off when converting to an
	// integral.
	// - If the most significant byte is < 0x80 and the value is negative,
	// add 0x80 to it, since it will be subtracted and interpreted as a
	// negative when converting to an integral.
	last := rv[len(rv)-1]
	if (last & 0x80) > 0 {
		rv = append(rv, byte(sigbyte))
	} else if n < 0 {
		rv[len(rv)-1] |= 0x80
	}

	return rv
}

func (n Number) Int32() int32 {
	if n > math.MaxInt32 {
		return math.MaxInt32
	}

	if n < math.MinInt32 {
		return math.MinInt32
	}

	return int32(n)
}
