// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO(rsc):
//	More precise error handling.
//	Presence functionality.
// TODO(mattn):
//  Add proxy authentication.

// Package xmpp implements a simple Google Talk client
// using the XMPP protocol described in RFC 3920 and RFC 3921.
package xmpp

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"http"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"xml"
)

const (
	nsStream = "http://etherx.jabber.org/streams"
	nsTLS    = "urn:ietf:params:xml:ns:xmpp-tls"
	nsSASL   = "urn:ietf:params:xml:ns:xmpp-sasl"
	nsBind   = "urn:ietf:params:xml:ns:xmpp-bind"
	nsClient = "jabber:client"
)

var DefaultConfig tls.Config

type Client struct {
	tls *tls.Conn // connection to server
	jid string    // Jabber ID for our connection
	p   *xml.Parser
}

// NewClient creates a new connection to a host given as "hostname" or "hostname:port".
// If host is not specified, the  DNS SRV should be used to find the host from the domainpart of the JID.
// Default the port to 5222. 
func NewClient(host, user, passwd string) (*Client, os.Error) {
	addr := host

	if strings.TrimSpace(host) == "" {
		a := strings.Split(user, "@", 2)
		if len(a) == 2 {
			host = a[1]
		}
	}
	a := strings.Split(host, ":", 2)
	if len(a) == 1 {
		host += ":5222"
	}
	proxy := os.Getenv("HTTP_PROXY")
	if proxy == "" {
		proxy = os.Getenv("http_proxy")
	}
	if proxy != "" {
		url, err := http.ParseRequestURL(proxy)
		if err == nil {
			addr = url.Host
		}
	}
	c, err := net.Dial("tcp", "", addr)
	if err != nil {
		return nil, err
	}

	if proxy != "" {
		fmt.Fprintf(c, "CONNECT %s HTTP/1.1\r\n", host)
		fmt.Fprintf(c, "Host: %s\r\n", host)
		fmt.Fprintf(c, "\r\n")
		br := bufio.NewReader(c)
		resp, err := http.ReadResponse(br, "CONNECT")
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			f := strings.Split(resp.Status, " ", 2)
			return nil, os.ErrorString(f[1])
		}
	}

	tlsconn := tls.Client(c, &DefaultConfig)
	if err = tlsconn.Handshake(); err != nil {
		return nil, err
	}

	if strings.LastIndex(host, ":") > 0 {
		host = host[:strings.LastIndex(host, ":")]
	}
	if err = tlsconn.VerifyHostname(host); err != nil {
		return nil, err
	}

	client := new(Client)
	client.tls = tlsconn
	if err := client.init(user, passwd); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

func (c *Client) Close() os.Error {
	return c.tls.Close()
}

