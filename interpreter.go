package bscript

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	. "github.com/detailyang/go-bprimitives"
)

var (
	ErrInterpreterScriptSize                         = errors.New("interpreter: over script size 1000")
	ErrInterpreterScriptOPCount                      = errors.New("interpreter: over script op count 210")
	ErrInterpreterInvalidStackOperation              = errors.New("interpreter: invalid stack operation")
	ErrInterpreterOperandsSize                       = errors.New("interpreter: operands size are not equal")
	ErrInterpreterVerifyFailed                       = errors.New("interpreter: verify failed")
	ErrInterpreterDivZero                            = errors.New("interpreter: div zero")
	ErrInterpreterModZero                            = errors.New("interpreter: mod zero")
	ErrInterpreterBadInstruction                     = errors.New("interpreter: bad instruction")
	ErrInterpreterStackSizeNotEnough                 = errors.New("interpreter: stack not enough")
	ErrInterpreterNoMatchConditional                 = errors.New("interpreter: no match conditional")
	ErrInterpreterBadOPCode                          = errors.New("interpreter: bad opcode")
	ErrInterpreterIllegalOPCode                      = errors.New("interpreter: illegal opcode")
	ErrInterpreterStackOverflow                      = errors.New("interpreter: data stack overflow")
	ErrInterpreterUnbalancedConditional              = errors.New("interpreter: unbalanced conditional")
	ErrInterpreterDisabledOPCode                     = errors.New("interpreter: disabled opcode")
	ErrInterpreterNegativeLocktime                   = errors.New("interpreter: negative locktime")
	ErrInterpreterUnsatisfiedLocktime                = errors.New("interpreter: unsatisfied locktime")
	ErrInterpreterDiscourageUpgradableNops           = errors.New("interpreter: discourage upgradable nops")
	ErrInterpreterSignaturePushOnly                  = errors.New("interpreter: signature push only")
	ErrInterpreterP2SHBadStack                       = errors.New("interpreter: p2sh bad stack")
	ErrInterpreterParseWitnessFailed                 = errors.New("interpreter: parse witness program failed")
	ErrInterpreterPushSize                           = errors.New("interpreter: bad push size")
	ErrInterpreterWitnessMalleatedP2SH               = errors.New("interpreter: malleated P2SH")
	ErrInterpreterDiscourageUpgradableWitnessProgram = errors.New("interpreter: discourage upgradable witness program")
	ErrInterpreterWitnessProgramWitnessEmpty         = errors.New("interpreter: witness program witness empty")
	ErrInterpreterWitnessProgramMismatch             = errors.New("interpreter: witness program mismatch")
	ErrInterpreterWitnessProgramWrongLength          = errors.New("interpreter: witness program wrong length")
	ErrInterpreterWitnessVerifyFailed                = errors.New("interpreter: witness verify failed")
	ErrInterpreterWitnessMalleated                   = errors.New("interpreter: witness malleated")
	ErrInterpreterCleanStack                         = errors.New("interpreter: clean stack")
	ErrInterpreterWitnessUnexpected                  = errors.New("interpreter: wintess unexpected")
	ErrInterpreterScriptPubekyesPerMultisig          = errors.New("interpreter: too many multisig")
	ErrInterpreterSignatureNullDummy                 = errors.New("interpreter: siganture null dummy")
	ErrInterpreterBadSignatureDer                    = errors.New("interpreter: bad signature der format")
	ErrInterpreterSigantureHighS                     = errors.New("interpreter: signature invalid high s")
	ErrInterpreterBadSignatureHashType               = errors.New("interpreter: bad signature hash type")
	ErrInterpreterBadPubkey                          = errors.New("interpreter: bad public key")
	ErrInterpreterEvalFalse                          = errors.New("interpreter: eval false")
)

const (
	MaxInterpreterScriptSize                = 1000
	MaxInterpreterScriptOPS                 = 210
	MaxInterpreterScriptPubekyesPerMultisig = 20
)

type Interpreter struct {
	dstack  *Stack
	astack  *Stack
	cstack  []int
	pc      int
	codesep int
	traces  []Trace
}

type InterpreterContext struct {
	sigver  SignatureVersion
	flag    Flag
	script  *Script
	i       *Interpreter
	ins     *Instruction
	checker Checker
}

