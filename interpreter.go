package bscript

import (
	"errors"
)

var (
	ErrInterpreterScriptSize               = errors.New("intrepreter: over script size 1000")
	ErrIntrepreterScriptOPCount            = errors.New("intrepreter: over script op count 210")
	ErrIntrepreterInvalidStackOperation    = errors.New("intrepreter: invalid stack operation")
	ErrInterpreterOperandsSize             = errors.New("interpreter: operands size are not equal")
	ErrInterpreterVerifyFailed             = errors.New("interpreter: verify failed")
	ErrInterpreterDivZero                  = errors.New("interpreter: div zero")
	ErrInterpreterModZero                  = errors.New("interpreter: mod zero")
	ErrInterpreterBadInstruction           = errors.New("intrepreter: bad instruction")
	ErrInterpreterStackSizeNotEnough       = errors.New("intrepreter: stack not enough")
	ErrInterpreterNoMatchConditional       = errors.New("interpreter: no match conditional")
	ErrIntrepreterBadOPCode                = errors.New("intrepreter: bad opcode")
	ErrIntrepreterStackOverflow            = errors.New("intrepreter: data stack overflow")
	ErrInterpreterUnbalancedConditional    = errors.New("interpreter: unbalanced conditional")
	ErrIntrepreterDisabledOPCode           = errors.New("interpreter: disabled opcode")
	ErrIntrepreterNegativeLocktime         = errors.New("interpreter: negative locktime")
	ErrIntrepreterUnsatisfiedLocktime      = errors.New("intrepreter: unsatisfied locktime")
	ErrIntrepreterDiscourageUpgradableNops = errors.New("intrepreter: discourage upgradable nops")
)

const (
	MaxIntrepreterScriptSize = 1000
	MaxINtrepreterScriptOPS  = 210
)

type Interpreter struct {
	dstack  *Stack
	astack  *Stack
	cstack  []int
	pc      int
	codesep int
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		dstack:  NewStack(),
		astack:  NewStack(),
		cstack:  make([]int, 0, 128),
		pc:      0,
		codesep: 0,
	}
}

func (i *Interpreter) SetDStack(dstack *Stack) {
	i.dstack = dstack
}

func (i *Interpreter) GetDStack() *Stack {
	return i.dstack
}

func (i *Interpreter) Eval(script *Script, flag Flag, checker Checker) error {
	if script.Size() > MaxIntrepreterScriptSize {
		return ErrInterpreterScriptSize
	}

	nop := 0

	for {
		ins, err := script.Next()
		if err != nil {
			if err == ErrScriptEOF {
				break
			}
			return err
		}

		i.pc++

		opcode := ins.OPCode
		if opcode.IsCountable() {
			nop++
			if nop > MaxINtrepreterScriptOPS {
				return ErrIntrepreterScriptOPCount
			}
		}

		if !flag.Has(ScriptSkipDisabledOPCode) && opcode.IsDisabled() {
			return ErrIntrepreterDisabledOPCode
		}

		operator, ok := instructionOperator[opcode]
		if !ok {
			return ErrIntrepreterBadOPCode
		}

		if err := operator(i, ins, flag, checker); err != nil {
			return err
		}

		if i.dstack.Depth()+i.astack.Depth() > 1000 {
			return ErrIntrepreterStackOverflow
		}
	}

	if len(i.cstack) > 0 {
		return ErrInterpreterUnbalancedConditional
	}

	return nil
}

func (i *Interpreter) shouldSkip() bool {
	if len(i.cstack) == 0 {
		return false
	}

	return i.cstack[len(i.cstack)-1] != OpCondTrue
}
