package bscript

import (
	"bytes"
	"math/bits"
)

type StackElemnt []byte

func (v StackElemnt) Cat(n StackElemnt) {
	v = append(v, n...)
}

func (v StackElemnt) Invert() {
	for i := range v {
		v[i] = bits.Reverse8(v[i])
	}
}

func (v StackElemnt) Equal(n StackElemnt) bool {
	return bytes.Equal([]byte(v), []byte(n))
}

func (v StackElemnt) BitXor(n StackElemnt) {
	for i := range v {
		v[i] ^= n[i]
	}
}

func (v StackElemnt) BitOr(n StackElemnt) {
	for i := range v {
		v[i] |= n[i]
	}
}

func (v StackElemnt) BitAnd(n StackElemnt) {
	for i := range v {
		v[i] &= n[i]
	}
}

func (v StackElemnt) Bytes() []byte {
	return v
}

func (v StackElemnt) Number(required bool, limit int) (Number, error) {
	return NewNumberFromBytes(v.Bytes(), required, limit)
}

func (v StackElemnt) Boolean() Boolean {
	return NewBoolean(v.Bytes())
}

func (v StackElemnt) Size() int {
	return len(v)
}
