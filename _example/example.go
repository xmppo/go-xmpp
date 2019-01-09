package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/kjx98/go-xmpp"
	"log"
	"os"
	"strings"
	"time"
)

var server = flag.String("server", "", "server")
var username = flag.String("username", "", "username")
var password = flag.String("password", "", "password")
var status = flag.String("status", "xa", "status")
var statusMessage = flag.String("status-msg", "I for one welcome our new codebot overlords.", "status message")
var notls = flag.Bool("notls", true, "No TLS")
var debug = flag.Bool("debug", false, "debug output")
var session = flag.Bool("session", false, "use server session")

func serverName(host string) string {
	return strings.Split(host, ":")[0]
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: example [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	if *username == "" || *password == "" {
		if *debug && *username == "" && *password == "" {
			fmt.Fprintf(os.Stderr, "no username or password were given; attempting ANONYMOUS auth\n")
		} else if *username != "" || *password != "" {
			flag.Usage()
		}
	}

	if !*notls {
		xmpp.DefaultConfig = tls.Config{
			ServerName:         serverName(*server),
			InsecureSkipVerify: false,
		}
	}

	var talk *xmpp.Client
	var err error
	options := xmpp.Options{Host: *server,
		User:          *username,
		Password:      *password,
		NoTLS:         *notls,
		Debug:         *debug,
		Session:       *session,
		Status:        *status,
		Resource:      "bot",
		StatusMessage: *statusMessage,
	}

	talk, err = options.NewClient()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			chat, err := talk.Recv()
			if err != nil {
				log.Fatal(err)
			}
			switch v := chat.(type) {
			case xmpp.Chat:
				if v.Type == "roster" {
					fmt.Println("roster", v.Roster)
				} else {
					for _, element := range v.OtherElem {
						if element.XMLName.Space == "jabber:x:conference" {
							// if not join
							talk.JoinMUCNoHistory(v.Remote, "bot")
						}
						// composing, paused, active
						if element.XMLName.Space ==
							"http://jabber.org/protocol/chatstates" &&
							element.XMLName.Local == "composing" {
							fmt.Println(v.Remote, "is composing")
						}
					}
					if strings.TrimSpace(v.Text) != "" {
						fmt.Println(v.Remote, v.Text)
					}
				}
			case xmpp.Presence:
				fmt.Println("Presence:", v.From, v.Show, v.Type)
			case xmpp.Roster, xmpp.Contact:
				// TODO: update local roster
				fmt.Println("Roster/Contact:", v)
			case xmpp.IQ:
				// ping ignore
				if v.Type == "result" && v.ID == "c2s1" {
					fmt.Printf("Got pong from %s to %s\n", v.From, v.To)
				}
			default:
				fmt.Printf("def: %v\n", v)
			}
		}
	}()
	// get roster first
	talk.Roster()
	talk.SendOrg("<presence/>")
	for {
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			continue
		}
		if len(line) >= 4 && line[:4] == "quit" {
			break
		}
		line = strings.TrimRight(line, "\n")

		tokens := strings.SplitN(line, " ", 2)
		if len(tokens) == 2 {
			talk.Send(xmpp.Chat{Remote: tokens[0], Type: "chat", Text: tokens[1]})
		}
	}
	talk.SendOrg("</stream:stream")
	time.Sleep(time.Second * 2)
}
