package bscript

func instructionPushOPN(i *Interpreter, ins *Instruction, flag Flag) error {
	opcode := ins.OPCode
	i.dstack.Push(Number(int(opcode) - (int(OP_1) - 1)).Bytes())
	return nil
}

func instructionPushOP0(i *Interpreter, ins *Instruction, flag Flag) error {
	i.dstack.Push([]byte{})
	return nil
}

func instructionPushOPBytes(i *Interpreter, ins *Instruction, flag Flag) error {
	i.dstack.Push(ins.Data)
	return nil
}
