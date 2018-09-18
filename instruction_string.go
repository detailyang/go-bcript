package bscript

func instructionCAT(ctx *InterpreterContext) error {
	d1, err := ctx.i.dstack.Peek(-2)
	if err != nil {
		return err
	}

	d2, err := ctx.i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	d1.Cat(d2)

	ctx.i.dstack.Pop()

	return nil
}

func instructionSUBSTR(ctx *InterpreterContext) error {
	return nil
}

func instructionLEFT(ctx *InterpreterContext) error {
	return nil
}

func instructionRIGHT(ctx *InterpreterContext) error {
	return nil
}

func instructionSIZE(ctx *InterpreterContext) error {
	d, err := ctx.i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	ctx.i.dstack.Push(Number(d.Size()).Bytes())

	return nil
}
