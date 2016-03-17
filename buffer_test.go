package main

import (
	"testing"
	"time"
)

func TestBufferEmpty(t *testing.T) {
	b := NewBuffer(2)

	time.Sleep(1)

	if res := b.List(); res != nil {
		t.Errorf("List(), expected: %v, got: %v", nil, res)
	}

	if res := b.Full(); res {
		t.Errorf("Full(), expected: %v, got: %v", false, res)
	}
}

func TestBufferHalfFull(t *testing.T) {
	b := NewBuffer(2)
	b.Add("1")

	if res := b.List(); res != nil {
		t.Errorf("List(), expected: %v, got: %v", nil, res)
	}

	if res := b.Full(); res {
		t.Errorf("Full(), expected: %v, got: %v", false, res)
	}
}

func TestBufferFull(t *testing.T) {
	b := NewBuffer(2)
	b.Add("1")
	b.Add("2")

	res := b.List()
	if len(res) != 2 {
		t.Errorf("List() length, expected: %v, got: %v", 2, len(res))
	}
	if res[0] != "1" {
		t.Errorf("res[0] mismatch, expected: %v, got %v", "1", res[0])
	}
	if res[1] != "2" {
		t.Errorf("res[1] mismatch, expected: %v, got %v", "2", res[1])
	}

	if res := b.Full(); !res {
		t.Errorf("Full(), expected: %v, got: %v", true, res)
	}
}

func TestBufferOverflow(t *testing.T) {
	b := NewBuffer(2)
	b.Add("1")
	b.Add("2")
	b.Add("3")

	res := b.List()
	if len(res) != 2 {
		t.Errorf("List() length, expected: %v, got: %v", 2, len(res))
	}
	if res[0] != "2" {
		t.Errorf("res[0] mismatch, expected: %v, got %v", "1", res[0])
	}
	if res[1] != "3" {
		t.Errorf("res[1] mismatch, expected: %v, got %v", "2", res[1])
	}

	if res := b.Full(); !res {
		t.Errorf("Full(), expected: %v, got: %v", true, res)
	}
}