func NewInterpreterContext(script *Script, i *Interpreter, ins *Instruction, checker Checker, flag Flag, sigver SignatureVersion) *InterpreterContext {
	return &InterpreterContext{
		sigver:  sigver,
		flag:    flag,
		script:  script,
		i:       i,
		ins:     ins,
		checker: checker,
	}
}

func copySlice(s []byte) []byte {
	t := make([]byte, len(s))
	copy(t, s)
	return t
}

func verifyWitnessProgramm(
	scriptWitness ScriptWitness,
	witnessVersion uint8,
	wintessProgram []byte,
	flag Flag,
	checker Checker,
) error {
	if witnessVersion != 0 {
		if flag.Has(ScriptVerifyDiscourageUpgradeableWitnessProgram) {
			return ErrInterpreterDiscourageUpgradableWitnessProgram
		}

		return nil
	}

	witnessStack := NewStack()
	scriptPubkey := NewScript()

	if len(wintessProgram) == 32 {
		if scriptWitness.Size() == 0 {
			return ErrInterpreterWitnessProgramWitnessEmpty
		}

		pubkey := scriptWitness[scriptWitness.Size()-1]
		stack := scriptWitness[:scriptWitness.Size()-1]
		scriptPubkeyHash := Hash256(scriptPubkey.Bytes())

		if !scriptPubkeyHash.Equal(NewHash(wintessProgram[0:32])) {
			return ErrInterpreterWitnessProgramMismatch
		}

		for _, s := range stack {
			witnessStack.Push(copySlice(s))
		}

		scriptPubkey.PushBytes(pubkey)

	} else if len(wintessProgram) == 20 {
		if scriptWitness.Size() != 2 {
			return ErrInterpreterWitnessProgramMismatch
		}

		scriptPubkey.
			PushOPCode(OP_DUP).
			PushOPCode(OP_HASH160).
			PushBytesWithOP(wintessProgram).
			PushOPCode(OP_EQUALVERIFY).
			PushOPCode(OP_CHECKSIG)

		for _, s := range scriptWitness {
			witnessStack.Push(copySlice(s))
		}
	} else {
		return ErrInterpreterWitnessProgramWrongLength
	}

	ok := true
	witnessStack.Iter(func(e StackElemnt) {
		if len(e.Bytes()) > MaxInterpreterScriptSize {
			ok = false
		}
	})

	if !ok {
		return ErrInterpreterPushSize
	}

	interpreter := NewInterpreter()
	interpreter.SetDStack(witnessStack)
	err := interpreter.Eval(scriptPubkey, flag, checker, SignatureVersionWitnessV0)
	if err != nil {
		return err
	}

	d, err := witnessStack.Peek(-1)
	if err != nil {
		return err
	}

	if !d.Boolean() {
		return ErrInterpreterWitnessVerifyFailed
	}

	return nil
}

func EvalScript(script *Script, flag Flag, checker Checker) error {
	interpreter := NewInterpreter()
	return interpreter.Eval(script, flag, checker, SignatureVersionBase)
}

