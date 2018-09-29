package bscript

import (
	"errors"
	"fmt"
)

var (
	ErrStackEmpty        = errors.New("Stack: no data")
	ErrStackNotEnough    = errors.New("Stack: not enough depth")
	ErrStackEraseInvalid = errors.New("Stack: erase invalid")
)

// Stack hold byte vectors.
// When used as numbers, byte vectors are interpreted as little-endian variable-length integers with the most significant bit determining the sign of the integer. Thus 0x81 represents -1. 0x80 is another representation of zero (so called negative 0). Positive 0 is represented by a null-length vector.
// Byte vectors are interpreted as Booleans where False is represented by any representation of zero and True is represented by any representation of non-zero.
type Stack struct {
	data []StackElemnt
}

func NewStack() *Stack {
	return &Stack{
		data: make([]StackElemnt, 0, 1024),
	}
}

func (s *Stack) CloneFrom(ns *Stack) {
	s.data = ns.data
}

func (s *Stack) Reverse() {
	for i, j := 0, len(s.data); i < j; i, j = i+1, j-1 {
		s.data[i], s.data[j] = s.data[j], s.data[i]
	}
}

func (s *Stack) Combine(ns *Stack) {
	s.data = append(s.data, ns.data...)
}

func (s *Stack) Push(data StackElemnt) {
	s.data = append(s.data, data)
}

func (s *Stack) Swap(i, j int) error {
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

func (s *Stack) Replace(n int, data StackElemnt) error {
	depth := len(s.data)
	if n > 0 {
		if n >= depth {
			return ErrStackNotEnough
		}
		s.data[n] = data
	} else {
		if -n >= depth {
			return ErrStackNotEnough
		}
		s.data[n+depth] = data
	}

	return nil
}

func (s *Stack) Erase(start, end int) error {
	for i := start; i < end; i++ {
		if err := s.Remove(i); err != nil {
			return err
		}
	}

	return nil
}

func (s *Stack) Remove(n int) error {
	depth := len(s.data)
	if n > 0 {
		if n >= depth {
			return ErrStackNotEnough
		}

		s.data = append(s.data[:n], s.data[n+1:]...)
	} else {
		if -n >= depth {
			return ErrStackNotEnough
		}

		s.data = append(s.data[:depth+n], s.data[depth+n+1:]...)
	}

	return nil
}

func (s *Stack) InsertBefore(n int, data StackElemnt) error {
	depth := len(s.data)
	if n > 0 {
		if n >= depth {
			return ErrStackNotEnough
		}

		behind := s.data[n:]
		s.data = append(append([]StackElemnt{data}, s.data[:n]...), behind...)
	} else {
		if -n >= depth {
			return ErrStackNotEnough
		}

		behind := s.data[depth+n:]
		s.data = append(append([]StackElemnt{data}, s.data[:depth+n]...), behind...)
	}

	return nil
}

func (s *Stack) InsertAfter(n int, data StackElemnt) error {
	depth := len(s.data)
	if n > 0 {
		if n >= depth {
			return ErrStackNotEnough
		}

		behind := s.data[n:]
		s.data = append(append(s.data[:n], data), behind...)
	} else {
		if -n >= depth {
			return ErrStackNotEnough
		}

		behind := s.data[depth+n:]
		s.data = append(append(s.data[:depth+n], data), behind...)
	}

	return nil
}

func (s *Stack) Pop() (StackElemnt, error) {
	if len(s.data) > 0 {
		data := s.data[len(s.data)-1]
		s.data = s.data[:len(s.data)-1]
		return data, nil
	}

	return nil, ErrStackEmpty
}

func (s *Stack) Peek(n int) (StackElemnt, error) {
	depth := s.Depth()
	if n >= 0 {
		if n >= depth {
			return nil, ErrStackNotEnough
		}
		return s.data[n], nil
	}

	if -n > depth {
		return nil, ErrStackNotEnough
	}

	return s.data[n+depth], nil
}

func (s *Stack) Iter(fn func(e StackElemnt)) {
	for _, elem := range s.data {
		fn(elem)
	}
}

func (s *Stack) Clean() {
	s.data = s.data[:0]
}

func (s *Stack) Depth() int {
	return len(s.data)
}

func (s *Stack) Clone() *Stack {
	data := make([]StackElemnt, len(s.data))
	copy(data, s.data)

	return &Stack{
		data: data,
	}
}

func (s *Stack) String() string {
	var result string

	if len(s.data) == 0 {
		result += " <empty> "
	}
	for _, stack := range s.data {
		if len(stack) == 0 {
			result += " <empty> "
		} else {
			result += fmt.Sprintf(" <%x> ", stack)
		}
	}

	return result
}
