package bscript

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func scriptTestName(test []interface{}) (string, error) {
	// Account for any optional leading witness data.
	var witnessOffset int
	if _, ok := test[0].([]interface{}); ok {
		witnessOffset++
	}

	// In addition to the optional leading witness data, the test must
	// consist of at least a signature script, public key script, flags,
	// and expected error.  Finally, it may optionally contain a comment.
	if len(test) < witnessOffset+4 || len(test) > witnessOffset+5 {
		return "", fmt.Errorf("invalid test length %d", len(test))
	}

	// Use the comment for the test name if one is specified, otherwise,
	// construct the name based on the signature script, public key script,
	// and flags.
	var name string
	if len(test) == witnessOffset+5 {
		name = fmt.Sprintf("test (%s)", test[witnessOffset+4])
	} else {
		name = fmt.Sprintf("test ([%s, %s, %s])", test[witnessOffset],
			test[witnessOffset+1], test[witnessOffset+2])
	}
	return name, nil
}

// parseWitnessStack parses a json array of witness items encoded as hex into a
// slice of witness elements.
func parseWitnessStack(elements []interface{}) ([][]byte, error) {
	witness := make([][]byte, len(elements))
	for i, e := range elements {
		witElement, err := hex.DecodeString(e.(string))
		if err != nil {
			return nil, err
		}

		witness[i] = witElement
	}

	return witness, nil
}

func TestInterpreter(t *testing.T) {
	f, err := ioutil.ReadFile("testdata/script_tests.json")
	if err != nil {
		t.Fatal(err)
	}

	var tests [][]interface{}
	err = json.Unmarshal(f, &tests)
	if err != nil {
		t.Fatalf("TestScripts couldn't Unmarshal: %v", err)
	}

	for i, test := range tests {
		if len(test) == 1 {
			continue
		}
		name, err := scriptTestName(test)
		if err != nil {
			t.Errorf("TestScripts: invalid test #%d: %v", i, err)
			continue
		}
		// scriptSigString, scriptPubKeyString, scriptFlagString, scriptErrorString := test[0], test[1], test[2], test[3]
		var (
			witness ScriptWitness
			amount  int64
		// inputAmt btcutil.Amount
		)

		// "0",
		// "IF 0 ELSE 1 ELSE 0 ENDIF",
		// "P2SH,STRICTENC",
		// "OK",
		// "Multiple ELSE's are valid and executed inverts on each ELSE encountered"

		// When the first field of the test data is a slice it contains
		// witness data and everything else is offset by 1 as a result.
		witnessOffset := 0
		if witnessData, ok := test[0].([]interface{}); ok {
			witnessOffset++

			// If this is a witness test, then the final element
			// within the slice is the input amount, so we ignore
			// all but the last element in order to parse the
			// witness stack.
			strWitnesses := witnessData[:len(witnessData)-1]
			witness, err = parseWitnessStack(strWitnesses)
			if err != nil {
				t.Errorf("%s: can't parse witness; %v", name, err)
				continue
			}

			amount = int64(witnessData[len(witnessData)-1].(float64) * 1e8)
			// fmt.Println(witness, "abcd", amount)
			_ = amount
			_ = witness
		}

		// Extract and parse the signature script from the test fields.
		scriptSigStr, ok := test[witnessOffset].(string)
		if !ok {
			t.Errorf("%s: signature script is not a string", name)
			continue
		}

		// Extract and parse the public key script from the test fields.
		scriptPubKeyStr, ok := test[witnessOffset+1].(string)
		if !ok {
			t.Errorf("%s: public key script is not a string", name)
			continue
		}

		// Extract and parse the script flags from the test fields.
		flagsStr, ok := test[witnessOffset+2].(string)
		if !ok {
			t.Errorf("%s: flags field is not a string", name)
			continue
		}

		scriptErrStr, ok := test[witnessOffset+3].(string)
		if !ok {
			continue
		}

		sigScript, err := parseScript(scriptSigStr)
		if err != nil {
			continue
		}

		pubkeyScript, err := parseScript(scriptPubKeyStr)
		if err != nil {
			continue
		}

		flag := NewFlagFromString(flagsStr)

		err = VerifyScript(
			sigScript,
			pubkeyScript,
			nil,
			flag,
			NewNoopChecker(),
			SignatureVersionBase,
		)

		if scriptErrStr != "OK" {
			if err != nil {
				continue
			}
			fmt.Println(sigScript, "qq", pubkeyScript)
			fmt.Println(scriptSigStr, scriptPubKeyStr, "shit")
			fmt.Println(err, scriptErrStr, "fuck")
			t.Fatal("fuck")
		}
	}
}

var scriptFlagMap = map[string]Flag{
	"NONE":                                  ScriptVerifyNone,
	"P2SH":                                  ScriptVerifyP2SH,
	"STRICTENC":                             ScriptVerifyStrictEncoding,
	"DERSIG":                                ScriptVerifyDERSignatures,
	"LOW_S":                                 ScriptVerifyLowS,
	"SIGPUSHONLY":                           ScriptVerifySigPushOnly,
	"MINIMALDATA":                           ScriptVerifyMinimalData,
	"NULLDUMMY":                             ScriptVerifyNullFail,
	"DISCOURAGE_UPGRADABLE_NOPS":            ScriptDiscourageUpgradableNops,
	"CLEANSTACK":                            ScriptVerifyCleanStack,
	"MINIMALIF":                             ScriptVerifyMinimalIf,
	"NULLFAIL":                              ScriptVerifyNullFail,
	"CHECKLOCKTIMEVERIFY":                   ScriptVerifyCheckLockTimeVerify,
	"CHECKSEQUENCEVERIFY":                   ScriptVerifyCheckSequenceVerify,
	"DISCOURAGE_UPGRADABLE_WITNESS_PROGRAM": ScriptVerifyDiscourageUpgradeableWitnessProgram,
	"COMPRESSED_PUBKEYTYPE":                 ScriptVerifyCompressedPubkeyType,
	"SIGHASH_FORKID":                        ScriptEnableSigHashForkID,
	"REPLAY_PROTECTION":                     ScriptEnableReplayProtection,
	"MONOLITH_OPCODES":                      ScriptEnableMonolithOpcodes,
}

func parseScript(str string) (*Script, error) {
	b := make([]string, len(str))
	for _, s := range strings.Split(str, " ") {
		_, err := NewOPCodeFromString("OP_" + s)
		if err != nil {
			b = append(b, s)
			continue
		}
		b = append(b, "OP_"+s)
	}

	fmt.Println(b)

	return NewScriptFromString(strings.Join(b, " "))
}

func NewFlagFromString(f string) Flag {
	flag := Flag(0)
	for _, s := range strings.Split(f, ",") {
		if len(s) == 0 {
			continue
		}

		v, ok := scriptFlagMap[s]
		if !ok {
			continue
		}
		flag.Enable(v)
	}

	return flag
}
