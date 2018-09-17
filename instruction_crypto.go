package bscript

import (
	"crypto/sha1"
	"crypto/sha256"

	"github.com/tokublock/tokucore/xcrypto/ripemd160"
)

func instructionRIPEMD160(i *Interpreter, ins *Instruction, flag Flag) error {
	d, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	hasher := ripemd160.New()
	hasher.Write(d.Bytes())
	hash := hasher.Sum(nil)

	i.dstack.Push(valtype(hash))

	return nil
}

func instructionSHA1(i *Interpreter, ins *Instruction, flag Flag) error {
	d, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	hasher := sha1.New()
	hasher.Write(d.Bytes())
	hash := hasher.Sum(nil)

	i.dstack.Push(valtype(hash))

	return nil
}

func instructionSHA256(i *Interpreter, ins *Instruction, flag Flag) error {
	d, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	hash := sha256.Sum256(d.Bytes())
	i.dstack.Push(valtype(hash[:]))

	return nil
}

func instructionHASH160(i *Interpreter, ins *Instruction, flag Flag) error {
	d, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	hash := sha256.Sum256(d.Bytes())
	hasher := ripemd160.New()
	hasher.Write(hash[:])
	i.dstack.Push(valtype(hasher.Sum(nil)))

	return nil
}

func instructionHASH256(i *Interpreter, ins *Instruction, flag Flag) error {
	d, err := i.dstack.Pop()
	if err != nil {
		return err
	}

	hash := sha256.Sum256(d.Bytes())
	hash = sha256.Sum256(hash[:])
	i.dstack.Push(valtype(hash[:]))

	return nil
}

func instructionCODESEPARATOR(i *Interpreter, ins *Instruction, flag Flag) error {
	i.codesep = i.pc
	return nil
}

func instructionCHECKSIG(i *Interpreter, ins *Instruction, flag Flag) error {
	return nil
}

func instructionCHECKSIGVERIFY(i *Interpreter, ins *Instruction, flag Flag) error {
	return nil
}

func instructionCHECKMULTISIG(i *Interpreter, ins *Instruction, flag Flag) error {
	return nil
}

func instructionCHECKMULTISIGVERIFY(i *Interpreter, ins *Instruction, flag Flag) error {
	return nil
}
