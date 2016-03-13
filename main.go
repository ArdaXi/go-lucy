package main

import (
	"flag"
	"log"

	"github.com/nickvanw/ircx"
)

var (
	name     = flag.String("name", "Lucy", "Nick to use on IRC")
	server   = flag.String("server", "chat.freenode.org:6667", "Host to connect to")
	channels = flag.String("chan", "#lucy-test", "Channels to join")
)

func init() {
	flag.Parse()
}

func main() {
	cfg := ircx.Config{
		User: *name,
	}
	bot := ircx.New(*server, *name, cfg)
	if err := bot.Connect(); err != nil {
		log.Panicln("Unable to connect to server,", err)
	}
	log.Println("Connected to server:", *server)
	RegisterHandlers(bot)
	bot.HandleLoop()
	log.Println("Exiting..")
}
