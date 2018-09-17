package bscript

// Conditional execution constants.
const (
	OpCondFalse = 0
	OpCondTrue  = 1
	OpCondSkip  = 2
)

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

func instructionIF(i *Interpreter, ins *Instruction, flag Flag) error {
	cond := OpCondFalse
	if !i.shouldSkip() {
		d, err := i.dstack.Pop()
		if err != nil {
			return err
		}

		b := d.Boolean()

		if ins.OPCode == OP_IF {
			if b {
				cond = OpCondTrue
			}
		} else {
			if !b {
				cond = OpCondTrue
			}
		}
	} else {
		cond = OpCondSkip
	}

	i.cstack = append(i.cstack, cond)

	return nil
}

func instructionENDIF(i *Interpreter, ins *Instruction, flag Flag) error {
	if len(i.cstack) == 0 {
		return ErrInterpreterNoMatchConditional
	}

	i.cstack = i.cstack[:len(i.cstack)-1]
	return nil
}

func instructionELSE(i *Interpreter, ins *Instruction, flag Flag) error {
	if len(i.cstack) == 0 {
		return ErrInterpreterNoMatchConditional
	}

	if i.cstack[len(i.cstack)-1] == OpCondTrue {
		i.cstack[len(i.cstack)-1] = OpCondFalse
	} else if i.cstack[len(i.cstack)-1] == OpCondFalse {
		i.cstack[len(i.cstack)-1] = OpCondTrue
	}

	return nil
}
