package bscript

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestScript(t *testing.T) {
	script, err := NewScriptFromBytes([]byte{0x02, 0x10, 0x00})
	if err != nil {
		t.Error(err)
	}

	for {
		ins, err := script.Next()
		if err != nil {
			if err == ErrScriptEOF {
				break
			}
			t.Error(err)
		}
		if ins.OPCode != OP_PUSHBYTES_2 {
			t.Error("expect OP_PUSHBYTES_2")
		}
		if ins.Step != 3 {
			t.Error("expect 3")
		}
		if ins.Data[0] != 16 && ins.Data[1] != 0 {
			t.Error("expect {16, 0")
		}
	}
}

func TestScriptFromBytes(t *testing.T) {
	b, err := hex.DecodeString("0101010293")
	if err != nil {
		t.Error(err)
	}

	script, err := NewScriptFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	instructions := make([]*Instruction, 0)
	for {
		ins, err := script.Next()
		if err != nil {
			if err == ErrScriptEOF {
				break
			}
			t.Error(err)
			break
		}
		instructions = append(instructions, ins)
	}

	if instructions[0].OPCode != OP_PUSHBYTES_1 {
		t.Error("expect OP_PUSHBYTES_1")
	}
	if instructions[1].OPCode != OP_PUSHBYTES_1 {
		t.Error("expect OP_PUSHBYTES_1")
	}
	if instructions[2].OPCode != OP_ADD {
		t.Error("expect OP_ADD")
	}
}

func TestScriptFromString(t *testing.T) {
	code := `1 2 OP_ADD`
	s, err := NewScriptFromString(code)
	if err != nil {
		t.Error(err)
	}

	if !(bytes.Equal(s.Bytes(), []byte{0x01, 0x01, 0x01, 0x02, 0x93})) {
		t.Error("expect []byte{0x01, 0x01, 0x01, 0x02, 0x93}")
	}
}
