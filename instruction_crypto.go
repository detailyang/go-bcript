package bscript

import (
	"crypto/sha1"
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

// Below flags apply in the context of BIP 68
const (
	SequenceLockTimeDisabledFlag = 1 << 31
)

func instructionRIPEMD160(ctx *InterpreterContext) error {
	i := ctx.i
	d, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	hasher := ripemd160.New()
	hasher.Write(d.Bytes())
	hash := hasher.Sum(nil)

	i.dstack.Push(hash)

	return nil
}

func instructionSHA1(ctx *InterpreterContext) error {
	i := ctx.i
	d, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	hasher := sha1.New()
	hasher.Write(d.Bytes())
	hash := hasher.Sum(nil)

	i.dstack.Push(hash)

	return nil
}

func instructionSHA256(ctx *InterpreterContext) error {
	i := ctx.i
	d, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	hash := sha256.Sum256(d.Bytes())
	i.dstack.Push(hash[:])

	return nil
}

func instructionHASH160(ctx *InterpreterContext) error {
	i := ctx.i
	d, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	hash := sha256.Sum256(d.Bytes())
	hasher := ripemd160.New()
	hasher.Write(hash[:])
	i.dstack.Push(hasher.Sum(nil))

	return nil
}

func instructionHASH256(ctx *InterpreterContext) error {
	i := ctx.i

	d, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	hash := sha256.Sum256(d.Bytes())
	hash = sha256.Sum256(hash[:])
	i.dstack.Push(hash[:])

	return nil
}

func instructionCODESEPARATOR(ctx *InterpreterContext) error {
	i := ctx.i
	i.codesep = i.pc
	return nil
}

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
	hashtype := sig[len(sig)-1]

	if len(sig) < 1 {
		i.dstack.Push(Boolean(false).Bytes())
		return nil
	}

	if err := CheckHashTypeEncoding(hashtype, flag); err != nil {
		return err
	}

	if err := CheckPubkeyEncoding(pubkey, flag); err != nil {
		return err
	}

	if err := CheckSignatureEncoding(sig, flag); err != nil {
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
		return err
	}

	if ctx.ins.OPCode == OP_CHECKSIGVERIFY {
		return instructionVERIFY(ctx)
	}

	return nil
}

func instructionCHECKMULTISIG(ctx *InterpreterContext) error {
	return nil
}

func instructionCHECKMULTISIGVERIFY(ctx *InterpreterContext) error {
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
