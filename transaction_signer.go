package bscript

// One thing all signature hash types sign is the transaction’s locktime. (Called nLockTime in the Bitcoin Core source code.) The locktime indicates the earliest time a transaction can be added to the block chain.

// Locktime allows signers to create time-locked transactions which will only become valid in the future, giving the signers a chance to change their minds.

// If any of the signers change their mind, they can create a new non-locktime transaction. The new transaction will use, as one of its inputs, one of the same outputs which was used as an input to the locktime transaction. This makes the locktime transaction invalid if the new transaction is added to the block chain before the time lock expires.

// Care must be taken near the expiry time of a time lock. The peer-to-peer network allows block time to be up to two hours ahead of real time, so a locktime transaction can be added to the block chain up to two hours before its time lock officially expires. Also, blocks are not created at guaranteed intervals, so any attempt to cancel a valuable transaction should be made a few hours before the time lock expires.

// Previous versions of Bitcoin Core provided a feature which prevented transaction signers from using the method described above to cancel a time-locked transaction, but a necessary part of this feature was disabled to prevent denial of service attacks. A legacy of this system are four-byte sequence numbers in every input. Sequence numbers were meant to allow multiple signers to agree to update a transaction; when they finished updating the transaction, they could agree to set every input’s sequence number to the four-byte unsigned maximum (0xffffffff), allowing the transaction to be added to a block even if its time lock had not expired.

// Even today, setting all sequence numbers to 0xffffffff (the default in Bitcoin Core) can still disable the time lock, so if you want to use locktime, at least one input must have a sequence number below the maximum. Since sequence numbers are not used by the network for any other purpose, setting any sequence number to zero is sufficient to enable locktime.

// Locktime itself is an unsigned 4-byte integer which can be parsed two ways:

// If less than 500 million, locktime is parsed as a block height. The transaction can be added to any block which has this height or higher.

// If greater than or equal to 500 million, locktime is parsed using the Unix epoch time format (the number of seconds elapsed since 1970-01-01T00:00 UTC—currently over 1.395 billion). The transaction can be added to any block whose block time is greater than the locktime.

import (
	"errors"

	"github.com/detailyang/go-bcore"
	"github.com/detailyang/go-bcrypto"
	. "github.com/detailyang/go-bprimitives"
)

type SigHash uint32

var (
	// the default, signs all the inputs and outputs, protecting everything except the signature scripts against modification.
	SigHashAll SigHash = 1
	// signs all of the inputs but none of the outputs, allowing anyone to change where the satoshis are going unless other signatures using other signature hash flags protect the outputs
	SigHashNone SigHash = 2
	// the only output signed is the one corresponding to this input (the output with the same output index number as this input), ensuring nobody can change your part of the transaction but allowing other signers to change their part of the transaction. The corresponding output must exist or the value “1” will be signed, breaking the security scheme. This input, as well as other inputs, are included in the signature. The sequence numbers of other inputs are not included in the signature, and can be updated.
	SigHashSingle       SigHash = 3
	SigHashForkId       SigHash = 0x40
	SigHashAnyoneCanPay SigHash = 0x80
)

func NewSigHash(s uint32) SigHash {
	return SigHash(s)
}

func (s *SigHash) Enable(flag SigHash) SigHash {
	*s = *s | flag
	return *s
}

func (s SigHash) Has(flag SigHash) bool {
	return (s & flag) == flag
}

func (s SigHash) HasBaseType(flag SigHash) bool {
	return s&0x1F == flag
}

type SignatureVersion uint32

const (
	SignatureVersionBase SignatureVersion = 1 << iota
	SignatureVersionWitnessV0
	SignatureVersionForkId
)

const (
	TransactionSignerLocktimeThreshold = 500000000
	// TransactionSignerLocktimeSequenceFinal Setting nSequence to this value for every input in a transaction
	// disables nLockTime.
	TransactionSignerLocktimeSequenceFinal = 0xffffffff

	// TransactionSequenceLocktimeDisableFlag apply in the context of BIP 68
	// If this flag set, sequence is NOT interpreted as a
	// relative lock-time.
	TransactionSequenceLocktimeDisableFlag = 1 << 31

	// If txinput->Sequence encodes a relative lock-time and this flag
	// is set, the relative lock-time has units of 512 seconds,
	// otherwise it specifies blocks with a granularity of 1.
	TransactionSequenceLockTimeTypeFlag = uint32(1 << 22)
	// If txinput->sequence encodes a relative lock-time, this mask is
	// applied to extract that lock-time from the sequence field.
	TransactionSequenceLocktimeMask = uint32(0x0000ffff)
)

var (
	ErrTransactionSignerLockTimeThreshold     = errors.New("transaction signer: locktime < threshold")
	ErrTransactionSignerLockTimeNotArrived    = errors.New("transaction signer: locktime has not arrived")
	ErrTransactionSignerLocktimeSequenceFinal = errors.New("transaction signer: final sequence")
	ErrTransactionSignerSequenceLowVersion    = errors.New("transaction signer: transaction version below 2")
	ErrTransactionSignerSequenceDisabled      = errors.New("transaction signer: sequence disabled")
	ErrTransactionSignerSequenceThresold      = errors.New("transaction signer: sequence < threshold")
	ErrTransactionSignerSequenceNotArrived    = errors.New("transaction signer: tosequnce < sequence")
	ErrTransactionSignerEmptySignature        = errors.New("transaction signer: zero signature")
	ErrTransactionSignerVerifySignatureFailed = errors.New("transaction signer: verify signature failed")
)

