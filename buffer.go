package main

import (
	"github.com/zfjagann/golang-ring"
)

type Buffer interface {
	Add(string)
	List() []interface{}
	Full() bool
}

type buffer struct {
	ring *ring.Ring
	len  int
	cap  int
	full bool
	msgs chan string
	done chan struct{}
}

func NewBuffer(capacity int) *buffer {
	b := &buffer{
		ring: &ring.Ring{},
	}

	b.ring.SetCapacity(capacity)
	b.cap = capacity
	b.msgs = make(chan string)
	b.done = make(chan struct{})
	go func() {
		for msg := range b.msgs {
			b.ring.Enqueue(msg)
			if !b.full {
				b.len++
				if b.len == b.cap {
					b.full = true
				}
			}
			b.done <- struct{}{}
		}
	}()
	return b
}

func (b *buffer) Add(i string) {
	b.msgs <- i
	<-b.done
}

func (b *buffer) List() []interface{} {
	if !b.full {
		return nil
	}
	return b.ring.Values()
}

func (b *buffer) Full() bool {
	return b.full
}
