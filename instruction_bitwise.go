package bscript

func instructionEQUAL(i *Interpreter, ins *Instruction, flag Flag) error {
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

func instructionINVERT(i *Interpreter, ins *Instruction, flag Flag) error {
	d, err := i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	d.Invert()

	return nil
}

func instructionBITOP(i *Interpreter, ins *Instruction, flag Flag) error {
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

func instructionEQUALVERIFY(i *Interpreter, ins *Instruction, flag Flag) error {
	if err := instructionEQUAL(i, ins, flag); err != nil {
		return err
	}

	return instructionVERIFY(i, ins, flag)
}

func instructionLSFHIT(i *Interpreter, ins *Instruction, flag Flag) error {
	return nil
}

func instructionRSHIFT(i *Interpreter, ins *Instruction, flag Flag) error {
	return nil
}
