package bscript

func instructionTOTALSTACK(i *Interpreter, ins *Instruction, flag Flag) error {
	data, err := i.dstack.Pop()
	if err != nil {
		return err
	}
	i.astack.Push(data)

	return nil
}

func instructionFROMALTSTACK(i *Interpreter, ins *Instruction, flag Flag) error {
	data, err := i.astack.Pop()
	if err != nil {
		return err
	}
	i.dstack.Push(data)

	return nil
}

func instruction2DROP(i *Interpreter, ins *Instruction, flag Flag) error {
	if i.dstack.Depth() < 2 {
		return ErrInterpreterStackSizeNotEnough
	}

	i.dstack.Pop()
	i.dstack.Pop()

	return nil
}

func instruction2DUP(i *Interpreter, ins *Instruction, flag Flag) error {
	d1, err := i.dstack.Peek(-2)
	if err != nil {
		return err
	}
	d2, err := i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	i.dstack.Push(d1)
	i.dstack.Push(d2)

	return nil
}

func instruction3DUP(i *Interpreter, ins *Instruction, flag Flag) error {
	d1, err := i.dstack.Peek(-3)
	if err != nil {
		return err
	}
	d2, err := i.dstack.Peek(-2)
	if err != nil {
		return err
	}
	d3, err := i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	i.dstack.Push(d1)
	i.dstack.Push(d2)
	i.dstack.Push(d3)

	return nil
}

func instruction2OVER(i *Interpreter, ins *Instruction, flag Flag) error {
	d1, err := i.dstack.Peek(-4)
	if err != nil {
		return err
	}
	d2, err := i.dstack.Peek(-3)
	if err != nil {
		return err
	}

	i.dstack.Push(d1)
	i.dstack.Push(d2)

	return nil
}

func instruction2ROT(i *Interpreter, ins *Instruction, flag Flag) error {
	d1, err := i.dstack.Peek(-6)
	if err != nil {
		return err
	}
	d2, err := i.dstack.Peek(-5)
	if err != nil {
		return err
	}

	i.dstack.Erase(-6, -4)
	i.dstack.Push(d1)
	i.dstack.Push(d2)

	return nil
}

func instruction2SWAP(i *Interpreter, ins *Instruction, flag Flag) error {
	if err := i.dstack.Swap(-4, -2); err != nil {
		return err
	}
	if err := i.dstack.Swap(-3, -1); err != nil {
		return err
	}

	return nil
}

func instructionIFDUP(i *Interpreter, ins *Instruction, flag Flag) error {
	d, err := i.dstack.Peek(0)
	if err != nil {
		return err
	}

	b := d.Boolean()
	if b {
		i.dstack.Push(d)
	}

	return nil
}

func instructionDEPTH(i *Interpreter, ins *Instruction, flag Flag) error {
	i.dstack.Push(Number(i.dstack.Depth()).Bytes())
	return nil
}

func instructionDROP(i *Interpreter, ins *Instruction, flag Flag) error {
	_, err := i.dstack.Pop()
	return err
}

func instructionDUP(i *Interpreter, ins *Instruction, flag Flag) error {
	d, err := i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	i.dstack.Push(d)

	return nil
}

func instructionNIP(i *Interpreter, ins *Instruction, flag Flag) error {
	return i.dstack.Remove(-2)
}

func instructionOVER(i *Interpreter, ins *Instruction, flag Flag) error {
	d, err := i.dstack.Peek(-2)
	if err != nil {
		return err
	}

	i.dstack.Push(d)
	return nil
}

func instructionROLL(i *Interpreter, ins *Instruction, flag Flag) error {
	d, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	n, err := NewNumberFromBytes(d.Bytes(), flag.Has(ScriptVerifyMinimalData), NumberDefaultElementSize)
	if err != nil {
		return err
	}

	npos := int(n)
	if npos < 0 || npos >= i.dstack.Depth() {
		return ErrIntrepreterInvalidStackOperation
	}

	t, err := i.dstack.Peek(-npos - 1)
	if err != nil {
		return err
	}

	if ins.OPCode == OP_ROLL {
		i.dstack.Remove(-npos - 1)
	}

	i.dstack.Push(t)

	return nil
}

func instructionROT(i *Interpreter, ins *Instruction, flag Flag) error {
	if err := i.dstack.Swap(-3, -2); err != nil {
		return err
	}

	return i.dstack.Swap(-2, -1)
}

func instructionSWAP(i *Interpreter, ins *Instruction, flag Flag) error {
	return i.dstack.Swap(-2, -1)
}

func instructionTUCK(i *Interpreter, ins *Instruction, flag Flag) error {
	d, err := i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	return i.dstack.InsertBefore(0, d)
}
