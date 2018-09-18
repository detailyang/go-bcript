package bscript

func instructionPushOPN(ctx *InterpreterContext) error {
	i := ctx.i
	ins := ctx.ins
	opcode := ins.OPCode

	i.dstack.Push(Number(int(opcode) - (int(OP_1) - 1)).Bytes())
	return nil
}

func instructionPushOP0(ctx *InterpreterContext) error {
	ctx.i.dstack.Push([]byte{})
	return nil
}

func instructionPushOPBytes(ctx *InterpreterContext) error {
	ctx.i.dstack.Push(ctx.ins.Data)
	return nil
}
