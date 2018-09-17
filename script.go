package bscript

import (
	"errors"
)

var (
	ErrScriptEOF            = errors.New("script: EOF")
	ErrScriptBadInstruction = errors.New("script: bad instruction")
	ErrScriptTakeOverflow   = errors.New("script: take overflow")
)

type Script struct {
	Data []byte
	Pos  int
}

func NewScriptFromBytes(src []byte) (*Script, error) {
	return &Script{
		Data: src,
		Pos:  0,
	}, nil
}

func NewScriptFromString(src string) (*Script, error) {
	return &Script{}, nil
}

func (s *Script) take(offset, n int) ([]byte, error) {
	if offset+n > len(s.Data) {
		return nil, ErrScriptTakeOverflow
	}

	return s.Data[offset : offset+n], nil
}

func (s *Script) Next() (*Instruction, error) {
	if s.Pos >= len(s.Data) {
		return nil, ErrScriptEOF
	}

	u8 := s.Data[s.Pos]
	opcode, err := NewOPCode(u8)
	if err != nil {
		return nil, err
	}

	switch opcode {
	case OP_PUSHDATA1:
		fallthrough
	case OP_PUSHDATA2:
		fallthrough
	case OP_PUSHDATA4:
		nbytes := 1
		if opcode == OP_PUSHDATA2 {
			nbytes = 2
		} else if opcode == OP_PUSHDATA4 {
			nbytes = 4
		}

		slice, err := s.take(s.Pos+1, nbytes)
		if err != nil {
			return nil, err
		}

		n, err := readNBytes(slice, nbytes)
		if err != nil {
			return nil, err
		}

		data, err := s.take(s.Pos+1+nbytes, n)
		if err != nil {
			return nil, err
		}

		step := nbytes + n + 1
		s.Pos += step

		return &Instruction{
			Data:   data,
			OPCode: opcode,
			Step:   step,
		}, nil

	default:
		if OP_0 <= opcode && opcode <= OP_PUSHBYTES_75 {
			data, err := s.take(s.Pos+1, int(opcode))
			if err != nil {
				return nil, err
			}

			step := int(opcode) + 1
			s.Pos += step

			return &Instruction{
				OPCode: opcode,
				Step:   step,
				Data:   data,
			}, nil
		}
	}

	return nil, ErrScriptBadInstruction
}

func (s *Script) Size() int {
	return len(s.Data)
}

func (s *Script) String() string {
	return ""
}
