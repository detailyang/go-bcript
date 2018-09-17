package main

import (
	"flag"
	"fmt"

	"github.com/detailyang/go-bscript"
)

func main() {
	flag.Parse()
	args := flag.Args()
	code := args[0]
	script, err := bscript.NewScriptFromString(code)
	if err != nil {
		panic(err)
	}

	interpreter := bscript.NewInterpreter()
	err = interpreter.Eval(script, bscript.ScriptSkipDisabledOPCode)
	if err != nil {
		panic(err)
	}

	dstack := interpreter.GetDStack()
	fmt.Println(dstack)
}