type TransactionSigner struct {
	Transaction *bcore.Transaction
	InputIndex  int
	InputValue  uint64
}

func NewTransactionSigner(tx *bcore.Transaction, InputIndex int, InputValue uint64) *TransactionSigner {
	return &TransactionSigner{
		Transaction: tx,
		InputIndex:  InputIndex,
		InputValue:  InputValue,
	}
}

func (ts *TransactionSigner) CheckLockTime(locktime uint32) error {
	// There are two kinds of nLockTime: lock-by-blockheight
	// and lock-by-blocktime, distinguished by whether
	// nLockTime < LOCKTIME_THRESHOLD.
	//
	// We want to compare apples to apples, so fail the script
	// unless the type of nLockTime being tested is the same as
	// the nLockTime in the transaction.

	if !((ts.Transaction.Locktime < TransactionSignerLocktimeThreshold &&
		locktime < TransactionSignerLocktimeThreshold) ||
		(ts.Transaction.Locktime >= TransactionSignerLocktimeThreshold &&
			locktime >= TransactionSignerLocktimeThreshold)) {
		return ErrTransactionSignerLockTimeThreshold
	}

	if locktime > ts.Transaction.Locktime {
		return ErrTransactionSignerLockTimeNotArrived
	}

	if bcore.TransactionFinalSequence != ts.Transaction.Inputs[ts.InputIndex].Sequence {
		return ErrTransactionSignerLocktimeSequenceFinal
	}

	return nil
}

func (ts *TransactionSigner) CheckSequence(sequence uint32) error {
	// Fail if the transaction's version number is not set high
	// enough to trigger BIP 68 rules.
	if ts.Transaction.Version < 2 {
		return ErrTransactionSignerSequenceLowVersion
	}

	if ts.Transaction.Inputs[ts.InputIndex].Sequence&TransactionSequenceLocktimeDisableFlag ==
		TransactionSequenceLocktimeDisableFlag {
		return ErrTransactionSignerSequenceDisabled
	}

	locktimeMask := TransactionSequenceLockTimeTypeFlag | TransactionSequenceLocktimeMask
	tsequence := ts.Transaction.Inputs[ts.InputIndex].Sequence & locktimeMask
	sequence = sequence & locktimeMask

	// There are two kinds of nSequence: lock-by-blockheight
	// and lock-by-blocktime, distinguished by whether
	// nSequenceMasked < TransactionSequenceLockTimeTypeFlag.

	// We want to compare apples to apples, so fail the script
	// unless the type of nSequenceMasked being tested is the same as
	// the nSequenceMasked in the transaction.
	if !((tsequence < TransactionSequenceLockTimeTypeFlag && locktimeMask < TransactionSequenceLockTimeTypeFlag) ||
		(tsequence >= TransactionSequenceLockTimeTypeFlag && locktimeMask >= TransactionSequenceLockTimeTypeFlag)) {
		return ErrTransactionSignerSequenceThresold
	}

	if tsequence < sequence {
		return ErrTransactionSignerSequenceNotArrived
	}

	return nil
}

func (ts *TransactionSigner) computeHashPrevOuts(sighash SigHash) Hash {
	if sighash.Has(SigHashAnyoneCanPay) {
		buffer := NewBuffer()
		for _, input := range ts.Transaction.Inputs {
			buffer.PutBytes(input.PrevOutput.Bytes())
		}

		return DHash256(buffer.Bytes())
	}

	return HashZero
}

func (ts *TransactionSigner) computeHashSequence(sighash SigHash) Hash {
	if sighash.Has(SigHashAll) && sighash.Has(SigHashAnyoneCanPay) {
		buffer := NewBuffer()
		for _, input := range ts.Transaction.Inputs {
			buffer.PutUint32(input.Sequence)
		}

		return DHash256(buffer.Bytes())
	}

	return HashZero
}

func (ts *TransactionSigner) computeHashOutputs(sighash SigHash) Hash {
	if sighash&SigHashAll == SigHashAll {
		buffer := NewBuffer()
		for _, output := range ts.Transaction.Outputs {
			buffer.PutBytes(output.Bytes())
		}

		return DHash256(buffer.Bytes())
	} else if sighash&SigHashSingle == SigHashSingle {
		if ts.InputIndex < len(ts.Transaction.Outputs) {
			return DHash256(ts.Transaction.Outputs[ts.InputIndex].Bytes())
		}
	}

	return HashZero
}

