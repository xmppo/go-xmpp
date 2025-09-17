package xmpp

import (
	"fmt"
)

func (c *Client) ApproveSubscription(jid string) {
	fmt.Fprintf(c.stanzaWriter, "<presence to='%s' type='subscribed'/>\n",
		xmlEscape(jid))
}

func (c *Client) RevokeSubscription(jid string) {
	fmt.Fprintf(c.stanzaWriter, "<presence to='%s' type='unsubscribed'/>\n",
		xmlEscape(jid))
}

// DEPRECATED: Use RevertSubscription instead.
func (c *Client) RetrieveSubscription(jid string) {
	c.RevertSubscription(jid)
}

func (c *Client) RevertSubscription(jid string) {
	fmt.Fprintf(c.conn, "<presence to='%s' type='unsubscribe'/>\n",
		xmlEscape(jid))
}

func (c *Client) RequestSubscription(jid string) {
	fmt.Fprintf(c.stanzaWriter, "<presence to='%s' type='subscribe'/>\n",
		xmlEscape(jid))
}
