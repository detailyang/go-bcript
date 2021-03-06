package main

import (
	"flag"

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

	flag := bscript.ScriptSkipDisabledOPCode | bscript.ScriptEnableTrace
	interpreter := bscript.NewInterpreter()
	err = interpreter.Eval(script, flag, bscript.NewNoopChecker(), bscript.SignatureVersionBase)
	if err != nil {
		panic(err)
	}

	interpreter.PrintTraces()
}
