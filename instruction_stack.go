package bscript

func instructionTOTALSTACK(ctx *InterpreterContext) error {
	data, err := ctx.i.dstack.Pop()
	if err != nil {
		return err
	}
	ctx.i.astack.Push(data)

	return nil
}

func instructionFROMALTSTACK(ctx *InterpreterContext) error {
	data, err := ctx.i.astack.Pop()
	if err != nil {
		return err
	}
	ctx.i.dstack.Push(data)

	return nil
}

func instruction2DROP(ctx *InterpreterContext) error {
	if ctx.i.dstack.Depth() < 2 {
		return ErrInterpreterStackSizeNotEnough
	}

	ctx.i.dstack.Pop()
	ctx.i.dstack.Pop()

	return nil
}

func instruction2DUP(ctx *InterpreterContext) error {
	d1, err := ctx.i.dstack.Peek(-2)
	if err != nil {
		return err
	}
	d2, err := ctx.i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	ctx.i.dstack.Push(d1)
	ctx.i.dstack.Push(d2)

	return nil
}

func instruction3DUP(ctx *InterpreterContext) error {
	d1, err := ctx.i.dstack.Peek(-3)
	if err != nil {
		return err
	}
	d2, err := ctx.i.dstack.Peek(-2)
	if err != nil {
		return err
	}
	d3, err := ctx.i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	ctx.i.dstack.Push(d1)
	ctx.i.dstack.Push(d2)
	ctx.i.dstack.Push(d3)

	return nil
}

func instruction2OVER(ctx *InterpreterContext) error {
	d1, err := ctx.i.dstack.Peek(-4)
	if err != nil {
		return err
	}
	d2, err := ctx.i.dstack.Peek(-3)
	if err != nil {
		return err
	}

	ctx.i.dstack.Push(d1)
	ctx.i.dstack.Push(d2)

	return nil
}

func instruction2ROT(ctx *InterpreterContext) error {
	d1, err := ctx.i.dstack.Peek(-6)
	if err != nil {
		return err
	}
	d2, err := ctx.i.dstack.Peek(-5)
	if err != nil {
		return err
	}

	ctx.i.dstack.Erase(-6, -4)
	ctx.i.dstack.Push(d1)
	ctx.i.dstack.Push(d2)

	return nil
}

func instruction2SWAP(ctx *InterpreterContext) error {
	if err := ctx.i.dstack.Swap(-4, -2); err != nil {
		return err
	}
	if err := ctx.i.dstack.Swap(-3, -1); err != nil {
		return err
	}

	return nil
}

func instructionIFDUP(ctx *InterpreterContext) error {
	d, err := ctx.i.dstack.Peek(0)
	if err != nil {
		return err
	}

	b := d.Boolean()
	if b {
		ctx.i.dstack.Push(d)
	}

	return nil
}

func instructionDEPTH(ctx *InterpreterContext) error {
	ctx.i.dstack.Push(Number(ctx.i.dstack.Depth()).Bytes())
	return nil
}

func instructionDROP(ctx *InterpreterContext) error {
	_, err := ctx.i.dstack.Pop()
	return err
}

func instructionDUP(ctx *InterpreterContext) error {
	d, err := ctx.i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	ctx.i.dstack.Push(d)

	return nil
}

func instructionNIP(ctx *InterpreterContext) error {
	return ctx.i.dstack.Remove(-2)
}

func instructionOVER(ctx *InterpreterContext) error {
	d, err := ctx.i.dstack.Peek(-2)
	if err != nil {
		return err
	}

	ctx.i.dstack.Push(d)
	return nil
}

func instructionROLL(ctx *InterpreterContext) error {
	d, err := ctx.i.dstack.Pop()
	if err != nil {
		return err
	}

	n, err := NewNumberFromBytes(d.Bytes(), ctx.flag.Has(ScriptVerifyMinimalData), NumberDefaultElementSize)
	if err != nil {
		return err
	}

	npos := int(n)
	if npos < 0 || npos >= ctx.i.dstack.Depth() {
		return ErrIntrepreterInvalidStackOperation
	}

	t, err := ctx.i.dstack.Peek(-npos - 1)
	if err != nil {
		return err
	}

	if ctx.ins.OPCode == OP_ROLL {
		ctx.i.dstack.Remove(-npos - 1)
	}

	ctx.i.dstack.Push(t)

	return nil
}

func instructionROT(ctx *InterpreterContext) error {
	if err := ctx.i.dstack.Swap(-3, -2); err != nil {
		return err
	}

	return ctx.i.dstack.Swap(-2, -1)
}

func instructionSWAP(ctx *InterpreterContext) error {
	return ctx.i.dstack.Swap(-2, -1)
}

func instructionTUCK(ctx *InterpreterContext) error {
	d, err := ctx.i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	return ctx.i.dstack.InsertBefore(0, d)
}
