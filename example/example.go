package main

import (
	".."
	"bufio"
	"flag"
	"fmt"
	"github.com/mattn/go-iconv"
	"log"
	"os"
	"strings"
)

var server = flag.String("server", "talk.google.com:443", "server")
var username = flag.String("username", "", "username")
var password = flag.String("password", "", "password")
var insecure = flag.Bool("insecure", false, "ignore certificates")
var trace = flag.Bool("trace", false, "prints raw data to stderr")
var room = flag.String("room", "", "joins the specified room")

func fromUTF8(s string) string {
	ic, err := iconv.Open("char", "UTF-8")
	if err != nil {
		return s
	}
	defer ic.Close()
	ret, _ := ic.Conv(s)
	return ret
}

func toUTF8(s string) string {
	ic, err := iconv.Open("UTF-8", "char")
	if err != nil {
		return s
	}
	defer ic.Close()
	ret, _ := ic.Conv(s)
	return ret
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: example [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	if *username == "" || *password == "" {
		flag.Usage()
	}
	if *insecure {
		xmpp.DefaultConfig.TLS.InsecureSkipVerify = true
	}
	if *trace {
		xmpp.DefaultConfig.Debug.R = os.Stderr
	}
	talk, err := xmpp.NewClient(*server, *username, *password)
	if err != nil {
		log.Fatal(err)
	}
	if *room != "" {
		var p xmpp.Presence
		nick := strings.SplitN(*username, "@", 2)[0]
		p.From = talk.Jid()
		p.To = *room + "/" + nick
		err = talk.Send(&p)
		if err != nil {
			log.Fatal(err)
		}
	}
	go func() {
		for {
			evt, err := talk.Recv()
			if err != nil {
				log.Fatal(err)
			}
			if chat, ok := evt.(*xmpp.Message); ok {
				fmt.Println(chat.From, fromUTF8(chat.Body))
			}
		}
	}()
	for {
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			continue
		}
		line = strings.TrimRight(line, "\n")

		tokens := strings.SplitN(line, " ", 2)
		if len(tokens) == 2 {
			var msg xmpp.Message
			msg.To, msg.Body = tokens[0], toUTF8(tokens[1])
			if msg.To == *room {
				msg.Type = "groupchat"
			} else {
				msg.Type = "chat"
			}
			talk.Send(&msg)
		}
	}
}
