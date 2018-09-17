package bscript

import (
	"testing"
)

func TestStack(t *testing.T) {
	stack := NewStack()

	stack.Push(Number(1).Bytes())
	stack.Push(Number(2).Bytes())
	stack.Push(Number(3).Bytes())

	d0, _ := stack.Peek(0)
	d1, _ := stack.Peek(-3)
	d2, _ := stack.Peek(-2)
	d3, _ := stack.Peek(-1)

	n0, _ := NewNumberFromBytes(d0, false, 4)
	n1, _ := NewNumberFromBytes(d1, false, 4)
	n2, _ := NewNumberFromBytes(d2, false, 4)
	n3, _ := NewNumberFromBytes(d3, false, 4)
	if n0 != 1 {
		t.Error("expected 1")
	}
	if n1 != 1 {
		t.Error("expected 3")
	}
	if n2 != 2 {
		t.Error("expected 2")
	}
	if n3 != 3 {
		t.Error("expected 1")
	}
}
