package bscript

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/detailyang/go-bcore"
	bcrypto "github.com/detailyang/go-bcrypto"

	. "github.com/detailyang/go-bprimitives"
)

type TestBuilder struct {
	script       *Script
	redeemscript *Script
	comment      string
	amount       uint64
	flag         Flag

	scriptError error

	havePush bool
	push     []byte

	creditTx *bcore.Transaction
	spendTx  *bcore.Transaction
}

func NewCreditingTransaction(script *Script, value uint64) *bcore.Transaction {
	var tx bcore.Transaction
	tx.Version = 1
	tx.Locktime = 0
	tx.Inputs = make([]*bcore.TransactionInput, 1)
	tx.Outputs = make([]*bcore.TransactionOutput, 1)
	var input bcore.TransactionInput
	var output bcore.TransactionOutput
	input.PrevOutput = bcore.NewDefaultOutPoint()
	input.ScriptSig = NewScript().PushNumber(0).PushNumber(0).Bytes()
	input.Sequence = bcore.TransactionFinalSequence
	output.ScriptPubkey = script.Bytes()
	output.Value = value

	tx.Inputs[0] = &input
	tx.Outputs[0] = &output

	return &tx
}

func NewSpendingTransaction(script *Script, creditTx *bcore.Transaction) *bcore.Transaction {
	var tx bcore.Transaction
	tx.Version = 1
	tx.Locktime = 0
	tx.Inputs = make([]*bcore.TransactionInput, 1)
	tx.Outputs = make([]*bcore.TransactionOutput, 1)
	var input bcore.TransactionInput
	var output bcore.TransactionOutput
	input.PrevOutput = bcore.NewOutPoint(creditTx.ID(), 0)
	input.ScriptSig = script.Bytes()
	input.Sequence = bcore.TransactionFinalSequence
	output.ScriptPubkey = make([]byte, 0)
	output.Value = creditTx.Outputs[0].Value

	tx.Inputs[0] = &input
	tx.Outputs[0] = &output

	return &tx
}

func NewTestBuilder(script *Script, comment string, flag Flag, P2SH bool, amount uint64) *TestBuilder {
	var redeemscript *Script
	scriptPubkey := script

	if P2SH {
		redeemscript = scriptPubkey
		scriptPubkey = NewScript().PushOPCode(OP_HASH160).
			PushBytesWithOP(Hash160(redeemscript.Bytes())).
			PushOPCode(OP_EQUAL)
	}

	// OP_HASH160 [20-byte-hash-value] OP_EQUAL
	creditTx := NewCreditingTransaction(scriptPubkey, amount)
	spendTx := NewSpendingTransaction(NewScript(), creditTx)

	return &TestBuilder{
		spendTx:      spendTx,
		creditTx:     creditTx,
		amount:       amount,
		flag:         flag,
		script:       script,
		comment:      comment,
		redeemscript: redeemscript,
	}
}

func clone(s []byte) []byte {
	h := make([]byte, len(s))
	copy(h, s)
	return h
}

// 30 + 02 + len(r) + r + 02 + len(s) + s
func negateSigantureS(sig []byte) []byte {
	r := clone(sig[4 : 4+sig[3]])
	s := clone(sig[6+sig[3] : 6+sig[3]+sig[5+sig[3]]])

	for len(s) < 33 {
		s = append([]byte{0x00}, s...)
	}

	order := []byte{
		0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0xBA, 0xAE, 0xDC, 0xE6, 0xAF,
		0x48, 0xA0, 0x3B, 0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41, 0x41}

	carry := 0
	for p := 32; p >= 1; p-- {
		n := int(order[p]) - int(s[p]) - carry
		s[p] = byte(int(n+256) & 0xFF)
		if n < 0 {
			carry = 1
		} else {
			carry = 0
		}
	}

	if len(s) > 1 && s[0] == 0 && s[1] < 0x80 {
		s = s[1:]
	}

	newsig := make([]byte, 0, len(sig))
	newsig = append(newsig, 0x30)
	newsig = append(newsig, byte(4+len(r)+len(s)))
	newsig = append(newsig, 0x02)
	newsig = append(newsig, byte(len(r)))
	newsig = append(newsig, r...)
	newsig = append(newsig, 0x02)
	newsig = append(newsig, byte(len(s)))
	newsig = append(newsig, s...)

	return newsig
}

func (tb *TestBuilder) Num(n int) *TestBuilder {
	tb.DoPush()
	tb.spendTx.Inputs[0].ScriptSig = NewScriptFromBytes(
		tb.spendTx.Inputs[0].ScriptSig,
	).PushInt64(int64(n)).Bytes()
	return tb
}

func (tb *TestBuilder) PushPubkey(pubkey *bcrypto.PublicKey) *TestBuilder {
	return tb.DoPushBytes(pubkey.Bytes())
}