// Double SHA256 of the serialization of:
// 1. nVersion of the transaction (4-byte little endian)
// 2. hashPrevouts (32-byte hash)
// 3. hashSequence (32-byte hash)
// 4. outpoint (32-byte hash + 4-byte little endian)
// 5. scriptCode of the input (serialized as scripts inside CTxOuts)
// 6. value of the output spent by this input (8-byte little endian)
// 7. nSequence of the input (4-byte little endian)
// 8. hashOutputs (32-byte hash)
// 9. nLocktime of the transaction (4-byte little endian)
// 10. sighash type of the signature (4-byte little endian)
func (ts *TransactionSigner) signatureHashWitnessV0(script *Script, sighash SigHash) Hash {
	hashPrevouts := ts.computeHashPrevOuts(sighash)
	hashSequence := ts.computeHashSequence(sighash)
	hashOutputs := ts.computeHashOutputs(sighash)

	return DHash256(NewBuffer().PutUint32(ts.Transaction.Version).
		PutHash(hashPrevouts).
		PutHash(hashSequence).
		PutBytes(ts.Transaction.Inputs[ts.InputIndex].PrevOutput.Bytes()).
		PutVarBytes(script.Bytes()).
		PutUint64(ts.InputValue).
		PutUint32(ts.Transaction.Inputs[ts.InputIndex].Sequence).
		PutHash(hashOutputs).
		PutUint32(ts.Transaction.Locktime).
		PutUint32(uint32(sighash)).Bytes())
}

func (ts *TransactionSigner) signatureHashForkId(script *Script, sighash SigHash) Hash {
	if ts.InputIndex >= len(ts.Transaction.Inputs) {
		return HashOne
	}

	if sighash&0x1f == SigHashSingle {
		return HashOne
	}

	return ts.signatureHashWitnessV0(script, sighash)
}

func (ts *TransactionSigner) signatureHashOriginal(script *Script, sighash SigHash) Hash {
	if ts.InputIndex >= len(ts.Transaction.Inputs) {
		return HashOne
	}

	if sighash.HasBaseType(SigHashSingle) && ts.InputIndex >= len(ts.Transaction.Outputs) {
		return HashOne
	}

	script = script.WithoutSep()

	var inputs []*bcore.TransactionInput

	if sighash.Has(SigHashAnyoneCanPay) {
		input := ts.Transaction.Inputs[ts.InputIndex]
		inputs = []*bcore.TransactionInput{
			&bcore.TransactionInput{
				PrevOutput: input.PrevOutput,
				ScriptSig:  script.Bytes(),
				Sequence:   input.Sequence,
			},
		}

	} else {
		inputs = make([]*bcore.TransactionInput, len(ts.Transaction.Inputs))
		for i, input := range ts.Transaction.Inputs {
			scriptSig := []byte{}
			if i == ts.InputIndex {
				scriptSig = script.Bytes()
			}

			sequence := input.Sequence
			if i != ts.InputIndex && (sighash.HasBaseType(SigHashSingle) ||
				sighash.HasBaseType(SigHashNone)) {
				sequence = 0
			}

			inputs[i] = &bcore.TransactionInput{
				PrevOutput: input.PrevOutput.Clone(),
				ScriptSig:  scriptSig,
				Sequence:   sequence,
			}
		}
	}

	var outputs []*bcore.TransactionOutput

	switch sighash & 0x1f {
	case SigHashNone:
	case SigHashSingle:
		outputs = make([]*bcore.TransactionOutput, 0, len(ts.Transaction.Outputs))
		for i := 0; i < ts.InputIndex+1; i++ {
			output := ts.Transaction.Outputs[i]
			if i == ts.InputIndex {
				outputs = append(outputs, output.Clone())
			} else {
				outputs = append(outputs, bcore.NewDefaultTransactionOutput())
			}
		}
	default:
		outputs = make([]*bcore.TransactionOutput, len(ts.Transaction.Outputs))
		for i, output := range ts.Transaction.Outputs {
			outputs[i] = output.Clone()
		}
	}

	tx := &bcore.Transaction{
		Version:  ts.Transaction.Version,
		Inputs:   inputs,
		Outputs:  outputs,
		Locktime: ts.Transaction.Locktime,
	}

	return DHash256(append(tx.Bytes(), byte(sighash), byte(sighash>>8), byte(sighash>>16), byte(sighash>>24)))
}

func (ts *TransactionSigner) CheckSignature(sig, pubkey []byte, script *Script, version SignatureVersion) error {
	if len(sig) == 0 {
		return ErrTransactionSignerEmptySignature
	}

	sighash := NewSigHash(uint32(sig[len(sig)-1]))
	sig = sig[:len(sig)-1]

	var hash Hash

	switch version {
	case SignatureVersionBase:
		hash = ts.signatureHashOriginal(script, sighash)
	case SignatureVersionWitnessV0:
		hash = ts.signatureHashWitnessV0(script, sighash)
	case SignatureVersionForkId:
		hash = ts.signatureHashForkId(script, sighash)
	}

	pk, err := bcrypto.NewPubkeyFromBytes(pubkey)
	if err != nil {
		return err
	}

	signature, err := bcrypto.NewSignatureFromBytes(sig)
	if err != nil {
		return err
	}

	ok := signature.Verify(hash.Bytes(), pk)
	if !ok {
		return ErrTransactionSignerVerifySignatureFailed
	}

	return nil
}
