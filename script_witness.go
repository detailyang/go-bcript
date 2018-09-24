package bscript

import (
	"bytes"

	. "github.com/detailyang/go-bprimitives"
)

type ScriptWitness [][]byte

func NewScriptWitness(b [][]byte) ScriptWitness {
	witness := ScriptWitness(b)
	return witness
}

func NewScriptWitnessFromBuffer(buffer *Buffer) (ScriptWitness, error) {
	n, err := buffer.GetVarInt()
	if err != nil {
		return nil, err
	}

	witness := make([][]byte, n)
	for i := 0; i < int(n); i++ {
		b, err := buffer.GetVarBytes()
		if err != nil {
			return nil, err
		}
		witness[i] = b
	}

	return NewScriptWitness(witness), nil
}

func (s ScriptWitness) Equal(t ScriptWitness) bool {
	if s.Size() != t.Size() {
		return false
	}

	for i, _ := range s {
		if !bytes.Equal(s[i], t[i]) {
			return false
		}
	}

	return true
}

func (s ScriptWitness) Clone() ScriptWitness {
	b := make([][]byte, len(s))
	for i := 0; i < len(s); i++ {
		copy(b[i], s[i])
	}

	return NewScriptWitness(b)
}

func (s ScriptWitness) Size() int { return len(s) }

func (s ScriptWitness) Bytes() []byte {
	n := s.Size()
	buffer := NewBuffer().PutVarInt(uint64(n))

	for i := 0; i < n; i++ {
		buffer.PutVarBytes(s[i])
	}

	return buffer.Bytes()
}
