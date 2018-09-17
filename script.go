package bscript

import (
	"encoding/binary"
	"errors"
)

var (
	ErrScriptEOF            = errors.New("script: EOF")
	ErrScriptBadInstruction = errors.New("script: bad instruction")
	ErrScriptTakeOverflow   = errors.New("script: take overflow")
	ErrScriptBadTypeCast    = errors.New("script: bad type cast")
)

type Script struct {
	Data []byte
	Pos  int
}

func NewScript() *Script {
	return &Script{
		Data: make([]byte, 0, 1024),
		Pos:  0,
	}
}

func NewScriptFromBytes(src []byte) (*Script, error) {
	return &Script{
		Data: src,
		Pos:  0,
	}, nil
}

func NewScriptFromString(src string) (*Script, error) {
	lexer := NewLexer(src)
	script := NewScript()

	for {
		tok, err := lexer.Scan()
		if err != nil {
			if err == ErrLexerReachEOF {
				break
			}
			return nil, err
		}

		switch tok.kind {
		case TOKEN_CODE:
			opcode, ok := tok.value.(OPCode)
			if !ok {
				return nil, ErrScriptBadTypeCast
			}
			script.AddOPCode(opcode)
		case TOKEN_NUMBER:
			value, ok := tok.value.(int64)
			if !ok {
				return nil, ErrScriptBadTypeCast
			}
			script.AddInt64(value)
		case TOKEN_HEXSTRING:
			value, ok := tok.value.([]byte)
			if !ok {
				return nil, ErrScriptBadTypeCast
			}
			script.AddBytes(value)
		}
	}

	return script, nil
}

func (s *Script) AddBytes(b []byte) *Script {
	n := len(b)
	switch {
	case n < 2^8:
		s.AddOPCode(OP_PUSHDATA1)
	case n < 2^16:
		s.AddOPCode(OP_PUSHDATA2)
	case n < 2^32:
		s.AddOPCode(OP_PUSHDATA4)
	default:
		panic("too many bytes")
	}

	s.Data = append(s.Data, b...)

	return s
}

func (s *Script) AddOPCode(opcode OPCode) *Script {
	s.Data = append(s.Data, byte(opcode))
	return s
}

func (s *Script) AddInt64(n int64) *Script {
	s.AddOPCode(OP_PUSHBYTES_8)
	binary.LittleEndian.PutUint64(s.Data, uint64(n))
	return s
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

		s.Pos += 1

		return &Instruction{
			OPCode: opcode,
			Step:   1,
			Data:   make([]byte, 0),
		}, nil
	}

	return nil, ErrScriptBadInstruction
}

func (s *Script) Bytes() []byte {
	return s.Data
}

func (s *Script) Size() int {
	return len(s.Data)
}

func (s *Script) String() string {
	return ""
}
