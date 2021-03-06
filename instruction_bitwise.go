package bscript

func instructionEQUAL(ctx *InterpreterContext) error {
	i := ctx.i
	d1, err := i.dstack.Peek(-2)
	if err != nil {
		return err
	}
	d2, err := i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	i.dstack.Pop()
	i.dstack.Pop()

	if d1.Equal(d2) {
		i.dstack.Push(Number(1).Bytes())
	} else {
		i.dstack.Push(Number(0).Bytes())
	}

	return nil
}

func instructionINVERT(ctx *InterpreterContext) error {
	i := ctx.i
	d, err := i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	d.Invert()

	return nil
}

func instructionBITOP(ctx *InterpreterContext) error {
	i := ctx.i
	ins := ctx.ins
	d1, err := i.dstack.Peek(-2)
	if err != nil {
		return err
	}
	d2, err := i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	if d1.Size() != d2.Size() {
		return ErrInterpreterOperandsSize
	}

	switch ins.OPCode {
	case OP_AND:
		d1.BitAnd(d2)
	case OP_OR:
		d1.BitOr(d2)
	case OP_XOR:
		d1.BitXor(d2)
	}

	i.dstack.Pop()

	return nil
}

func instructionEQUALVERIFY(ctx *InterpreterContext) error {
	if err := instructionEQUAL(ctx); err != nil {
		return err
	}

	return instructionVERIFY(ctx)
}

func instructionSFHIT(ctx *InterpreterContext) error {
	i := ctx.i
	ins := ctx.ins
	flag := ctx.flag
	d1, err := i.dstack.Peek(-2)
	if err != nil {
		return err
	}
	d2, err := i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	n1, err := d1.Number(flag.Has(ScriptVerifyMinimalData), NumberDefaultElementSize)
	if err != nil {
		return err
	}

	n2, err := d2.Number(flag.Has(ScriptVerifyMinimalData), NumberDefaultElementSize)
	if err != nil {
		return err
	}

	var n0 int
	if ins.OPCode == OP_LSHIFT {
		n0 = int(n1) << uint(n2)
	} else {
		n0 = int(n1) >> uint(n2)
	}

	i.dstack.Pop()
	i.dstack.Pop()
	i.dstack.Push(Number(n0).Bytes())

	return nil
}
