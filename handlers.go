package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

var (
	sanitizer = regexp.MustCompile(`[^\w ']+`)
	buffer    = NewBuffer(10)
)

func AsyncHandleFunc(b *ircx.Bot, cmd string, handler func(s ircx.Sender, m *irc.Message)) {
	b.Handle(cmd, AsyncHandlerFunc(handler))
}

type AsyncHandlerFunc func(s ircx.Sender, m *irc.Message)

func (f AsyncHandlerFunc) Handle(s ircx.Sender, m *irc.Message) {
	go f(s, m)
}

func RegisterHandlers(bot *ircx.Bot) {
	AsyncHandleFunc(bot, irc.PING, PingHandler)
	AsyncHandleFunc(bot, irc.RPL_WELCOME, WelcomeHandler)
	AsyncHandleFunc(bot, irc.ERR_NICKNAMEINUSE, NickCollisionHandler)
	AsyncHandleFunc(bot, irc.JOIN, JoinHandler)
	AsyncHandleFunc(bot, irc.PRIVMSG, MsgHandler)
}

func PingHandler(s ircx.Sender, m *irc.Message) {
	log.Println("ping", m)
	s.Send(&irc.Message{
		Command:  irc.PONG,
		Params:   m.Params,
		Trailing: m.Trailing,
	})
}

func WelcomeHandler(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  []string{*channels},
	})
}

func NickCollisionHandler(s ircx.Sender, m *irc.Message) {
	nick := m.Params[1]
	*name = nick + "_"
	log.Println(nick, "already in use, changing nick to", *name)
	s.Send(&irc.Message{
		Command: irc.NICK,
		Params:  []string{*name},
	})
}

func JoinHandler(s ircx.Sender, m *irc.Message) {
	nick := m.Prefix.Name
	if nick == *name {
		log.Println("Joined", m.Trailing)
	} else {
		log.Println(nick, "joined", m.Trailing)
	}
}

func MsgHandler(s ircx.Sender, m *irc.Message) {
	nick := m.Prefix.Name
	//channel := m.Params[0]
	msg := m.Trailing
	args := strings.Fields(msg)
	buffer.Add(sanitizer.ReplaceAllString(msg, " "))
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
