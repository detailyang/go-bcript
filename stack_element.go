package bscript

import (
	"bytes"
	"math/bits"
)

type valtype []byte

func (v valtype) Cat(n valtype) {
	v = append(v, n...)
}

func (v valtype) Invert() {
	for i := range v {
		v[i] = bits.Reverse8(v[i])
	}
}

func (v valtype) Equal(n valtype) bool {
	return bytes.Equal([]byte(v), []byte(n))
}

func (v valtype) BitXor(n valtype) {
	for i := range v {
		v[i] ^= n[i]
	}
}

func (v valtype) BitOr(n valtype) {
	for i := range v {
		v[i] |= n[i]
	}
}

func (v valtype) BitAnd(n valtype) {
	for i := range v {
		v[i] &= n[i]
	}
}

func (v valtype) Bytes() []byte {
	return v
}

func (v valtype) ToNumber(required bool, limit int) (Number, error) {
	return NewNumberFromBytes(v.Bytes(), required, limit)
}

func (v valtype) ToBoolean() Boolean {
	return NewBoolean(v.Bytes())
}

func (v valtype) Size() int {
	return len(v)
}
