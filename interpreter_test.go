package bscript

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func getTests(testfile string) ([][]interface{}, error) {
	file, err := ioutil.ReadFile(testfile)
	if err != nil {
		return nil, err
	}

	var tests [][]interface{}
	if err := json.Unmarshal(file, &tests); err != nil {
		return nil, err
	}
	return tests, nil
}

func run(code string) (*Stack, error) {
	script, err := NewScriptFromString(code)
	if err != nil {
		return nil, err
	}

	interpreter := NewInterpreter()
	err = interpreter.Eval(script, ScriptSkipDisabledOPCode)
	if err != nil {
		return nil, err
	}

	return interpreter.GetDStack(), nil
}

func TestArithmetic(t *testing.T) {
	tests, err := getTests("testdata/arithmetic.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		name := test[0]
		code := test[1].(string)
		expect := test[2].(float64)
		if err != nil {
			t.Fatal(err)
		}
		s, err := run(code)
		if err != nil {
			t.Fatal(err)
		}

		d, err := s.Peek(-1)
		if err != nil {
			t.Fatal(err)
		}

		r, err := NewNumberFromBytes(d.Bytes(), false, NumberDefaultElementSize)
		if err != nil {
			t.Fatal(err)
		}

		if float64(r) != expect {
			t.Errorf("%s: expect %f got %d with \"%s\"", name, expect, r, code)
		}
	}
}

func TestBitwise(t *testing.T) {
	tests, err := getTests("testdata/bitwise.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		name := test[0]
		code := test[1].(string)
		expect := test[2].(float64)
		if err != nil {
			t.Fatal(err)
		}
		s, err := run(code)
		if err != nil {
			t.Fatal(err)
		}

		d, err := s.Peek(-1)
		if err != nil {
			t.Fatal(err)
		}

		r, err := NewNumberFromBytes(d.Bytes(), false, NumberDefaultElementSize)
		if err != nil {
			t.Fatal(err)
		}

		if float64(r) != expect {
			t.Errorf("%s: expect %f got %d with \"%s\"", name, expect, r, code)
		}
	}
}

func TestPush(t *testing.T) {
	tests, err := getTests("testdata/push.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		name := test[0]
		code := test[1].(string)
		expect := test[2].(float64)
		if err != nil {
			t.Fatal(err)
		}
		s, err := run(code)
		if err != nil {
			t.Fatal(err)
		}

		d, err := s.Peek(-1)
		if err != nil {
			t.Fatal(err)
		}

		r, err := NewNumberFromBytes(d.Bytes(), false, NumberDefaultElementSize)
		if err != nil {
			t.Fatal(err)
		}

		if float64(r) != expect {
			t.Errorf("%s: expect %f got %d with \"%s\"", name, expect, r, code)
		}
	}
}

func TestFlowControl(t *testing.T) {
	tests, err := getTests("testdata/flowcontrol.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		name := test[0]
		code := test[1].(string)
		expect := test[2].(float64)
		if err != nil {
			t.Fatal(err)
		}
		s, err := run(code)
		if err != nil {
			t.Fatal(err)
		}

		d, err := s.Peek(-1)
		if err != nil {
			t.Fatal(err)
		}

		r, err := NewNumberFromBytes(d.Bytes(), false, NumberDefaultElementSize)
		if err != nil {
			t.Fatal(err)
		}

		if float64(r) != expect {
			t.Errorf("%s: expect %f got %d with \"%s\"", name, expect, r, code)
		}
	}
}

func TestPushData(t *testing.T) {
	tests, err := getTests("testdata/pushdata.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		name := test[0]
		code := test[1].(string)
		expect := test[2].(float64)
		if err != nil {
			t.Fatal(err)
		}
		s, err := run(code)
		if err != nil {
			t.Fatal(err)
		}

		d, err := s.Peek(-1)
		if err != nil {
			t.Fatal(err)
		}

		r, err := NewNumberFromBytes(d.Bytes(), false, NumberDefaultElementSize)
		if err != nil {
			t.Fatal(err)
		}

		if float64(r) != expect {
			t.Errorf("%s: expect %f got %d with \"%s\"", name, expect, r, code)
		}
	}
}

