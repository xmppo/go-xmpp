package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mattn/go-iconv"
	"github.com/mattn/go-xmpp"
	"log"
	"os"
	"strings"
)

var server = flag.String("server", "talk.google.com:443", "server")
var username = flag.String("username", "", "username")
var password = flag.String("password", "", "password")
var insecure = flag.Bool("insecure", false, "ignore certificates")
var trace = flag.Bool("trace", false, "prints raw data to stderr")

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

	go func() {
		for {
			chat, err := talk.Recv()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(chat.Remote, fromUTF8(chat.Text))
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
			talk.Send(xmpp.Chat{Remote: tokens[0], Type: "chat", Text: toUTF8(tokens[1])})
		}
	}
}
