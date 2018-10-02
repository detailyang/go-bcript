package bscript

import (
	"encoding/hex"
	"fmt"
)

func instructionCHECKSIG(ctx *InterpreterContext) error {
	script := ctx.script
	i := ctx.i
	flag := ctx.flag
	sigver := ctx.sigver
	checker := ctx.checker

	d1, err := i.dstack.Pop()
	if err != nil {
		return err
	}
	d2, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	pubkey := d1.Bytes()
	sig := d2.Bytes()
	if len(sig) < 1 {
		i.dstack.Push(Boolean(false).Bytes())
		return nil
	}

	hashtype := sig[len(sig)-1]

	if err := CheckHashTypeEncoding(hashtype, flag); err != nil {
		return err
	}

	if err := CheckPubkeyEncoding(pubkey, flag); err != nil {
		return err
	}

	if err := CheckSignatureEncoding(sig, flag, sigver); err != nil {
		return err
	}

	subscript, err := script.SubScript(i.codesep)
	if err != nil {
		return err
	}

	if sigver == SignatureVersionBase {
		sigscript := NewScript().PushBytesWithOP(sig)
		subscript = subscript.Filter(sigscript)
	}

	if err := checker.CheckSignature(sig, pubkey, subscript, sigver); err != nil {
		i.dstack.Push(Number(0).Bytes())
	} else {
		i.dstack.Push(Number(1).Bytes())
	}

	if ctx.ins.OPCode == OP_CHECKSIGVERIFY {
		return instructionVERIFY(ctx)
	}

	return nil
}

// Compares the first signature against each public key until it finds an ECDSA match.
// Starting with the subsequent public key, it compares the second signature against each remaining public key until it finds an ECDSA match.
// The process is repeated until all signatures have been checked or not enough public keys remain to produce a successful result.
// All signatures need to match a public key. Because public keys are not checked again if they fail any signature comparison,
// signatures must be placed in the scriptSig using the same order as their corresponding public keys were placed in the scriptPubKey or redeemScript.
// If all signatures are valid, 1 is returned, 0 otherwise. Due to a bug, one extra unused value is removed from the stack.
func instructionCHECKMULTISIG(ctx *InterpreterContext) error {
	d, err := ctx.i.dstack.Pop()
	if err != nil {
		return err
	}

	n, err := d.Number(ctx.flag.Has(ScriptVerifyMinimalData), 4)
	if err != nil {
		return err
	}

	if n <= 0 || n > MaxInterpreterScriptPubekyesPerMultisig {
		return ErrInterpreterBadInstruction
	}

	keys := make([][]byte, n)
	for i := 0; i < int(n); i++ {
		d, err := ctx.i.dstack.Pop()
		if err != nil {
			return err
		}

		keys[i] = d.Bytes()
	}

	d, err = ctx.i.dstack.Pop()
	if err != nil {
		return err
	}

	n, err = d.Number(ctx.flag.Has(ScriptVerifyMinimalData), 4)
	if err != nil {
		return err
	}

	if n == 0 {
		return ErrInterpreterBadInstruction
	}

	sigs := make([][]byte, n)
	for i := 0; i < int(n); i++ {
		d, err := ctx.i.dstack.Pop()
		if err != nil {
			return err
		}

		sigs[i] = d.Bytes()
	}

	subscript, err := ctx.script.SubScript(ctx.i.codesep)
	if err != nil {
		return err
	}

	for _, sig := range sigs {
		if ctx.sigver == SignatureVersionBase|SignatureVersionForkId {
			sigscript := NewScript().PushBytesWithOP(sig)
			subscript = subscript.Filter(sigscript)
		}
	}

	success := true
	k := 0
	s := 0
	for s < len(sigs) && success {
		key := keys[k]
		sig := sigs[s]

		if err := CheckSignatureEncoding(sig, ctx.flag, ctx.sigver); err != nil {
			return err
		}

		if err := CheckPubkeyEncoding(key, ctx.flag); err != nil {
			return err
		}

		fmt.Println("abccsig", hex.EncodeToString(sig))
		fmt.Println("abcckey", hex.EncodeToString(key))
		fmt.Println("abccscript", hex.EncodeToString(subscript.Bytes()))
		err := ctx.checker.CheckSignature(sig, key, subscript, ctx.sigver)
		fmt.Println(err)
		if err == nil {
			s += 1
		}
		k += 1

		success = len(sigs)-s <= len(keys)-k
	}

	if ctx.i.dstack.Depth() > 0 && ctx.flag.Has(ScriptVerifyNullFail) {
		return ErrInterpreterSignatureNullDummy
	}

	if success {
		ctx.i.dstack.Push(Number(1).Bytes())
	} else {
		ctx.i.dstack.Push(Number(0).Bytes())
	}

	if ctx.ins.OPCode == OP_CHECKMULTISIGVERIFY {
		return instructionVERIFY(ctx)
	}

	return nil
}

func instructionCHECKLOCKTIMEVERIFY(ctx *InterpreterContext) error {
	i := ctx.i
	flag := ctx.flag
	checker := ctx.checker

	if flag.Has(ScriptVerifyCheckLockTimeVerify) {
		d, err := i.dstack.Pop()
		if err != nil {
			return err
		}

		locktime, err := d.Number(flag.Has(ScriptVerifyMinimalData), 5)
		if err != nil {
			return err
		}

		if locktime.IsNegative() {
			return ErrInterpreterNegativeLocktime
		}

		if err := checker.CheckLockTime(uint32(locktime)); err != nil {
			return ErrInterpreterUnsatisfiedLocktime
		}

	} else if flag.Has(ScriptDiscourageUpgradableNops) {
		return ErrInterpreterDiscourageUpgradableNops
	}

	return nil
}

func instructionCHECKSEQUENCEVERIFY(ctx *InterpreterContext) error {
	i := ctx.i
	flag := ctx.flag
	checker := ctx.checker

	if flag.Has(ScriptVerifyCheckSequenceVerify) {
		d, err := i.dstack.Pop()
		if err != nil {
			return err
		}

		sequence, err := d.Number(flag.Has(ScriptVerifyMinimalData), 5)
		if err != nil {
			return err
		}

		if sequence.IsNegative() {
			return ErrInterpreterNegativeLocktime
		}

		if sequence&SequenceLockTimeDisabledFlag == SequenceLockTimeDisabledFlag {
			if err := checker.CheckSequence(uint32(sequence)); err != nil {
				return ErrInterpreterUnsatisfiedLocktime
			}
		}

	} else if flag.Has(ScriptDiscourageUpgradableNops) {
		return ErrInterpreterDiscourageUpgradableNops
	}

	return nil
}
