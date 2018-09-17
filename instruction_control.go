package bscript

func instructionVERIFY(i *Interpreter, ins *Instruction, flag Flag) error {
	d, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	ok := NewBoolean(d)
	if !ok {
		return ErrInterpreterVerifyFailed
	}

	return nil
}