func VerifyScript(scriptSig, scriptPubkey *Script, scriptWitness ScriptWitness, flag Flag, checker Checker, sigversion SignatureVersion) error {
	if flag.Has(ScriptVerifySigPushOnly) && !scriptSig.IsPushOnly() {
		return ErrInterpreterSignaturePushOnly
	}

	interpreter := NewInterpreter()
	stack := NewStack()
	interpreter.SetDStack(stack)
	stackCopy := NewStack()
	hadWitness := false
	cleanStack := flag.Has(ScriptVerifyCleanStack)

	fmt.Println("unlocking", scriptSig.String())
	fmt.Println("locking", scriptPubkey.String())

	err := interpreter.Eval(scriptSig, flag, checker, sigversion)
	if err != nil {
		return err
	}

	if flag.Has(ScriptVerifyP2SH) {
		stackCopy = stack.Clone()
	}

	interpreter.astack.Clean()

	err = interpreter.Eval(scriptPubkey, flag, checker, sigversion)
	if err != nil {
		return err
	}

	if interpreter.dstack.Depth() == 0 {
		return ErrInterpreterEvalFalse
	}

	d, err := interpreter.dstack.Peek(-1)
	if err != nil {
		return err
	}

	if d.Boolean() == false {
		return ErrInterpreterEvalFalse
	}

	// Verify witness program
	if flag.Has(ScriptVerifyWitness) {
		witnessVersion, witnessProgram, ok := scriptPubkey.ParseWitnessProgram()
		if !ok {
			return ErrInterpreterWitnessMalleated
		}

		hadWitness = true
		cleanStack = false

		err = verifyWitnessProgramm(
			scriptWitness,
			witnessVersion,
			witnessProgram,
			flag,
			checker)
		if err != nil {
			return err
		}
	}

	if flag.Has(ScriptVerifyP2SH) && scriptPubkey.IsPayToScriptHash() {
		if !scriptSig.IsPushOnly() {
			return ErrInterpreterSignaturePushOnly
		}

		stack.CloneFrom(stackCopy)

		// stack cannot be empty here, because if it was the
		// P2SH  HASH <> EQUAL  scriptPubKey would be evaluated with
		// an empty stack and the EvalScript above would return false.
		if stack.Depth() == 0 {
			return ErrInterpreterP2SHBadStack
		}

		d, err := stack.Pop()
		if err != nil {
			return err
		}

		pubkey := NewScriptFromBytes(d.Bytes())
		err = interpreter.Eval(pubkey, flag, checker, sigversion)
		if err != nil {
			return err
		}

		if flag.Has(ScriptVerifyWitness) {
			witnessVersion, witnessProgram, ok := pubkey.ParseWitnessProgram()
			if !ok {
				return ErrInterpreterParseWitnessFailed
			}

			if !bytes.Equal(scriptSig.Bytes(), NewScript().PushBytesWithOP(pubkey.Bytes()).Bytes()) {
				return ErrInterpreterWitnessMalleatedP2SH
			}

			hadWitness = true
			cleanStack = false
			err = verifyWitnessProgramm(
				scriptWitness,
				witnessVersion,
				witnessProgram,
				flag,
				checker)
			if err != nil {
				return err
			}

		}
	}

	// The CLEANSTACK check is only performed after potential P2SH evaluation,
	// as the non-P2SH evaluation of a P2SH script will obviously not result in
	// a clean stack (the P2SH inputs remain). The same holds for witness evaluation.
	if cleanStack {
		if stack.Depth() != 1 {
			return ErrInterpreterCleanStack
		}
	}

	if flag.Has(ScriptVerifyWitness) {
		if !hadWitness && scriptWitness.Size() == 0 {
			return ErrInterpreterWitnessUnexpected
		}
	}

	return nil
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

func (i *Interpreter) Eval(script *Script, flag Flag, checker Checker, sigversion SignatureVersion) error {
	if script.Size() > MaxInterpreterScriptSize {
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
			if nop > MaxInterpreterScriptOPS {
				return ErrInterpreterScriptOPCount
			}
		}

		if !flag.Has(ScriptSkipDisabledOPCode) && opcode.IsDisabled() {
			return ErrInterpreterDisabledOPCode
		}

		if ins.IsIllegal() {
			return ErrInterpreterIllegalOPCode
		}

		if i.shouldSkip() && !ins.IsConditional() {
			continue
		}

		operator, ok := instructionOperator[opcode]
		if !ok {
			return ErrInterpreterBadOPCode
		}

		ctx := NewInterpreterContext(script, i, ins, checker, flag, sigversion)
		if err := operator(ctx); err != nil {
			return err
		}

		if i.dstack.Depth()+i.astack.Depth() > 1000 {
			return ErrInterpreterStackOverflow
		}

		if flag.Has(ScriptEnableTrace) {
			trace := Trace{
				Step:      i.pc,
				Executed:  ins.OPCode.String(),
				Stack:     i.dstack.String(),
				Remaining: script.Disassemble(" "),
			}
			i.traces = append(i.traces, trace)
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

func (i *Interpreter) PrintTraces() {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
	for i, trace := range i.traces {
		if i == 0 {
			fmt.Fprintln(w, "\n#Step\tExecuted OP Code\tResulted Stack\tRemaining OP Codes\t")
		}
		row := fmt.Sprintf("%04d\t%s\t%s\t%s\t", trace.Step, trace.Executed, trace.Stack, trace.Remaining)
		fmt.Fprintln(w, row)
	}
	w.Flush()
}
