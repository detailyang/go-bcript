package bscript

import (
	"fmt"
	"strings"
)

type Disassembler struct {
}

func NewDisassembler() *Disassembler {
	return &Disassembler{}
}

func (d *Disassembler) Disassemble(script *Script, sep string) (string, error) {
	tmp := script.Pos
	script.Pos = 0
	defer func() {
		script.Pos = tmp
	}()

	rv := make([]string, 0, 64)
	for {
		ins, err := script.Next()
		if err != nil {
			if err == ErrScriptEOF {
				break
			}
			return "", err
		}

		rv = append(rv, fmt.Sprintf("%s", ins))
	}

	return strings.Join(rv, sep), nil
}
