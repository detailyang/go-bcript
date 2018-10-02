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
