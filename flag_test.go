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
	newsig = append(newsig, 0x20)
	newsig = append(newsig, byte(len(r)))
	newsig = append(newsig, r...)
	newsig = append(newsig, 0x20)
	newsig = append(newsig, byte(len(s)))
	newsig = append(newsig, s...)

	return newsig
}

func (tb *TestBuilder) Num(n int) *TestBuilder {
	tb.DoPush()
	tb.spendTx.Inputs[0].ScriptSig = NewScriptFromBytes(
		tb.spendTx.Inputs[0].ScriptSig,
	).PushNumber(Number(n)).Bytes()
	return tb
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

	left := tb.push[:pos]
	right := tb.push[pos+len(datain):]

	tb.push = append(left, dataout...)
	tb.push = append(tb.push, right...)

	return tb
}

func (tb *TestBuilder) DoPushBytes(data []byte) {
	tb.DoPush()
	tb.push = data
	tb.havePush = true
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
	flags Flag,
) *TestBuilder {
	ts := NewTransactionSigner(tb.spendTx, 0, amount)
	hash := ts.signatureHashOriginal(tb.script, sigHash)
	sig := tb.DoSign(key, hash, 32, 32)
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
	pubkey0c, _ := key0c.GetPubkey()

	key1 := bcrypto.NewKey(key1bytes[:], false)
	key1c := bcrypto.NewKey(key1bytes[:], true)
	pubkey1, _ := key1.GetPubkey()
	pubkey1c, _ := key1c.GetPubkey()

	key2 := bcrypto.NewKey(key2bytes[:], false)
	key2c := bcrypto.NewKey(key2bytes[:], true)
	// pubkey2, _ := key1.GetPubkey()
	pubkey2c, _ := key2c.GetPubkey()

	var flag Flag
	flag.Enable(ScriptEnableSigHashForkID)

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
	).PushSig(key1, SigHashAll.Enable(SigHashAnyoneCanPay), 32, 32, 0, flag))

	tests = append(tests, NewTestBuilder(
		NewScript().PushBytesWithOP(pubkey1.Bytes()).
			PushOPCode(OP_CHECKSIG),
		"P2PK anyonecanpy marked with normal hashtyp",
		0,
		false,
		0,
	).PushSig(key1, SigHashAll.Enable(SigHashAnyoneCanPay), 32, 32, 0, flag).
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
		ScriptError(ErrTransactionSignerEmptySignature))

	// for _, test := range tests {
	tests[len(tests)-1].Test(t)
	// }

}