func TestStackOPS(t *testing.T) {
	tests, err := getTests("testdata/stackops.json")
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		name := test[0]
		code := test[1].(string)
		expect := test[2].([]interface{})
		if err != nil {
			t.Fatal(err)
		}
		s, err := run(code)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < s.Depth(); i++ {
			d, err := s.Pop()
			if err != nil {
				t.Fatal(err)
				break
			}
			n, err := NewNumberFromBytes(d.Bytes(), false, NumberDefaultElementSize)
			if err != nil {
				t.Fatal(err)
			}

			if expect[len(expect)-i-1].(float64) != float64(n) {
				t.Fatalf("%s: expect %f got %d with \"%s\"", name, expect[len(expect)-i-1], d, code)
			}
		}
	}
}

func TestStackOP(t *testing.T) {
	tests, err := getTests("testdata/stack.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		name := test[0]
		code := test[1].(string)
		expect := test[2].([]interface{})
		if err != nil {
			t.Fatal(err)
		}
		s, err := run(code)
		if err != nil {
			t.Fatal(err)
			break
		}

		for i := 0; i < s.Depth(); i++ {
			d, err := s.Pop()
			if err != nil {
				t.Fatal(err)
				break
			}
			n, err := NewNumberFromBytes(d.Bytes(), false, NumberDefaultElementSize)
			if err != nil {
				t.Fatal(err)
			}

			if expect[len(expect)-i-1].(float64) != float64(n) {
				t.Fatalf("%s: expect %f got %d with \"%s\"", name, expect[len(expect)-i-1], d, code)
			}
		}
	}
}

func TestConstants(t *testing.T) {
	tests, err := getTests("testdata/constants.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		name := test[0]
		code := test[1].(string)
		expect := test[2].(float64)
		if err != nil {
			t.Fatal(err)
		}
		s, err := run(code)
		if err != nil {
			t.Fatal(err)
		}

		d, err := s.Peek(-1)
		if err != nil {
			t.Fatal(err)
		}

		r, err := NewNumberFromBytes(d.Bytes(), false, NumberDefaultElementSize)
		if err != nil {
			t.Fatal(err)
		}

		if float64(r) != expect {
			t.Errorf("%s: expect %f got %d with \"%s\"", name, expect, r, code)
		}
	}
}

func TestCrypto(t *testing.T) {
	tests, err := getTests("testdata/crypto.json")
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		name := test[0]
		code := test[1].(string)
		tmp := test[2].(string)
		expect, err := hex.DecodeString(tmp[2:])
		if err != nil {
			t.Fatal(err)
		}
		s, err := run(code)
		if err != nil {
			t.Fatal(err)
			break
		}
		d, err := s.Peek(-1)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(expect, d.Bytes()) {
			t.Fatalf("%s: expect %#v got %#v with \"%s\"", name, expect, d, code)
		}
		break
	}
}

func TestInterpreterOPNumberic(t *testing.T) {
	type Test struct {
		code string
		rv   int64
	}
	tests := []Test{
		{
			code: "1 2 OP_ADD",
			rv:   3,
		},
		{
			code: "0x01 0x02 OP_ADD",
			rv:   3,
		},
		{
			code: "0x01 2 OP_ADD",
			rv:   3,
		},
		{
			code: "1 2 OP_ADD 4 OP_ADD",
			rv:   7,
		},
		{
			code: "1 2 OP_SUB",
			rv:   -1,
		},
		{
			code: "0x01 0x02 OP_SUB",
			rv:   -1,
		},
		{
			code: "0x01 2 OP_SUB",
			rv:   -1,
		},
		{
			code: "1 2 OP_MUL",
			rv:   2,
		},
		{
			code: "1 2 OP_MUL 3 OP_MUL 4 OP_MUL 5 OP_MUL 10 OP_MUL",
			rv:   1200,
		},
	}

	for _, test := range tests {
		script, err := NewScriptFromString(test.code)
		if err != nil {
			t.Fatal(err)
		}

		interpreter := NewInterpreter()
		err = interpreter.Eval(script, ScriptSkipDisabledOPCode)
		if err != nil {
			t.Fatal(err)
		}

		d, err := interpreter.dstack.Peek(-1)
		if err != nil {
			t.Fatal(err)
		}

		n, err := NewNumberFromBytes(d.Bytes(), false, NumberDefaultElementSize)
		if err != nil {
			t.Fatal(err)
		}

		if int64(n) != test.rv {
			t.Fatalf("expect %d", test.rv)
		}
	}
}
