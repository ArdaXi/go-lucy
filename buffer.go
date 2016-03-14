package main

import (
	"github.com/zfjagann/golang-ring"
)

type Buffer struct {
	ring *ring.Ring
	len  int
	cap  int
	full bool
	msgs chan string
}

func NewBuffer(capacity int) *Buffer {
	b := &Buffer{
		ring: &ring.Ring{},
	}

	b.ring.SetCapacity(capacity)
	b.cap = capacity
	b.msgs = make(chan string)
	go func() {
		for msg := range b.msgs {
			b.ring.Enqueue(msg)
			if !b.full {
				b.len++
				if b.len == b.cap {
					b.full = true
				}
			}
		}
	}()
	return b
}

func (b *Buffer) Add(i string) {
	b.msgs <- i
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
