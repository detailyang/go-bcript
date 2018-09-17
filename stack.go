package bscript

import (
	"errors"
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
	data []valtype
}

func NewStack() *Stack {
	return &Stack{
		data: make([]valtype, 0, 1024),
	}
}

func (s *Stack) Reverse() {
	for i, j := 0, len(s.data); i < j; i, j = i+1, j-1 {
		s.data[i], s.data[j] = s.data[j], s.data[i]
	}
}

func (s *Stack) Combine(ns *Stack) {
	s.data = append(s.data, ns.data...)
}

func (s *Stack) Push(data valtype) {
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

func (s *Stack) Replace(n int, data valtype) error {
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

func (s *Stack) InsertBefore(n int, data valtype) error {
	depth := len(s.data)
	if n > 0 {
		if n >= depth {
			return ErrStackNotEnough
		}

		behind := s.data[n:]
		s.data = append(append([]valtype{data}, s.data[:n]...), behind...)
	} else {
		if -n >= depth {
			return ErrStackNotEnough
		}

		behind := s.data[depth+n:]
		s.data = append(append([]valtype{data}, s.data[:depth+n]...), behind...)
	}

	return nil
}

func (s *Stack) InsertAfter(n int, data valtype) error {
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

func (s *Stack) Pop() (valtype, error) {
	if len(s.data) > 0 {
		data := s.data[len(s.data)-1]
		s.data = s.data[:len(s.data)-1]
		return data, nil
	}

	return nil, ErrStackEmpty
}

func (s *Stack) Peek(n int) (valtype, error) {
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

func (s *Stack) Clean() {
	s.data = s.data[:0]
}

func (s *Stack) Depth() int {
	return len(s.data)
}
