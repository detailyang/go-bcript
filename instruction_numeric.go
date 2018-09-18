package bscript

func instructionBINARY(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
	n1, err := i.dstack.Peek(-2)
	if err != nil {
		return err
	}
	n2, err := i.dstack.Peek(-1)
	if err != nil {
		return err
	}

	required := flag.Has(ScriptVerifyMinimalData)
	d1, err := NewNumberFromBytes(n1.Bytes(), required, NumberDefaultElementSize)
	if err != nil {
		return err
	}
	d2, err := NewNumberFromBytes(n2.Bytes(), required, NumberDefaultElementSize)
	if err != nil {
		return err
	}

	var d0 Number
	switch ins.OPCode {
	case OP_ADD:
		d0 = d1 + d2
	case OP_SUB:
		d0 = d1 - d2
	case OP_MUL:
		d0 = d1 * d2
	case OP_DIV:
		if d2 == 0 {
			return ErrInterpreterDivZero
		}
		d0 = d1 / d2
	case OP_MOD:
		if d2 == 0 {
			return ErrInterpreterModZero
		}

		d0 = d1 % d2
	case OP_BOOLAND:
		if d1 != 0 && d2 != 0 {
			d0 = 1
		} else {
			d0 = 0
		}
	case OP_BOOLOR:
		if d1 != 0 || d2 != 0 {
			d0 = 1
		} else {
			d0 = 0
		}
	case OP_NUMEQUAL:
		if d1 == d2 {
			d0 = 1
		} else {
			d0 = 0
		}
	case OP_NUMEQUALVERIFY:
		if d1 == d2 {
			d0 = 1
		} else {
			d0 = 0
		}
	case OP_NUMNOTEQUAL:
		if d1 != d2 {
			d0 = 1
		} else {
			d0 = 0
		}
	case OP_LESSTHAN:
		if d1 < d2 {
			d0 = 1
		} else {
			d0 = 0
		}
	case OP_GREATERTHAN:
		if d1 > d2 {
			d0 = 1
		} else {
			d0 = 0
		}
	case OP_LESSTHANOREQUAL:
		if d1 <= d2 {
			d0 = 1
		} else {
			d0 = 0
		}
	case OP_GREATERTHANOREQUAL:
		if d1 >= d2 {
			d0 = 1
		} else {
			d0 = 0
		}
	case OP_MIN:
		if d1 < d2 {
			d0 = d1
		} else {
			d0 = d2
		}
	case OP_MAX:
		if d1 < d2 {
			d0 = d2
		} else {
			d0 = d1
		}
	default:
		return ErrInterpreterBadInstruction
	}

	i.dstack.Pop()
	i.dstack.Pop()
	i.dstack.Push(d0.Bytes())

	if ins.OPCode == OP_NUMEQUALVERIFY {
		return instructionVERIFY(i, ins, flag, checker)
	}

	return nil
}

func instructionUNARY(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
	b, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	d, err := NewNumberFromBytes(b, flag.Has(ScriptVerifyMinimalData), NumberDefaultElementSize)
	if err != nil {
		return err
	}

	switch ins.OPCode {
	case OP_1ADD:
		d = d + 1
	case OP_1SUB:
		d = d - 1
	case OP_NEGATE:
		d = -d
	case OP_ABS:
		if d < 0 {
			d = -d
		}
	case OP_NOT:
		if d == 0 {
			d = 1
		} else {
			d = 0
		}
	case OP_0NOTEQUAL:
		if d != 0 {
			d = 1
		} else {
			d = 0
		}
	default:
		return ErrInterpreterBadInstruction
	}

	i.dstack.Push(d.Bytes())

	return nil
}

func instructionWITHIN(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
	required := flag.Has(ScriptVerifyMinimalData)

	d1, err := i.dstack.Peek(-3)
	if err != nil {
		return err
	}
	n1, err := NewNumberFromBytes(d1.Bytes(), required, NumberDefaultElementSize)
	if err != nil {
		return err
	}

	d2, err := i.dstack.Peek(-2)
	if err != nil {
		return err
	}
	n2, err := NewNumberFromBytes(d2.Bytes(), required, NumberDefaultElementSize)
	if err != nil {
		return err
	}

	d3, err := i.dstack.Peek(-1)
	if err != nil {
		return err
	}
	n3, err := NewNumberFromBytes(d3.Bytes(), required, NumberDefaultElementSize)
	if err != nil {
		return err
	}

	i.dstack.Pop()
	i.dstack.Pop()
	i.dstack.Pop()

	if n2 <= n1 && n1 < n3 {
		i.dstack.Push(Number(1).Bytes())
	} else {
		i.dstack.Push(Number(0).Bytes())
	}

	return nil
}
