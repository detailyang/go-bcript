package bscript

import (
	"encoding/binary"
	"errors"
	"math"
)

var (
	ErrScriptEOF              = errors.New("script: EOF")
	ErrScriptBadInstruction   = errors.New("script: bad instruction")
	ErrScriptTakeOverflow     = errors.New("script: take overflow")
	ErrScriptBadTypeCast      = errors.New("script: bad type cast")
	ErrScriptPushSizeOverflow = errors.New("script: push size overflow")
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

	needPushSize := true

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

			if OP_PUSHBYTES_1 <= opcode && opcode <= OP_PUSHDATA4 {
				needPushSize = false
			} else {
				needPushSize = true
			}
		case TOKEN_NUMBER:
			value, ok := tok.value.(int64)
			if !ok {
				return nil, ErrScriptBadTypeCast
			}

			var n Number
			var size int
			if value != 0 {
				n = Number(value)
				size = len(n.Bytes())
			} else {
				n = Number(0)
				size = 1
			}

			if needPushSize {
				switch {
				case size <= math.MaxInt8:
					script.AddOPCode(OP_PUSHDATA1)
					script.AddBytes([]byte{byte(size)})
				case size <= math.MaxInt16:
					script.AddOPCode(OP_PUSHDATA2)
					buf := make([]byte, 2)
					binary.LittleEndian.PutUint16(buf, uint16(size))
					script.AddBytes(buf)
				case size <= math.MaxInt32:
					script.AddOPCode(OP_PUSHDATA4)
					buf := make([]byte, 4)
					binary.LittleEndian.PutUint32(buf, uint32(size))
					script.AddBytes(buf)
				default:
					return nil, ErrScriptPushSizeOverflow
				}
			}

			if value != 0 {
				script.AddNumber(n)
			} else {
				script.AddBytes([]byte{0})
			}
		case TOKEN_HEXSTRING:
			value, ok := tok.value.([]byte)
			if !ok {
				return nil, ErrScriptBadTypeCast
			}

			size := len(value)

			if needPushSize {
				switch {
				case size <= math.MaxInt8:
					script.AddOPCode(OP_PUSHDATA1)
					script.AddBytes([]byte{byte(size)})
				case size <= math.MaxInt16:
					script.AddOPCode(OP_PUSHDATA2)
					buf := make([]byte, 2)
					binary.LittleEndian.PutUint16(buf, uint16(size))
					script.AddBytes(buf)
				case size <= math.MaxInt32:
					script.AddOPCode(OP_PUSHDATA4)
					buf := make([]byte, 4)
					binary.LittleEndian.PutUint32(buf, uint32(size))
					script.AddBytes(buf)
				default:
					return nil, ErrScriptPushSizeOverflow
				}
			}

			script.AddBytes(value)
		}
	}

	return script, nil
}

func (s *Script) AddBytes(b []byte) *Script {
	s.Data = append(s.Data, b...)
	return s
}

func (s *Script) AddOPCode(opcode OPCode) *Script {
	s.Data = append(s.Data, byte(opcode))
	return s
}

func (s *Script) AddNumber(n Number) *Script {
	s.Data = append(s.Data, n.Bytes()...)
	return s
}

func (s *Script) takePushBytes(offset, n int) ([]byte, error) {
	if offset+n >= len(s.Data) {
		return nil, ErrScriptTakeOverflow
	}

	nbytes := 0
	for i := offset; i < offset+n; i++ {
		nbytes |= int(s.Data[i]) << uint8(8*(i-offset))
	}

	if offset+n+nbytes > len(s.Data) {
		return nil, ErrScriptTakeOverflow
	}

	return s.Data[offset+n : offset+n+nbytes], nil
}

func (s *Script) takeBytes(offset, n int) ([]byte, error) {
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

		data, err := s.takePushBytes(s.Pos+1, nbytes)
		if err != nil {
			return nil, err
		}

		step := nbytes + len(data) + 1
		s.Pos += step

		return &Instruction{
			Data:   data,
			OPCode: opcode,
			Step:   step,
		}, nil

	default:
		if OP_PUSHBYTES_1 <= opcode && opcode <= OP_PUSHBYTES_75 {
			data, err := s.takeBytes(s.Pos+1, int(opcode))
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
