package main

import (
	"testing"

	"github.com/sorcix/irc"
)

type testBuffer struct{}

func (b testBuffer) Add(s string) {}

func (b testBuffer) List() []interface{} {
	return []interface{}{}
}

func (b testBuffer) Full() bool {
	return false
}

type testSender struct {
	lastMessage *irc.Message
}

func (s *testSender) Send(m *irc.Message) error {
	s.lastMessage = m
	return nil
}

func TestPingHandler(t *testing.T) {
	l := &Lucy{buffer: &testBuffer{}}
	s := &testSender{}
	m := &irc.Message{
		Command:  irc.PING,
		Params:   []string{},
		Trailing: ":server",
	}

	PingHandler(l, s, m)
	resp := *s.lastMessage
	if resp.Command != irc.PONG {
		t.Errorf("Command mismatch, expected: %v, got: %v", irc.PONG, resp.Command)
	}
	if len(resp.Params) != len(m.Params) {
		t.Errorf("Params mismatch, expected: %v, got: %v", m.Params, resp.Params)
	}
	if resp.Trailing != m.Trailing {
		t.Errorf("Trailing mismatch, expected: %v, got: %v", m.Trailing, resp.Trailing)
	}
}

func TestWelcomeHandler(t *testing.T) {
	*channels = "#chan"
	l := &Lucy{buffer: &testBuffer{}}
	s := &testSender{}
	m := &irc.Message{
		Command: irc.RPL_WELCOME,
		Params:  []string{},
	}

	WelcomeHandler(l, s, m)
	resp := *s.lastMessage
	if resp.Command != irc.JOIN {
		t.Errorf("Command mismatch, expected: %v, got: %v", irc.JOIN, resp.Command)
	}
	if len(resp.Params) != 1 || resp.Params[0] != "#chan" {
		t.Errorf("Params mismatch, expected: %v, got: %v", []string{"#chan"}, resp.Params)
	}
}