func (tb *TestBuilder) PushHex(hexstring string) *TestBuilder {
	data, err := hex.DecodeString(hexstring)
	if err != nil {
		panic(err)
	}

	return tb.DoPushBytes(data)
}

func (tb *TestBuilder) Push(publickey *bcrypto.PublicKey) *TestBuilder {
	tb.DoPushBytes(publickey.Bytes())
	return tb
}

func (tb *TestBuilder) DoPush() {
	if tb.havePush {
		tb.spendTx.Inputs[0].ScriptSig = NewScriptFromBytes(tb.spendTx.Inputs[0].ScriptSig).
			PushBytesWithOP(tb.push).
			Bytes()
		tb.havePush = false
	}
}

func (tb *TestBuilder) EditPush(pos int, hexin, hexout string) *TestBuilder {
	datain, _ := hex.DecodeString(hexin)
	dataout, _ := hex.DecodeString(hexout)

	if !bytes.Equal(tb.push[pos:pos+len(datain)], datain) {
		panic(tb.comment)
	}

	left := clone(tb.push[:pos])
	right := clone(tb.push[pos+len(datain):])

	tb.push = append(left, dataout...)
	tb.push = append(tb.push, right...)

	return tb
}

func (tb *TestBuilder) Add(script *Script) *TestBuilder {
	tb.DoPush()
	tb.spendTx.Inputs[0].ScriptSig = NewScriptFromBytes(tb.spendTx.Inputs[0].ScriptSig).
		PushBytes(script.Bytes()).Bytes()
	return tb
}

func (tb *TestBuilder) DoPushBytes(data []byte) *TestBuilder {
	tb.DoPush()
	tb.push = data
	tb.havePush = true

	return tb
}

func (tb *TestBuilder) PushRedeem() *TestBuilder {
	tb.DoPushBytes(tb.redeemscript.Bytes())
	return tb
}

func (tb *TestBuilder) PushSig(
	key *bcrypto.Key,
	sigHash SigHash,
	nr, ns int,
	amount uint64,
	flag Flag,
) *TestBuilder {
	ts := NewTransactionSigner(tb.spendTx, 0, amount)
	hash := ts.SiagntureHash(tb.script, sigHash, flag)
	sig := tb.DoSign(key, hash, nr, ns)
	sig = append(sig, byte(sigHash))

	tb.DoPushBytes(sig)

	return tb
}

func (tb *TestBuilder) DamagePush(pos int) *TestBuilder {
	tb.push[pos] ^= 1
	return tb
}

func (tb *TestBuilder) ScriptError(e error) *TestBuilder {
	tb.scriptError = e
	return tb
}

func DoTest(
	t *testing.T,
	scriptPubkey *Script,
	scriptSig *Script,
	flag Flag,
	message string,
	scriptError error,
	value uint64,
) {
	creditTx := NewCreditingTransaction(
		scriptPubkey, value,
	)
	spentTx := NewSpendingTransaction(scriptSig, creditTx)

	err := VerifyScript(scriptSig, scriptPubkey,
		NewScriptWitness(nil),
		flag,
		NewTransactionSigner(spentTx, 0, creditTx.Outputs[0].Value),
		SignatureVersionBase,
	)

	if err != scriptError {
		t.Errorf("%s - got %s", message, err)
	}
}

func (tb *TestBuilder) Test(t *testing.T) {
	tb.DoPush()
	DoTest(
		t,
		NewScriptFromBytes(tb.creditTx.Outputs[0].ScriptPubkey),
		NewScriptFromBytes(tb.spendTx.Inputs[0].ScriptSig),
		tb.flag,
		tb.comment,
		tb.scriptError,
		tb.amount,
	)
}

func (tb *TestBuilder) DoSign(key *bcrypto.Key, hash Hash, nr, ns int) []byte {
	var sig, r, s []byte

	for i := 0; ; i++ {
		sig, _ = key.Signature(hash.Bytes(), uint32(i))

		if (ns == 33) != (sig[5]+sig[3] == 33) {
			sig = negateSigantureS(sig)
		}

		r = clone(sig[4 : 4+sig[3]])
		s = clone(sig[6+sig[3] : 6+sig[3]+sig[5+sig[3]]])

		if nr == len(r) && ns == len(s) {
			break
		}
	}

	return sig
}

var (
	key0bytes = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

	key1bytes = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}

	key2bytes = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0}
)