func (c *Client) init(user, passwd string) os.Error {
	// For debugging: the following causes the plaintext of the connection to be duplicated to stdout.
	// c.p = xml.NewParser(tee{c.tls, os.Stdout});
	c.p = xml.NewParser(c.tls)

	a := strings.Split(user, "@", 2)
	if len(a) != 2 {
		return os.ErrorString("xmpp: invalid username (want user@domain): " + user)
	}
	user = a[0]
	domain := a[1]

	// Declare intent to be a jabber client.
	fmt.Fprintf(c.tls, "<?xml version='1.0'?>\n"+
		"<stream:stream to='%s' xmlns='%s'\n"+
		" xmlns:stream='%s' version='1.0'>\n",
		xmlEscape(domain), nsClient, nsStream)

	// Server should respond with a stream opening.
	se, err := nextStart(c.p)
	if err != nil {
		return err
	}
	if se.Name.Space != nsStream || se.Name.Local != "stream" {
		return os.ErrorString("xmpp: expected <stream> but got <" + se.Name.Local + "> in " + se.Name.Space)
	}

	// Now we're in the stream and can use Unmarshal.
	// Next message should be <features> to tell us authentication options.
	// See section 4.6 in RFC 3920.
	var f streamFeatures
	if err = c.p.Unmarshal(&f, nil); err != nil {
		return os.ErrorString("unmarshal <features>: " + err.String())
	}
	havePlain := false
	for _, m := range f.Mechanisms.Mechanism {
		if m == "PLAIN" {
			havePlain = true
			break
		}
	}
	if !havePlain {
		return os.ErrorString(fmt.Sprintf("PLAIN authentication is not an option: %v", f.Mechanisms.Mechanism))
	}

	// Plain authentication: send base64-encoded \x00 user \x00 password.
	raw := "\x00" + user + "\x00" + passwd
	enc := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
	base64.StdEncoding.Encode(enc, []byte(raw))
	fmt.Fprintf(c.tls, "<auth xmlns='%s' mechanism='PLAIN'>%s</auth>\n",
		nsSASL, enc)

	// Next message should be either success or failure.
	name, val, err := next(c.p)
	switch v := val.(type) {
	case *saslSuccess:
	case *saslFailure:
		// v.Any is type of sub-element in failure,
		// which gives a description of what failed.
		return os.ErrorString("auth failure: " + v.Any.Local)
	default:
		return os.ErrorString("expected <success> or <failure>, got <" + name.Local + "> in " + name.Space)
	}

	// Now that we're authenticated, we're supposed to start the stream over again.
	// Declare intent to be a jabber client.
	fmt.Fprintf(c.tls, "<stream:stream to='%s' xmlns='%s'\n"+
		" xmlns:stream='%s' version='1.0'>\n",
		xmlEscape(domain), nsClient, nsStream)

	// Here comes another <stream> and <features>.
	se, err = nextStart(c.p)
	if err != nil {
		return err
	}
	if se.Name.Space != nsStream || se.Name.Local != "stream" {
		return os.ErrorString("expected <stream>, got <" + se.Name.Local + "> in " + se.Name.Space)
	}
	if err = c.p.Unmarshal(&f, nil); err != nil {
		// TODO: often stream stop. 
		//return os.ErrorString("unmarshal <features>: " + err.String())
	}

	// Send IQ message asking to bind to the local user name.
	fmt.Fprintf(c.tls, "<iq type='set' id='x'><bind xmlns='%s'/></iq>\n", nsBind)
	var iq clientIQ
	if err = c.p.Unmarshal(&iq, nil); err != nil {
		return os.ErrorString("unmarshal <iq>: " + err.String())
	}
	if iq.Bind == nil {
		return os.ErrorString("<iq> result missing <bind>")
	}
	c.jid = iq.Bind.Jid // our local id

	// We're connected and can now receive and send messages.
	fmt.Fprintf(c.tls, "<presence xml:lang='en'><show>xa</show><status>I for one welcome our new codebot overlords.</status></presence>")
	return nil
}

type Chat struct {
	Remote string
	Type   string
	Text   string
}

// Recv wait next token of chat.
func (c *Client) Recv() (chat Chat, err os.Error) {
	for {
		_, val, err := next(c.p)
		if err != nil {
			return Chat{}, err
		}
		if v, ok := val.(*clientMessage); ok {
			return Chat{v.From, v.Type, v.Body}, nil
		}
	}
	panic("unreachable")
}

// Send sends message text.
func (c *Client) Send(chat Chat) {
	fmt.Fprintf(c.tls, "<message to='%s' from='%s' type='chat' xml:lang='en'>"+
		"<body>%s</body></message>",
		xmlEscape(chat.Remote), xmlEscape(c.jid),
		xmlEscape(chat.Text))
}


// RFC 3920  C.1  Streams name space

type streamFeatures struct {
	XMLName    xml.Name "http://etherx.jabber.org/streams features"
	StartTLS   *tlsStartTLS
	Mechanisms *saslMechanisms
	Bind       *bindBind
	Session    bool
}

type streamError struct {
	XMLName xml.Name "http://etherx.jabber.org/streams error"
	Any     xml.Name
	Text    string
}

// RFC 3920  C.3  TLS name space

type tlsStartTLS struct {
	XMLName  xml.Name ":ietf:params:xml:ns:xmpp-tls starttls"
	Required bool
}

type tlsProceed struct {
	XMLName xml.Name "urn:ietf:params:xml:ns:xmpp-tls proceed"
}

type tlsFailure struct {
	XMLName xml.Name "urn:ietf:params:xml:ns:xmpp-tls failure"
}

// RFC 3920  C.4  SASL name space

type saslMechanisms struct {
	XMLName   xml.Name "urn:ietf:params:xml:ns:xmpp-sasl mechanisms"
	Mechanism []string
}

type saslAuth struct {
	XMLName   xml.Name "urn:ietf:params:xml:ns:xmpp-sasl auth"
	Mechanism string   "attr"
}

type saslChallenge string

type saslResponse string

type saslAbort struct {
	XMLName xml.Name "urn:ietf:params:xml:ns:xmpp-sasl abort"
}

type saslSuccess struct {
	XMLName xml.Name "urn:ietf:params:xml:ns:xmpp-sasl success"
}

type saslFailure struct {
	XMLName xml.Name "urn:ietf:params:xml:ns:xmpp-sasl failure"
	Any     xml.Name
}

// RFC 3920  C.5  Resource binding name space

