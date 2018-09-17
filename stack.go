package bscript

import (
	"bytes"
	"errors"
	"math/bits"
)

var (
	ErrStackEmpty        = errors.New("stack: no data")
	ErrStackNotEnough    = errors.New("stack: not enough depth")
	ErrStackEraseInvalid = errors.New("stack: erase invalid")
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

func (v valtype) ToNumber() (Number, error) {
	return NewNumber(v.Bytes())
}

func (v valtype) ToBoolean() Boolean {
	return NewBoolean(v.Bytes())
}

func (v valtype) Size() int {
	return len(v)
}

// stack hold byte vectors.
// When used as numbers, byte vectors are interpreted as little-endian variable-length integers with the most significant bit determining the sign of the integer. Thus 0x81 represents -1. 0x80 is another representation of zero (so called negative 0). Positive 0 is represented by a null-length vector.
// Byte vectors are interpreted as Booleans where False is represented by any representation of zero and True is represented by any representation of non-zero.
type stack struct {
	stack []valtype
}

func NewStack() *stack {
	return &stack{
		stack: make([]valtype, 0, 1024),
	}
}

func (s *stack) Copy() *stack {
	return &stack{
		stack: s.stack,
	}
}

func (s *stack) Reverse() {
	for i, j := 0, len(s.stack); i < j; i, j = i+1, j-1 {
		s.stack[i], s.stack[j] = s.stack[j], s.stack[i]
	}
}

func (s *stack) Combine(ns *stack) {
	s.stack = append(s.stack, ns.stack...)
}

func (s *stack) Push(data valtype) {
	s.stack = append(s.stack, data)
}

func (s *stack) Swap(i, j int) error {
	d1, err := s.Peek(i)
	if err != nil {
		return err
	}

	d2, err := s.Peek(j)
	if err != nil {
		return err
	}

	if err := s.Replace(i, d2); err != nil {
		return err
	}
	if err := s.Replace(j, d1); err != nil {
		return err
	}

	return nil
}

func (s *stack) Replace(n int, data valtype) error {
	depth := len(s.stack)
	if n > 0 {
		if n >= depth {
			return ErrStackNotEnough
		}
		s.stack[n] = data
	} else {
		if -n >= depth {
			return ErrStackNotEnough
		}
		s.stack[n+depth] = data
	}

	return nil
}

func (s *stack) Erase(start, end int) error {
	for i := start; i < end; i++ {
		if err := s.Remove(i); err != nil {
			return err
		}
	}

	return nil
}

func (s *stack) Remove(n int) error {
	depth := len(s.stack)
	if n > 0 {
		if n >= depth {
			return ErrStackNotEnough
		}

		s.stack = append(s.stack[:n], s.stack[n+1:]...)
	} else {
		if -n >= depth {
			return ErrStackNotEnough
		}

		s.stack = append(s.stack[:depth+n], s.stack[depth+n+1:]...)
	}

	return nil
}

func (s *stack) InsertBefore(n int, data valtype) error {
	depth := len(s.stack)
	if n > 0 {
		if n >= depth {
			return ErrStackNotEnough
		}

		behind := s.stack[n:]
		s.stack = append(append([]valtype{data}, s.stack[:n]...), behind...)
	} else {
		if -n >= depth {
			return ErrStackNotEnough
		}

		behind := s.stack[depth+n:]
		s.stack = append(append([]valtype{data}, s.stack[:depth+n]...), behind...)
	}

	return nil
}

func (s *stack) InsertAfter(n int, data valtype) error {
	depth := len(s.stack)
	if n > 0 {
		if n >= depth {
			return ErrStackNotEnough
		}

		behind := s.stack[n:]
		s.stack = append(append(s.stack[:n], data), behind...)
	} else {
		if -n >= depth {
			return ErrStackNotEnough
		}

		behind := s.stack[depth+n:]
		s.stack = append(append(s.stack[:depth+n], data), behind...)
	}

	return nil
}

func (s *stack) Pop() (valtype, error) {
	if len(s.stack) > 0 {
		data := s.stack[len(s.stack)-1]
		s.stack = s.stack[:len(s.stack)-1]
		return data, nil
	}

	return nil, ErrStackEmpty
}

func (s *stack) PeekNumber(n int) (Number, error) {
	d, err := s.Peek(n)
	if err != nil {
		return 0, err
	}

	return d.ToNumber()
}

func (s *stack) Peek(n int) (valtype, error) {
	depth := s.Depth()
	if n >= 0 {
		if n >= depth {
			return nil, ErrStackNotEnough
		}
		return s.stack[n], nil
	}

	if -n > depth {
		return nil, ErrStackNotEnough
	}

	return s.stack[n+depth], nil
}

func (s *stack) Clean() {
	s.stack = s.stack[:0]
}

func (s *stack) Depth() int {
	return len(s.stack)
}
