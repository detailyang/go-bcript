package bscript

import (
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
