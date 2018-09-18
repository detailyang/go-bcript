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

func instructionRIPEMD160(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
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

func instructionSHA1(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
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

func instructionSHA256(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
	d, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	hash := sha256.Sum256(d.Bytes())
	i.dstack.Push(hash[:])

	return nil
}

func instructionHASH160(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
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

func instructionHASH256(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
	d, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	hash := sha256.Sum256(d.Bytes())
	hash = sha256.Sum256(hash[:])
	i.dstack.Push(hash[:])

	return nil
}

func instructionCODESEPARATOR(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
	i.codesep = i.pc
	return nil
}

func instructionCHECKSIG(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
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

	return nil
}

func instructionCHECKSIGVERIFY(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
	return nil
}

func instructionCHECKMULTISIG(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
	return nil
}

func instructionCHECKMULTISIGVERIFY(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
	return nil
}

func instructionCHECKLOCKTIMEVERIFY(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
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
			return ErrIntrepreterNegativeLocktime
		}

		if err := checker.CheckLockTime(uint32(locktime)); err != nil {
			return ErrIntrepreterUnsatisfiedLocktime
		}

	} else if flag.Has(ScriptDiscourageUpgradableNops) {
		return ErrIntrepreterDiscourageUpgradableNops
	}

	return nil
}

func instructionCHECKSEQUENCEVERIFY(i *Interpreter, ins *Instruction, flag Flag, checker Checker) error {
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
			return ErrIntrepreterNegativeLocktime
		}

		if sequence&SequenceLockTimeDisabledFlag == SequenceLockTimeDisabledFlag {
			if err := checker.CheckSequence(uint32(sequence)); err != nil {
				return ErrIntrepreterUnsatisfiedLocktime
			}
		}

	} else if flag.Has(ScriptDiscourageUpgradableNops) {
		return ErrIntrepreterDiscourageUpgradableNops
	}

	return nil
}
