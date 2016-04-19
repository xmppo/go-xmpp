package xmpp

import (
	"crypto/sha1"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
)

const (
	nsComponentAccept = "jabber:component:accept"
)

// ComponentOptions specify additional options for component client.
type ComponentOptions struct {
	Options
}

// ComponentClient holds XMPP connection opitons for component client.
type ComponentClient struct {
	*Client
}

// jabber:component:accept
type componentMessage struct {
	*clientMessage
	XMLName xml.Name `xml:"jabber:component:accept message"`
}

type componentIQ struct {
	*clientIQ
	XMLName xml.Name `xml:"jabber:component:accept iq"`
}

// NewClient establishes a new component Client connection based on a set of Options.
func (o ComponentOptions) NewClient() (*ComponentClient, error) {
	host := o.Host
	conn, err := connect(host, "", "")
	if err != nil {
		return nil, err
	}

	c := &ComponentClient{new(Client)}
	c.conn = conn

	streamID, err := c.startComponentStream(&o)
	if err != nil {
		return nil, err
	}

	if err := c.startHandshake(streamID, o.Password); err != nil {
		return nil, err
	}

	return c, nil
}

// NewComponentClient creates a new connection to a host given as "hostname" or "hostname:port".
// The connection uses a component protocol, see http://www.xmpp.org/extensions/xep-0114.html
func NewComponentClient(host, cname, secret string, debug bool) (*ComponentClient, error) {
	opts := ComponentOptions{
		Options: Options{
			Host:     host,
			User:     cname,
			Password: secret,
			Debug:    debug,
			Session:  false,
		},
	}
	return opts.NewClient()
}

// startComponentStream will start a new XML decoder for the connection, signal the start of
// a *component* stream to the server and verify that the server has also started the stream.
// If o.Debug is true, startComponentStream will tee decoded XML data to stderr.
func (c *ComponentClient) startComponentStream(o *ComponentOptions) (string, error) {
	if o.Debug {
		c.p = xml.NewDecoder(tee{c.conn, os.Stderr})
	} else {
		c.p = xml.NewDecoder(c.conn)
	}

	_, err := fmt.Fprintf(c.conn, "<?xml version='1.0'?>\n"+
		"<stream:stream to='%s' xmlns='%s'\n"+
		" xmlns:stream='%s' version='1.0'>\n",
		xmlEscape(o.User), nsComponentAccept, nsStream)
	if err != nil {
		return "", err
	}

	// We expect the server to start a <stream>.
	se, err := nextStart(c.p)
	if err != nil {
		return "", err
	}
	if se.Name.Space != nsStream || se.Name.Local != "stream" {
		return "", fmt.Errorf("expected <stream> but got <%v> in %v", se.Name.Local, se.Name.Space)
	}

	// Find the stream id.
	var streamID string
	for _, attr := range se.Attr {
		if attr.Name.Local == "id" {
			streamID = attr.Value
			break
		}
	}
	if streamID == "" {
		return "", errors.New("missing stream id")
	}

	return streamID, nil
}

// startHandshake starts the component handshake on the provided stream and
// checks the authentication results.
func (c *ComponentClient) startHandshake(streamID, secret string) error {
	hash := sha1.New()
	hash.Write([]byte(streamID))
	hash.Write([]byte(secret))

	// Send handshake.
	if _, err := fmt.Fprintf(c.conn, "<handshake>%x</handshake>\n", hash.Sum(nil)); err != nil {
		return err
	}

	// Read handshake reply.
	se, err := nextStart(c.p)
	if err != nil {
		return err
	}

	if se.Name.Local != "handshake" {
		if se.Name.Local == "error" {
			se, err = nextStart(c.p)
			if err != nil {
				return err
			}
			return fmt.Errorf("handshake error: %s", se.Name.Local)
		}
		return fmt.Errorf("handshake error, expected <handshake> but got <%v>", se.Name.Local)
	}

	return nil
}
