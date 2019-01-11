package xmpp

import (
	"fmt"
	"time"
)

func (c *Client) RawVersion(from, to, id, version, osName string) error {
	body := "<name>go-xmpp</name><version>" + version + "</version><os>" +
		osName + "</os>"
	_, err := c.RawInformationQuery(from, to, id, "result", "jabber:iq:version",
		body)
	return err
}

func (c *Client) RawLast(from, to, id string, last int) error {
	body := fmt.Sprintf("<query xmlns='jabber:iq:last' "+
		"seconds='%d'>Working</query>", last)
	_, err := c.RawInformation(from, to, id, "result", body)
	return err
}

func (c *Client) RawIQtime(from, to, id string) error {
	tt := time.Now()
	zone, _ := tt.Zone()
	body := fmt.Sprintf("<time xmlns='urn:xmpp:time'>\n<tzo>%s</tzo><utc>%s"+
		"</utc></time>", zone, tt.UTC().Format("2006-01-02T15:03:04Z"))
	_, err := c.RawInformation(from, to, id, "result", body)
	return err
}
