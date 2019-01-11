package main

import (
	"bufio"
	"crypto/tls"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/kjx98/go-xmpp"
	"log"
	"os"
	"runtime"
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

type rosterItem struct {
	XMLName      xml.Name `xml:"item"`
	Jid          string   `xml:"jid,attr"`
	Name         string   `xml:"name,attr"`
	Subscription string   `xml:"subscription,attr"`
	Group        []string `"xml:"group"`
}

type contactType struct {
	Jid          string
	Name         string
	Subscription string
	Online       bool
}

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
	} else {
		xmpp.DefaultConfig = tls.Config{
			InsecureSkipVerify: true,
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
	loginTime := time.Now()

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
				switch v.Type {
				case "subscribe":
					// Approve all subscription
					fmt.Printf("Presence: %s Approve %s subscription\n",
						v.To, v.From)
					talk.ApproveSubscription(v.From)
					talk.RequestSubscription(v.From)
				case "unsubscribe":
					fmt.Printf("Presence: %s Revoke %s subscription\n",
						v.To, v.From)
					talk.RevokeSubscription(v.From)
				default:
					fmt.Printf("Presence: %s %s Type(%s)\n", v.From, v.Show, v.Type)
				}
			case xmpp.Roster, xmpp.Contact:
				// TODO: update local roster
				fmt.Println("Roster/Contact:", v)
			case xmpp.IQ:
				// ping ignore
				switch v.QueryName.Space {
				case "jabber:iq:version":
					if err := talk.RawVersion(v.To, v.From, v.ID,
						"0.1", runtime.GOOS); err != nil {
						fmt.Println("RawVersion:", err)
					}
					continue
				case "jabber:iq:last":
					tt := time.Now().Sub(loginTime)
					last := int(tt.Seconds())
					if err := talk.RawLast(v.To, v.From, v.ID, last); err != nil {
						fmt.Println("RawLast:", err)
					}
					continue
				case "urn:xmpp:time":
					if err := talk.RawIQtime(v.To, v.From, v.ID); err != nil {
						fmt.Println("RawIQtime:", err)
					}
					continue
				case "jabber:iq:roster":
					var item rosterItem
					if v.Type != "result" && v.Type != "set" {
						// only result and set processed
						fmt.Println("jabber:iq:roster, type:", v.Type)
						continue
					}
					vv := strings.Split(v.Query, "/>")
					for _, ss := range vv {
						if strings.TrimSpace(ss) == "" {
							continue
						}
						ss += "/>"
						if err := xml.Unmarshal([]byte(ss), &item); err != nil {
							fmt.Println("unmarshal roster <query>: ", err)
							continue
						} else {
							if item.Subscription == "remove" {
								continue
							}
							/*
								//may loop whiel presence is unavailable
								if item.Subscription == "from" {
									fmt.Printf("%s Approve %s subscription\n",
										v.To, item.Jid)
									talk.RequestSubscription(item.Jid)
								}
							*/
							fmt.Printf("roster item %s subscription(%s), %v\n",
								item.Jid, item.Subscription, item.Group)
							if v.Type == "set" && item.Subscription == "both" {
								// shall we check presence unavailable
								pr := xmpp.Presence{From: v.To, To: item.Jid,
									Show: "xa"}
								talk.SendPresence(pr)
							}
						}
					}
					continue
				}
				if v.Type == "result" && v.ID == "c2s1" {
					fmt.Printf("Got pong from %s to %s\n", v.From, v.To)
				} else {
					fmt.Printf("Got from %s to %s IQ, tag: (%v), query(%s)\n",
						v.From, v.To, v.QueryName, v.Query)
				}
			default:
				fmt.Printf("def: %v\n", v)
			}
		}
	}()
	// get roster first
	talk.Roster()
	//talk.RevokeSubscription("wkpb@hot-chilli.net")
	//talk.SendOrg("<presence from='wkpb@hot-chilli.net' to='kjx@hot-chilli.net' type='subscribe'/>")
	// test conf
	talk.JoinMUCNoHistory("test@conference.jabb3r.org", "bot")
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