func TestScriptFlag(t *testing.T) {
	key0 := bcrypto.NewKey(key0bytes[:], false)
	key0c := bcrypto.NewKey(key0bytes[:], true)
	pubkey0, _ := key0.GetPubkey()
	pubkey0h, _ := key0.GetPubkey()
	pubkey0h[0] = 0x06 | (pubkey0h[64] & 1)
	pubkey0c, _ := key0c.GetPubkey()

	key1 := bcrypto.NewKey(key1bytes[:], false)
	key1c := bcrypto.NewKey(key1bytes[:], true)
	pubkey1, _ := key1.GetPubkey()
	pubkey1c, _ := key1c.GetPubkey()

	key2 := bcrypto.NewKey(key2bytes[:], false)
	key2c := bcrypto.NewKey(key2bytes[:], true)
	// pubkey2, _ := key1.GetPubkey()
	pubkey2c, _ := key2c.GetPubkey()

	flag := ScriptEnableSigHashForkID

	tests := make([]*TestBuilder, 0, 512)

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey0.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK",
		0,
		false,
		0,
	).PushSig(key0, SigHashAll, 32, 32, 0, flag))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey0.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK, bad sig",
		0,
		false,
		0,
	).PushSig(key0, SigHashAll, 32, 32, 0, flag).
		DamagePush(10).
		ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_DUP).PushOPCode(OP_HASH160).
			PushBytesWithOP(pubkey1c.ID()).
			PushOPCode(OP_EQUALVERIFY).
			PushOPCode(OP_CHECKSIG),
		"P2PKH",
		0,
		false,
		0,
	).PushSig(key1, SigHashAll, 32, 32, 0, flag).
		Push(&pubkey1c))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_DUP).PushOPCode(OP_HASH160).
			PushBytesWithOP(pubkey2c.ID()).
			PushOPCode(OP_EQUALVERIFY).
			PushOPCode(OP_CHECKSIG),
		"P2PKH, bad pubkey",
		0,
		false,
		0,
	).PushSig(key2, SigHashAll, 32, 32, 0, flag).Push(&pubkey2c).
		DamagePush(5).ScriptError(ErrInterpreterVerifyFailed))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1.Bytes()).
			PushOPCode(OP_CHECKSIG),
		"P2PK anyonecanpy",
		0,
		false,
		0,
	).PushSig(key1, SigHashAll|SigHashAnyoneCanPay, 32, 32, 0, flag))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1.Bytes()).
			PushOPCode(OP_CHECKSIG),
		"P2PK anyonecanpy marked with normal hashtyp",
		0,
		false,
		0,
	).PushSig(key1, SigHashAll|SigHashAnyoneCanPay, 32, 32, 0, flag).
		EditPush(70, "81", "01").ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey0c.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2SH(P2PK)",
		ScriptVerifyP2SH,
		true,
		0).PushSig(key0, SigHashAll, 32, 32, 0, flag).PushRedeem())

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey0c.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2SH(P2PK)",
		ScriptVerifyP2SH,
		true,
		0).PushSig(key0, SigHashAll, 32, 32, 0, flag).PushRedeem().DamagePush(10).
		ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_DUP).PushOPCode(OP_HASH160).
			PushBytesWithOP(pubkey0.ID()).PushOPCode(OP_EQUALVERIFY).PushOPCode(OP_CHECKSIG),
		"P2SH(P2PK)",
		ScriptVerifyP2SH,
		true,
		0).PushSig(key0, SigHashAll, 32, 32, 0, flag).Push(&pubkey0).PushRedeem())

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_DUP).PushOPCode(OP_HASH160).
			PushBytesWithOP(pubkey1.ID()).PushOPCode(OP_EQUALVERIFY).PushOPCode(OP_CHECKSIG),
		"P2SH(P2PK), bad sig but no VERIFY_P2SH",
		0,
		true,
		0).PushSig(key0, SigHashAll, 32, 32, 0, flag).DamagePush(10).PushRedeem())

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_DUP).PushOPCode(OP_HASH160).
			PushBytesWithOP(pubkey1.ID()).PushOPCode(OP_EQUALVERIFY).PushOPCode(OP_CHECKSIG),
		"P2SH(P2PK), bad sig",
		ScriptVerifyP2SH,
		true,
		0).PushSig(key0, SigHashAll, 32, 32, 0, flag).DamagePush(10).PushRedeem().
		ScriptError(ErrInterpreterVerifyFailed))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_3).PushBytesWithOP(pubkey0c.Bytes()).
			PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_3).PushOPCode(OP_CHECKMULTISIG),
		"3-of-3",
		0,
		false,
		0,
	).Num(0).
		PushSig(key0, SigHashAll, 32, 32, 0, flag).
		PushSig(key1, SigHashAll, 32, 32, 0, flag).
		PushSig(key2, SigHashAll, 32, 32, 0, flag))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_3).PushBytesWithOP(pubkey0c.Bytes()).
			PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_3).PushOPCode(OP_CHECKMULTISIG),
		"3-of-3, 2 sigs",
		0,
		false,
		0,
	).Num(0).
		PushSig(key0, SigHashAll, 32, 32, 0, flag).
		PushSig(key1, SigHashAll, 32, 32, 0, flag).
		Num(0).
		ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).PushBytesWithOP(pubkey0c.Bytes()).
			PushBytesWithOP(pubkey1c.Bytes()).
			PushBytesWithOP(pubkey2c.Bytes()).PushOPCode(OP_3).
			PushOPCode(OP_CHECKMULTISIG),
		"P2SH(2-of-3)",
		ScriptVerifyP2SH,
		true,
		0,
	).Num(0).
		PushSig(key1, SigHashAll, 32, 32, 0, flag).
		PushSig(key2, SigHashAll, 32, 32, 0, flag).
		PushRedeem())

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).PushBytesWithOP(pubkey0c.Bytes()).
			PushBytesWithOP(pubkey1c.Bytes()).
			PushBytesWithOP(pubkey2c.Bytes()).PushOPCode(OP_3).
			PushOPCode(OP_CHECKMULTISIG),
		"P2SH(2-of-3), 1 sig",
		ScriptVerifyP2SH,
		true,
		0,
	).Num(0).
		PushSig(key1, SigHashAll, 32, 32, 0, flag).
		Num(0).
		PushRedeem().
		ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with too much R padding but no DERSIG",
		0,
		false,
		0,
	).PushSig(key1, SigHashAll, 31, 32, 0, flag).EditPush(1, "43021F", "44022000"))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with too much R padding",
		ScriptVerifyDERSignatures,
		false,
		0,
	).PushSig(key1, SigHashAll, 31, 32, 0, flag).EditPush(1, "43021F", "44022000").ScriptError(ErrInterpreterBadSignatureDer))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with too much s padding but no DERSIG",
		0,
		false,
		0,
	).PushSig(key1, SigHashAll, 32, 32, 0, flag).
		EditPush(1, "44", "45").
		EditPush(37, "20", "2100"))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with too much s padding but no DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).PushSig(key1, SigHashAll, 32, 32, 0, flag).
		EditPush(1, "44", "45").
		EditPush(37, "20", "2100").ScriptError(ErrInterpreterBadSignatureDer))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with too little R padding but no DERSIG",
		0,
		false,
		0,
	).PushSig(key1, SigHashAll, 33, 32, 0, flag).EditPush(1, "45022100", "440220"))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with too little R padding but no DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).PushSig(key1, SigHashAll, 33, 32, 0, flag).EditPush(1, "45022100", "440220").
		ScriptError(ErrInterpreterBadSignatureDer))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey2c.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"P2PK NOT with bad sig with too much R padding but no DERSIG",
		0,
		false,
		0,
	).PushSig(key2, SigHashAll, 31, 32, 0, flag).
		EditPush(1, "43021f", "44022000").
		DamagePush(10))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey2c.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"P2PK NOT with bad sig with too much R padding",
		ScriptVerifyDERSignatures,
		false,
		0,
	).PushSig(key2, SigHashAll, 31, 32, 0, flag).
		EditPush(1, "43021f", "44022000").
		DamagePush(10).ScriptError(ErrInterpreterBadSignatureDer))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey2c.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"P2PK NOT with bad sig with too much R padding but no DERSIG",
		0,
		false,
		0,
	).PushSig(key2, SigHashAll, 31, 32, 0, flag).
		EditPush(1, "43021f", "44022000").
		ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey2c.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"P2PK NOT with bad sig with too much R padding",
		ScriptVerifyDERSignatures,
		false,
		0,
	).PushSig(key2, SigHashAll, 31, 32, 0, flag).
		EditPush(1, "43021f", "44022000").
		ScriptError(ErrInterpreterBadSignatureDer))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG),
		"BIP66 example 1, without DERSIG",
		0,
		false,
		0,
	).PushSig(key1, SigHashAll, 33, 32, 0, flag).
		EditPush(1, "45022100", "440220"))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG),
		"BIP66 example 1, with DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).PushSig(key1, SigHashAll, 33, 32, 0, flag).
		EditPush(1, "45022100", "440220").ScriptError(ErrInterpreterBadSignatureDer))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"BIP66 example 2, without DERSIG",
		0,
		false,
		0,
	).PushSig(key1, SigHashAll, 33, 32, 0, flag).
		EditPush(1, "45022100", "440220").ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"BIP66 example 2, with DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).PushSig(key1, SigHashAll, 33, 32, 0, flag).
		EditPush(1, "45022100", "440220").ScriptError(ErrInterpreterBadSignatureDer))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG),
		"BIP66 example 3, without DERSIG",
		0,
		false,
		0,
	).Num(0).ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG),
		"BIP66 example 3, with DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).Num(0).ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"BIP66 example 4, without DERSIG",
		0,
		false,
		0,
	).Num(0))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"BIP66 example 4, with DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).Num(0))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"BIP66 example 4, without DERSIG, non-null DER-compliant signature",
		0,
		false,
		0,
	).PushHex("300602010102010101"))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"BIP66 example 4, without DERSIG, non-null DER-compliant signature",
		ScriptVerifyDERSignatures|ScriptVerifyNullFail,
		false,
		0,
	).Num(0))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"BIP66 example 4, with DERSIG, non-null DER-compliant signature",
		ScriptVerifyDERSignatures|ScriptVerifyNullFail,
		false,
		0,
	).PushHex("300602010102010101").ScriptError(ErrInterpreterSignatureNullFail))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG),
		"BIP66 example 5, without DERSIG",
		0,
		false,
		0,
	).Num(0).ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG),
		"BIP66 example 5, with DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).Num(1).ScriptError(ErrInterpreterBadSignatureDer))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"BIP66 example 6, without DERSIG",
		0,
		false,
		0,
	).Num(1))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"BIP66 example 6, without DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).Num(1).ScriptError(ErrInterpreterBadSignatureDer))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_2).
			PushOPCode(OP_CHECKMULTISIG),
		"BIP66 example 7, without DERSIG",
		0,
		false,
		0,
	).Num(0).PushSig(key1, SigHashAll, 33, 32, 0, flag).EditPush(1, "45022100", "440220").
		PushSig(key2, SigHashAll, 32, 32, 0, flag))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_2).
			PushOPCode(OP_CHECKMULTISIG),
		"BIP66 example 7, without DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).Num(0).PushSig(key1, SigHashAll, 33, 32, 0, flag).EditPush(1, "45022100", "440220").
		PushSig(key2, SigHashAll, 32, 32, 0, flag).ScriptError(ErrInterpreterBadSignatureDer))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_2).
			PushOPCode(OP_CHECKMULTISIG).PushOPCode(OP_NOT),
		"BIP66 example 8, without DERSIG",
		0,
		false,
		0,
	).Num(0).PushSig(key1, SigHashAll, 33, 32, 0, flag).EditPush(1, "45022100", "440220").
		PushSig(key2, SigHashAll, 32, 32, 0, flag).ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_2).
			PushOPCode(OP_CHECKMULTISIG).PushOPCode(OP_NOT),
		"BIP66 example 8, without DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).Num(0).PushSig(key1, SigHashAll, 33, 32, 0, flag).EditPush(1, "45022100", "440220").
		PushSig(key2, SigHashAll, 32, 32, 0, flag).ScriptError(ErrInterpreterBadSignatureDer))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_2).
			PushOPCode(OP_CHECKMULTISIG),
		"BIP66 example 9, without DERSIG",
		0,
		false,
		0,
	).Num(0).Num(0).PushSig(key2, SigHashAll, 33, 32, 0, flag).EditPush(1, "45022100", "440220").
		ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_2).
			PushOPCode(OP_CHECKMULTISIG),
		"BIP66 example 9, with DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).Num(0).Num(0).PushSig(key2, SigHashAll, 33, 32, 0, flag).EditPush(1, "45022100", "440220").
		ScriptError(ErrInterpreterBadSignatureDer))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_2).
			PushOPCode(OP_CHECKMULTISIG).PushOPCode(OP_NOT),
		"BIP66 example 10, without DERSIG",
		0,
		false,
		0,
	).Num(0).Num(0).PushSig(key2, SigHashAll, 33, 32, 0, flag).EditPush(1, "45022100", "440220"))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_2).
			PushOPCode(OP_CHECKMULTISIG).PushOPCode(OP_NOT),
		"BIP66 example 10, without DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).Num(0).Num(0).PushSig(key2, SigHashAll, 33, 32, 0, flag).EditPush(1, "45022100", "440220").
		ScriptError(ErrInterpreterBadSignatureDer))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_2).
			PushOPCode(OP_CHECKMULTISIG),
		"BIP66 example 11, without DERSIG",
		0,
		false,
		0,
	).Num(0).PushSig(key1, SigHashAll, 33, 32, 0, flag).EditPush(1, "45022100", "440220").Num(0).
		ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_2).
			PushOPCode(OP_CHECKMULTISIG),
		"BIP66 example 11, with DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).Num(0).PushSig(key1, SigHashAll, 33, 32, 0, flag).EditPush(1, "45022100", "440220").Num(0).
		ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_2).
			PushOPCode(OP_CHECKMULTISIG).PushOPCode(OP_NOT),
		"BIP66 example 12, without DERSIG",
		0,
		false,
		0,
	).Num(0).PushSig(key1, SigHashAll, 33, 32, 0, flag).EditPush(1, "45022100", "440220").Num(0))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_2).
			PushOPCode(OP_CHECKMULTISIG).PushOPCode(OP_NOT),
		"BIP66 example 12, without DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).Num(0).PushSig(key1, SigHashAll, 33, 32, 0, flag).EditPush(1, "45022100", "440220").Num(0))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey2c.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with multi-byte hashtype, without DERSIG",
		0,
		false,
		0,
	).Num(0).PushSig(key2, SigHashAll, 32, 32, 0, flag).EditPush(70, "01", "0101"))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey2c.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with multi-byte hashtype, with DERSIG",
		ScriptVerifyDERSignatures,
		false,
		0,
	).Num(0).PushSig(key2, SigHashAll, 32, 32, 0, flag).EditPush(70, "01", "0101").ScriptError(
		ErrInterpreterBadSignatureDer))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey2c.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with high S but no LOW_S",
		0,
		false,
		0,
	).PushSig(key2, SigHashAll, 32, 33, 0, flag))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey2c.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with high S",
		ScriptVerifyLowS,
		false,
		0,
	).PushSig(key2, SigHashAll, 32, 33, 0, flag).ScriptError(ErrInterpreterSigantureHighS))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey0h.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with hybrid pubkey but no STRICTENC",
		0,
		false,
		0,
	).PushSig(key0, SigHashAll, 32, 32, 0, flag))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey0h.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with hybrid pubkey",
		ScriptVerifyStrictEncoding,
		false,
		0,
	).PushSig(key0, SigHashAll, 32, 32, 0, flag).ScriptError(ErrInterpreterBadPubkey))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey0h.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"P2PK with hybrid pubkey but no STRICTENC",
		0,
		false,
		0,
	).PushSig(key0, SigHashAll, 32, 32, 0, flag).ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey0h.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"P2PK with hybrid pubkey",
		ScriptVerifyStrictEncoding,
		false,
		0,
	).PushSig(key0, SigHashAll, 32, 32, 0, flag).ScriptError(ErrInterpreterBadPubkey))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey0h.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"P2PK with invalid hybrid pubkey but no STRICTENC",
		0,
		false,
		0,
	).PushSig(key0, SigHashAll, 32, 32, 0, flag).DamagePush(10))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey0h.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"P2PK with invalid hybrid pubkey",
		ScriptVerifyStrictEncoding,
		false,
		0,
	).PushSig(key0, SigHashAll, 32, 32, 0, flag).DamagePush(10).ScriptError(ErrInterpreterBadPubkey))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_1).PushBytesWithOP(pubkey0h.Bytes()).PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_2).PushOPCode(OP_CHECKMULTISIG),
		"1-of-2 with the second 1 hybrid pubkey and no STRICTENC",
		0,
		false,
		0,
	).Num(0).PushSig(key0, SigHashAll, 32, 32, 0, flag))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_1).PushBytesWithOP(pubkey0h.Bytes()).PushBytesWithOP(pubkey1c.Bytes()).PushOPCode(OP_2).PushOPCode(OP_CHECKMULTISIG),
		"1-of-2 with the second 1 hybrid pubkey",
		ScriptVerifyStrictEncoding,
		false,
		0,
	).Num(0).PushSig(key1, SigHashAll, 32, 32, 0, flag))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_1).PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey0h.Bytes()).PushOPCode(OP_2).PushOPCode(OP_CHECKMULTISIG),
		"1-of-2 with the second 1 hybrid pubkey",
		ScriptVerifyStrictEncoding,
		false,
		0,
	).Num(0).PushSig(key1, SigHashAll, 32, 32, 0, flag).ScriptError(ErrInterpreterBadPubkey))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with undefined hashtype but no STRICTENC",
		0,
		false,
		0,
	).PushSig(key1, NewSigHash(5), 32, 32, 0, flag))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with undefined hashtype",
		ScriptVerifyStrictEncoding,
		false,
		0,
	).PushSig(key1, NewSigHash(5), 32, 32, 0, flag).ScriptError(ErrInterpreterBadSignatureHashType))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_DUP).PushOPCode(OP_HASH160).
			PushBytesWithOP(pubkey0.ID()).PushOPCode(OP_EQUALVERIFY).PushOPCode(OP_CHECKSIG),
		"P2PK with undefined hashtype",
		0,
		false,
		0,
	).PushSig(key0, NewSigHash(0x21), 32, 32, 0, 0).PushPubkey(&pubkey0))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_DUP).PushOPCode(OP_HASH160).
			PushBytesWithOP(pubkey0.ID()).PushOPCode(OP_EQUALVERIFY).PushOPCode(OP_CHECKSIG),
		"P2PK with undefined hashtype",
		ScriptVerifyStrictEncoding,
		false,
		0,
	).PushSig(key0, NewSigHash(0x21), 32, 32, 0, ScriptVerifyStrictEncoding).
		PushPubkey(&pubkey0).ScriptError(ErrInterpreterBadSignatureHashType))

	tests = append(tests, NewTestBuilder(
		NewScript().
			PushBytesWithOP(pubkey1.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2SH(P2PK) with invalid sighashtype",
		ScriptVerifyP2SH,
		true,
		0,
	).PushSig(key1, NewSigHash(0x21), 32, 32, 0, flag).PushRedeem())

	tests = append(tests, NewTestBuilder(
		NewScript().
			PushBytesWithOP(pubkey1.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2SH(P2PK) with invalid sighashtype",
		ScriptVerifyP2SH|ScriptVerifyStrictEncoding,
		true,
		0,
	).PushSig(key1, NewSigHash(0x21), 32, 32, 0, flag).PushRedeem().
		ScriptError(ErrInterpreterBadSignatureHashType))

	tests = append(tests, NewTestBuilder(
		NewScript().
			PushBytesWithOP(pubkey1.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"P2PK NOT with invalid sig and undefined hashtype but no STRICTENC",
		0,
		false,
		0,
	).PushSig(key1, NewSigHash(5), 32, 32, 0, flag).
		DamagePush(10))

	tests = append(tests, NewTestBuilder(
		NewScript().
			PushBytesWithOP(pubkey1.Bytes()).PushOPCode(OP_CHECKSIG).PushOPCode(OP_NOT),
		"P2PK NOT with invalid sig and undefined hashtype",
		ScriptVerifyStrictEncoding,
		false,
		0,
	).PushSig(key1, NewSigHash(5), 32, 32, 0, flag).
		DamagePush(10).ScriptError(ErrInterpreterBadSignatureHashType))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_3).PushBytesWithOP(pubkey0c.Bytes()).
			PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_3).PushOPCode(OP_CHECKMULTISIG),
		"3-of-3 with nonzero dummy but no NULLDUMMY",
		0,
		false,
		0,
	).Num(1).
		PushSig(key0, SigHashAll, 32, 32, 0, flag).
		PushSig(key1, SigHashAll, 32, 32, 0, flag).
		PushSig(key2, SigHashAll, 32, 32, 0, flag))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_3).PushBytesWithOP(pubkey0c.Bytes()).
			PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_3).PushOPCode(OP_CHECKMULTISIG),
		"3-of-3 with nonzero dummy",
		ScriptVerifyNullDummy,
		false,
		0,
	).Num(1).
		PushSig(key0, SigHashAll, 32, 32, 0, flag).
		PushSig(key1, SigHashAll, 32, 32, 0, flag).
		PushSig(key2, SigHashAll, 32, 32, 0, flag).
		ScriptError(ErrInterpreterSignatureNullDummy))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_3).PushBytesWithOP(pubkey0c.Bytes()).
			PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_3).PushOPCode(OP_CHECKMULTISIG).PushOPCode(OP_NOT),
		"3-of-3 NOT with invalid sig and nonzero dummy but no NULLDUMMY",
		0,
		false,
		0,
	).Num(1).
		PushSig(key0, SigHashAll, 32, 32, 0, flag).
		PushSig(key1, SigHashAll, 32, 32, 0, flag).
		PushSig(key2, SigHashAll, 32, 32, 0, flag).
		DamagePush(10))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_3).PushBytesWithOP(pubkey0c.Bytes()).
			PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_3).PushOPCode(OP_CHECKMULTISIG).PushOPCode(OP_NOT),
		"3-of-3 NOT with invalid sig and nonzero dummy but no NULLDUMMY",
		ScriptVerifyNullDummy,
		false,
		0,
	).Num(1).
		PushSig(key0, SigHashAll, 32, 32, 0, flag).
		PushSig(key1, SigHashAll, 32, 32, 0, flag).
		PushSig(key2, SigHashAll, 32, 32, 0, flag).
		DamagePush(10).
		ScriptError(ErrInterpreterSignatureNullDummy))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).
			PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey1c.Bytes()).
			PushOPCode(OP_2).PushOPCode(OP_CHECKMULTISIG),
		"2-of-2 with two identical keys and sigs pushed using OP_DUP but no SIGPUSHONLY",
		0,
		false,
		0,
	).Num(0).
		PushSig(key1, SigHashAll, 32, 32, 0, flag).
		Add(NewScript().PushOPCode(OP_DUP)))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).
			PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey1c.Bytes()).
			PushOPCode(OP_2).PushOPCode(OP_CHECKMULTISIG),
		"2-of-2 with two identical keys and sigs pushed using OP_DUP",
		ScriptVerifySigPushOnly,
		false,
		0,
	).Num(0).
		PushSig(key1, SigHashAll, 32, 32, 0, flag).
		Add(NewScript().PushOPCode(OP_DUP)).ScriptError(ErrInterpreterSignaturePushOnly))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_CHECKSIG),
		"P2SH(P2PK) with non-push scriptSig but no P2SH or SIGPUSHONLY",
		0,
		true,
		0,
	).PushSig(key2, SigHashAll, 32, 32, 0, flag).
		Add(NewScript().PushOPCode(OP_NOP8)).PushRedeem())

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_CHECKSIG),
		"P2PK with non-push scriptSig but with P2SH validation",
		0,
		false,
		0,
	).PushSig(key2, SigHashAll, 32, 32, 0, flag).
		Add(NewScript().PushOPCode(OP_NOP8)))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_CHECKSIG),
		"P2SH(P2PK) with non-push scriptSig but no SIGPUSHONLY",
		ScriptVerifyP2SH,
		true,
		0,
	).PushSig(key2, SigHashAll, 32, 32, 0, flag).
		Add(NewScript().PushOPCode(OP_NOP8)).PushRedeem().ScriptError(ErrInterpreterSignaturePushOnly))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey2c.Bytes()).
			PushOPCode(OP_CHECKSIG),
		"P2SH(P2PK) with non-push scriptSig but not P2SH",
		ScriptVerifySigPushOnly,
		true,
		0,
	).PushSig(key2, SigHashAll, 32, 32, 0, flag).
		Add(NewScript().PushOPCode(OP_NOP8)).PushRedeem().ScriptError(ErrInterpreterSignaturePushOnly))

	tests = append(tests, NewTestBuilder(
		NewScript().PushOPCode(OP_2).
			PushBytesWithOP(pubkey1c.Bytes()).PushBytesWithOP(pubkey1c.Bytes()).
			PushOPCode(OP_2).PushOPCode(OP_CHECKMULTISIG),
		"2-of-2 with two identical keys and sigs pushed",
		ScriptVerifySigPushOnly,
		false,
		0,
	).Num(0).PushSig(key1, SigHashAll, 32, 32, 0, flag).
		PushSig(key1, SigHashAll, 32, 32, 0, flag))

	tests = append(tests, NewTestBuilder(
		NewScript().
			PushBytesWithOP(pubkey0.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with unnecessary input but no CLEANSTACK",
		ScriptVerifyP2SH,
		false,
		0,
	).Num(11).PushSig(key0, SigHashAll, 32, 32, 0, flag))

	tests = append(tests, NewTestBuilder(
		NewScript().
			PushBytesWithOP(pubkey0.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with unnecessary input but no CLEANSTACK",
		ScriptVerifyP2SH|ScriptVerifyCleanStack,
		false,
		0,
	).Num(11).PushSig(key0, SigHashAll, 32, 32, 0, flag).
		ScriptError(ErrInterpreterCleanStack))

	tests = append(tests, NewTestBuilder(
		NewScript().
			PushBytesWithOP(pubkey0.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with unnecessary input but no CLEANSTACK",
		ScriptVerifyP2SH,
		true,
		0,
	).Num(11).PushSig(key0, SigHashAll, 32, 32, 0, flag).PushRedeem())

	tests = append(tests, NewTestBuilder(
		NewScript().
			PushBytesWithOP(pubkey0.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with unnecessary input",
		ScriptVerifyP2SH|ScriptVerifyCleanStack,
		true,
		0,
	).Num(11).PushSig(key0, SigHashAll, 32, 32, 0, flag).PushRedeem().ScriptError(
		ErrInterpreterCleanStack,
	))

	tests = append(tests, NewTestBuilder(
		NewScript().
			PushBytesWithOP(pubkey0.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK with CLEANSTACK",
		ScriptVerifyP2SH|ScriptVerifyCleanStack,
		true,
		0,
	).PushSig(key0, SigHashAll, 32, 32, 0, flag).PushRedeem())

	value := uint64(12345000000000)

	tests = append(tests, NewTestBuilder(
		NewScript().
			PushBytesWithOP(pubkey0.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK FORKID",
		ScriptEnableSigHashForkID,
		false,
		value,
	).PushSig(key0, SigHashAll|SigHashForkId, 32, 32, value, flag))

	tests = append(tests, NewTestBuilder(
		NewScript().
			PushBytesWithOP(pubkey0.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK INVALID AMOUNT",
		ScriptEnableSigHashForkID,
		false,
		value,
	).PushSig(key0, SigHashAll|SigHashForkId, 32, 32, value+1, flag).ScriptError(ErrInterpreterEvalFalse))

	tests = append(tests, NewTestBuilder(
		NewScript().
			PushBytesWithOP(pubkey0.Bytes()).PushOPCode(OP_CHECKSIG),
		"P2PK INVALID FORKID",
		ScriptVerifyStrictEncoding,
		false,
		value,
	).PushSig(key0, SigHashAll|SigHashForkId, 32, 32, value, flag).ScriptError(ErrInterpreterIllegalForkId))

	// tests[len(tests)-1].Test(t)
	for _, test := range tests {
		test.Test(t)
	}

}
