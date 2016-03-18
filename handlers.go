package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

var sanitizer = regexp.MustCompile(`[^\w ']+`)

type Lucy struct {
	buffer *Buffer
}

func AsyncHandleFunc(l *Lucy, b *ircx.Bot, cmd string,
	handlerFunc func(l *Lucy, s ircx.Sender, m *irc.Message)) {
	handler := AsyncHandler{
		handler: handlerFunc,
		lucy:    l,
	}
	b.Handle(cmd, handler)
}

type AsyncHandler struct {
	handler func(l *Lucy, s ircx.Sender, m *irc.Message)
	lucy    *Lucy
}

func (f AsyncHandler) Handle(s ircx.Sender, m *irc.Message) {
	go f.handler(f.lucy, s, m)
}

func RegisterHandlers(bot *ircx.Bot) {
	lucy := &Lucy{buffer: NewBuffer(10)}

	AsyncHandleFunc(lucy, bot, irc.PING, PingHandler)
	AsyncHandleFunc(lucy, bot, irc.RPL_WELCOME, WelcomeHandler)
	AsyncHandleFunc(lucy, bot, irc.ERR_NICKNAMEINUSE, NickCollisionHandler)
	AsyncHandleFunc(lucy, bot, irc.JOIN, JoinHandler)
	AsyncHandleFunc(lucy, bot, irc.PRIVMSG, MsgHandler)
}

func PingHandler(l *Lucy, s ircx.Sender, m *irc.Message) {
	log.Println("ping", m)
	s.Send(&irc.Message{
		Command:  irc.PONG,
		Params:   m.Params,
		Trailing: m.Trailing,
	})
}

func WelcomeHandler(l *Lucy, s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  []string{*channels},
	})
}

func NickCollisionHandler(l *Lucy, s ircx.Sender, m *irc.Message) {
	nick := m.Params[1]
	*name = nick + "_"
	log.Println(nick, "already in use, changing nick to", *name)
	s.Send(&irc.Message{
		Command: irc.NICK,
		Params:  []string{*name},
	})
}

func JoinHandler(l *Lucy, s ircx.Sender, m *irc.Message) {
	nick := m.Prefix.Name
	if nick == *name {
		log.Println("Joined", m.Trailing)
	} else {
		log.Println(nick, "joined", m.Trailing)
	}
}

func MsgHandler(l *Lucy, s ircx.Sender, m *irc.Message) {
	nick := m.Prefix.Name
	//channel := m.Params[0]
	msg := m.Trailing
	args := strings.Fields(msg)
	l.buffer.Add(sanitizer.ReplaceAllString(msg, " "))
	if strings.TrimRight(args[0], ",:") == *name {
		res := CommandHandler(nick, args[1:])
		if res != "" {
			s.Send(&irc.Message{
				Command:  irc.PRIVMSG,
				Params:   m.Params,
				Trailing: res,
			})
		}
	}
}

func CommandHandler(nick string, args []string) string {
	return fmt.Sprintf("<%v> %v", nick, strings.Join(args, " "))
}
