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

	script.Pos = tmp

	return strings.Join(rv, sep), nil
}
