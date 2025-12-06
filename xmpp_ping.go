package xmpp

import (
	"fmt"
	"time"
)

func (c *Client) PingC2S(jid, server string) error {
	if jid == "" {
		jid = c.jid
	}

	if server == "" {
		server = c.domain
	}

	_, err := fmt.Fprintf(c.stanzaWriter, "<iq from='%s' to='%s' id='%s' type='get'>"+
		"<ping xmlns='urn:xmpp:ping'/>"+
		"</iq>\n",
		xmlEscape(jid), xmlEscape(server), getUUID())

	return err
}

func (c *Client) PingS2S(fromServer, toServer string) error {
	_, err := fmt.Fprintf(c.stanzaWriter, "<iq from='%s' to='%s' id='%s' type='get'>"+
		"<ping xmlns='urn:xmpp:ping'/>"+
		"</iq>\n",
		xmlEscape(fromServer), xmlEscape(toServer), getUUID())

	return err
}

func (c *Client) SendResultPing(id, toServer string) error {
	_, err := fmt.Fprintf(c.stanzaWriter, "<iq type='result' to='%s' id='%s'/>\n",
		xmlEscape(toServer), xmlEscape(id))

	return err
}

func (c *Client) sendPeriodicPings() {
	for range c.periodicPingTicker.C {
		// Reset ticker for periodic pings if configured.
		if c.periodicPings {
			c.periodicPingTicker.Reset(c.periodicPingPeriod)
		}

		c.periodicPingID = getUUID()
		c.periodicPingReply = false

		_, err := fmt.Fprintf(c.stanzaWriter, "<iq from='%s' to='%s' id='%s' type='get'>"+
			"<ping xmlns='urn:xmpp:ping'/></iq>\n",
			xmlEscape(c.jid), xmlEscape(c.domain), c.periodicPingID)

		if err != nil {
			c.Close()
		}

		time.Sleep(c.periodicPingTimeout)

		if !c.periodicPingReply {
			c.shutdown = true
			fmt.Fprintf(c.stanzaWriter, "</stream:stream>\n")
			c.conn.Close()
		}
	}
}
