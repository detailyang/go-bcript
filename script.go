package bscript

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
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

func NewScriptFromHexString(hexstring string) (*Script, error) {
	b, err := hex.DecodeString(hexstring)
	if err != nil {
		return nil, err
	}

	return NewScriptFromBytes(b), nil
}

func NewScriptFromBytes(src []byte) *Script {
	return &Script{
		Data: src,
		Pos:  0,
	}
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
			script.PushOPCode(opcode)

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
					script.PushOPCode(OP_PUSHDATA1)
					script.PushBytes([]byte{byte(size)})
				case size <= math.MaxInt16:
					script.PushOPCode(OP_PUSHDATA2)
					buf := make([]byte, 2)
					binary.LittleEndian.PutUint16(buf, uint16(size))
					script.PushBytes(buf)
				case size <= math.MaxInt32:
					script.PushOPCode(OP_PUSHDATA4)
					buf := make([]byte, 4)
					binary.LittleEndian.PutUint32(buf, uint32(size))
					script.PushBytes(buf)
				default:
					return nil, ErrScriptPushSizeOverflow
				}
			}

			if value != 0 {
				script.PushNumber(n)
			} else {
				script.PushBytes([]byte{0})
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
					script.PushOPCode(OP_PUSHDATA1)
					script.PushBytes([]byte{byte(size)})
				case size <= math.MaxInt16:
					script.PushOPCode(OP_PUSHDATA2)
					buf := make([]byte, 2)
					binary.LittleEndian.PutUint16(buf, uint16(size))
					script.PushBytes(buf)
				case size <= math.MaxInt32:
					script.PushOPCode(OP_PUSHDATA4)
					buf := make([]byte, 4)
					binary.LittleEndian.PutUint32(buf, uint32(size))
					script.PushBytes(buf)
				default:
					return nil, ErrScriptPushSizeOverflow
				}
			}

			script.PushBytes(value)
		}
	}

	return script, nil
}

func (s Script) IsPushOnly() bool {
	for {
		ins, err := s.Next()
		if err != nil {
			break
		}

		if ins.OPCode > OP_16 {
			return false
		}
	}

	return true
}

func (s Script) IsPayToScriptHash() bool {
	return len(s.Data) == 23 &&
		s.Data[0] == byte(OP_HASH160) &&
		s.Data[1] == byte(OP_PUSHBYTES_20) &&
		s.Data[22] == byte(OP_EQUAL)
}

// A scriptPubKey (or redeemScript as defined in BIP16/P2SH) that consists of a 1-byte push opcode (for 0 to 16) followed by a data push between 2 and 40 bytes gets a new special meaning. The value of the first push is called the "version byte". The following byte vector pushed is called the "witness program".
func (s Script) ParseWitnessProgram() (uint8, []byte, bool) {
	if len(s.Data) < 4 || len(s.Data) > 42 || len(s.Data) != int(s.Data[1])+2 {
		return 0, nil, false
	}

	var version uint8
	if OPCode(s.Data[0]) == 0 {
	} else if OP_1 <= OPCode(s.Data[0]) && OPCode(s.Data[0]) <= OP_16 {
		version = s.Data[0] - uint8(OP_1) + 1
	} else {
		return 0, nil, false
	}

	return version, s.Data[2:], true
}

func (s *Script) PushBytesWithOP(b []byte) *Script {
	size := len(b)

	switch {
	case size <= int(OP_PUSHDATA1):
		s.PushBytes([]byte{byte(size)})
	case size <= math.MaxInt8:
		s.PushOPCode(OP_PUSHDATA1)
		s.PushBytes([]byte{byte(size)})
	case size <= math.MaxInt16:
		s.PushOPCode(OP_PUSHDATA2)
		buf := make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, uint16(size))
		s.PushBytes(buf)
	case size <= math.MaxInt32:
		s.PushOPCode(OP_PUSHDATA4)
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(size))
		s.PushBytes(buf)
	default:
		panic("Cannot push more than 2^32 bytes to script")
	}

	s.PushBytes(b)

	return s
}

func (s *Script) PushInstruction(ins *Instruction) *Script {
	s.Data = append(s.Data, byte(ins.OPCode))
	s.Data = append(s.Data, ins.Data...)
	return s
}

func (s *Script) PushBytes(b []byte) *Script {
	s.Data = append(s.Data, b...)
	return s
}

func (s *Script) PushOPCode(opcode OPCode) *Script {
	s.Data = append(s.Data, byte(opcode))
	return s
}

func (s *Script) PushInt64(n int64) *Script {
	if n == -1 || n >= 1 && n <= 16 {
		s.Data = append(s.Data, byte(n)+byte(OP_1)-1)
	} else if n == 0 {
		s.Data = append(s.Data, byte(OP_0))
	} else {
		return s.PushNumber(Number(n))
	}

	return s
}

func (s *Script) PushNumber(n Number) *Script {
	s.Data = append(s.Data, n.Bytes()...)
	return s
}

func (s *Script) Filter(fs *Script) *Script {
	pos := 0
	size := len(fs.Data)
	end := len(s.Data)
	rv := make([]byte, 0, end)

	if size > end || size == 0 {
		return NewScriptFromBytes(s.Bytes())
	}

	for pos < end-size {
		if bytes.Equal(s.Data[pos:pos+size], fs.Data) {
			pos += size
		} else {
			rv = append(rv, s.Data[pos])
			pos += 1
		}
	}

	return NewScriptFromBytes(rv)
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

func (s *Script) SubScript(from int) (*Script, error) {
	if from >= len(s.Data) {
		return nil, ErrScriptTakeOverflow
	}

	return NewScriptFromBytes(s.Data[from:]), nil
}

func (s Script) WithoutSep() *Script {
	ns := NewScript()

	for {
		ins, err := s.Next()
		if err != nil {
			break
		}

		opcode := ins.OPCode
		if opcode == OP_CODESEPARATOR {
			continue
		}

		ns.PushInstruction(ins)
	}

	return ns
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

func (s *Script) Disassemble(sep string) string {
	dis, _ := NewDisassembler().Disassemble(s, sep)
	return dis
}

func (s *Script) Reset() {
	s.Pos = 0
}

func (s *Script) Bytes() []byte {
	buf := make([]byte, len(s.Data))
	copy(buf, s.Data)
	return buf
}

func (s *Script) Size() int {
	return len(s.Data)
}

func (s *Script) Hex() string {
	return hex.EncodeToString(s.Data[s.Pos:])
}

func (s *Script) String() string {
	return s.Hex()
}
