package bscript

func instructionCAT(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
	d1, err := i.dstack.Peek(-2)
	if err != nil {
		return err
	}

	d2, err := i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	d1.Cat(d2)

	i.dstack.Pop()

	return nil
}

func instructionSUBSTR(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
	return nil
}

func instructionLEFT(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
	return nil
}

func instructionRIGHT(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
	return nil
}

func instructionSIZE(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
	d, err := i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	i.dstack.Push(Number(d.Size()).Bytes())

	return nil
}
