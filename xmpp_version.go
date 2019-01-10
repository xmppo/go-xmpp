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
	ss := fmt.Sprintf("<iq type='result' from='%s' to='%s'"+
		" id='%s'>\n", xmlEscape(fromU), xmlEscape(toServer), xmlEscape(id))
	tt := time.Now().Sub(c.loginTime)
	ss += fmt.Sprintf("<query xmlns='jabber:iq:last' "+
		"seconds='%d'>Working</query>\n</iq>", int(tt.Seconds()))
	_, err := fmt.Fprint(c.conn, ss)
	return err
}

func (c *Client) SendIQtime(id, toServer, fromU string) error {
	ss := fmt.Sprintf("<iq type='result' from='%s' to='%s'"+
		" id='%s'>\n", xmlEscape(fromU), xmlEscape(toServer), xmlEscape(id))
	tt := time.Now()
	zoneN, _ := tt.Zone()
	ss += fmt.Sprintf("<time xmlns='urn:xmpp:time'>\n<tzo>%s</tzo>"+
		"<utc>%s</utc></time>\n</iq>", zoneN,
		tt.UTC().Format("2006-01-02T15:03:04Z"))
	_, err := fmt.Fprint(c.conn, ss)
	return err
}
