package main

import (
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/detailyang/go-bscript"
)

func main() {
	flag.Parse()
	args := flag.Args()
	hexcode := args[0]
	code, err := hex.DecodeString(hexcode)
	if err != nil {
		panic(err)
	}

	script, err := bscript.NewScriptFromBytes([]byte(code))
	if err != nil {
		panic(err)
	}

	dissembler := bscript.NewDisassembler()
	s, err := dissembler.Disassemble(script)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", s)
}
