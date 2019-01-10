package xmpp

import (
	"fmt"
	"runtime"
	"time"
)

func (c *Client) SendVersion(id, toServer, fromU string) error {
	_, err := fmt.Fprintf(c.conn, "<iq type='result' from='%s' to='%s'"+
		" id='%s'>", xmlEscape(fromU), xmlEscape(toServer), xmlEscape(id))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(c.conn, "<query xmlns='jabber:iq:version'>"+
		"<name>go-xmpp</name><version>0.1</version><os>%s</os>"+
		"</query>\n</iq>", runtime.GOOS)
	return err
}

func (c *Client) SendIQLast(id, toServer, fromU string) error {
	_, err := fmt.Fprintf(c.conn, "<iq type='result' from='%s' to='%s'"+
		"id='%s' type='result'>\n", xmlEscape(fromU),
		xmlEscape(toServer), xmlEscape(id))
	if err != nil {
		return err
	}
	tt := time.Now().Sub(c.loginTime)
	_, err = fmt.Fprintf(c.conn, "<query xmlns='jabber:iq:last' "+
		"seconds='%d'>Working</query>\n</iq>", int(tt.Seconds()))
	return err
}
