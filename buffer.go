package main

import (
	"github.com/zfjagann/golang-ring"
)

type Buffer struct {
	ring *ring.Ring
	len  int
	cap  int
	full bool
}

func NewBuffer(capacity int) *Buffer {
	buffer := &Buffer{
		ring: &ring.Ring{},
	}

	buffer.ring.SetCapacity(capacity)
	buffer.cap = capacity
	return buffer
}

func (b *Buffer) Add(i interface{}) {
	b.ring.Enqueue(i)
	if !b.full {
		b.len++
		if b.len == b.cap {
			b.full = true
		}
	}
}

func (b *Buffer) List() []interface{} {
	if !b.full {
		return nil
	}
	return b.ring.Values()
}

func (b *Buffer) Full() bool {
	return b.full
}
