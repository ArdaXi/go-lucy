package main

import (
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

func RegisterHandlers(bot *ircx.Bot) {
	bot.HandleFunc(irc.PING, PingHandler)
	bot.HandleFunc(irc.RPL_WELCOME, WelcomeHandler)
	bot.HandleFunc(irc.JOIN, JoinHandler)
	bot.HandleFunc(irc.PRIVMSG, MsgHandler)
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
	channel := m.Params[0]
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
	log.Println(nick, channel, m.Trailing)
}

func CommandHandler(nick string, args []string) string {
	return ""
}
