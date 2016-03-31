package xmpp

import (
	"fmt"
)

const xmlIqGet = "<iq from='%s' to='%s' id='%d' type='get'><query xmlns='http://jabber.org/protocol/disco#items'/></iq>"

func (c *Client) Discovery() {
	cookie := getCookie()
	fmt.Fprintf(c.conn, xmlIqGet, xmlEscape(c.jid), xmlEscape(c.domain), cookie)
}