type bindBind struct {
	XMLName  xml.Name "urn:ietf:params:xml:ns:xmpp-bind bind"
	Resource string
	Jid      string
}

// RFC 3921  B.1  jabber:client

type clientMessage struct {
	XMLName xml.Name "jabber:client message"
	From    string   "attr"
	Id      string   "attr"
	To      string   "attr"
	Type    string   "attr" // chat, error, groupchat, headline, or normal

	// These should technically be []clientText,
	// but string is much more convenient.
	Subject string
	Body    string
	Thread  string
}

type clientText struct {
	Lang string "attr"
	Body string "chardata"
}

type clientPresence struct {
	XMLName xml.Name "jabber:client presence"
	From    string   "attr"
	Id      string   "attr"
	To      string   "attr"
	Type    string   "attr" // error, probe, subscribe, subscribed, unavailable, unsubscribe, unsubscribed
	Lang    string   "attr"

	Show     string // away, chat, dnd, xa
	Status   string // sb []clientText
	Priority string
	Error    *clientError
}

type clientIQ struct { // info/query
	XMLName xml.Name "jabber:client iq"
	From    string   "attr"
	Id      string   "attr"
	To      string   "attr"
	Type    string   "attr" // error, get, result, set
	Error   *clientError
	Bind    *bindBind
}

type clientError struct {
	XMLName xml.Name "jabber:client error"
	Code    string   "attr"
	Type    string   "attr"
	Any     xml.Name
	Text    string
}

// Scan XML token stream to find next StartElement.
func nextStart(p *xml.Parser) (xml.StartElement, os.Error) {
	for {
		t, err := p.Token()
		if err != nil {
			log.Fatal("token", err)
		}
		switch t := t.(type) {
		case xml.StartElement:
			return t, nil
		}
	}
	panic("unreachable")
}

// Prototypical nil pointers for specific XML element names.
var proto = map[string]interface{}{
	nsStream + " features": (*streamFeatures)(nil),
	nsStream + " error":    (*streamError)(nil),

	nsTLS + " starttls": (*tlsStartTLS)(nil),
	nsTLS + " proceed":  (*tlsProceed)(nil),
	nsTLS + " failure":  (*tlsFailure)(nil),

	nsSASL + " mechanisms": (*saslMechanisms)(nil),
	nsSASL + " challenge":  (*saslChallenge)(nil),
	nsSASL + " response":   (*saslResponse)(nil),
	nsSASL + " abort":      (*saslAbort)(nil),
	nsSASL + " success":    (*saslSuccess)(nil),
	nsSASL + " failure":    (*saslFailure)(nil),

	nsBind + " bind": (*bindBind)(nil),

	nsClient + " message":  (*clientMessage)(nil),
	nsClient + " presence": (*clientPresence)(nil),
	nsClient + " iq":       (*clientIQ)(nil),
	nsClient + " error":    (*clientError)(nil),
}

// Scan XML token stream for next element and save into val.
// If val == nil, allocate new element based on proto map.
// Either way, return val.
func next(p *xml.Parser) (xml.Name, interface{}, os.Error) {
	// Read start element to find out what type we want.
	se, err := nextStart(p)
	if err != nil {
		return xml.Name{}, nil, err
	}
	v, ok := proto[se.Name.Space+" "+se.Name.Local]
	if !ok {
		return xml.Name{}, nil, os.ErrorString("unexpected XMPP message " +
			se.Name.Space + " <" + se.Name.Local + "/>")
	}

	// The map lookup got us a pointer.
	// Put it in an interface and allocate one.
	pv := reflect.NewValue(v).(*reflect.PtrValue)
	zv := reflect.MakeZero(pv.Type().(*reflect.PtrType).Elem())
	pv.PointTo(zv)

	// Unmarshal into that storage.
	if err = p.Unmarshal(pv.Interface(), &se); err != nil {
		return xml.Name{}, nil, err
	}
	return se.Name, pv.Interface(), err
}

var xmlSpecial = map[byte]string{
	'<':  "&lt;",
	'>':  "&gt;",
	'"':  "&quot;",
	'\'': "&apos;",
	'&':  "&amp;",
}

func xmlEscape(s string) string {
	var b bytes.Buffer
	for i := 0; i < len(s); i++ {
		c := s[i]
		if s, ok := xmlSpecial[c]; ok {
			b.WriteString(s)
		} else {
			b.WriteByte(c)
		}
	}
	return b.String()
}

type tee struct {
	r io.Reader
	w io.Writer
}

func (t tee) Read(p []byte) (n int, err os.Error) {
	n, err = t.r.Read(p)
	if n > 0 {
		t.w.Write(p[0:n])
	}
	return
}
