package xmpp

import (
	"fmt"
)

func (c *Client) PingC2S(jid, server string) error {
	if jid == "" {
		jid = c.jid
	}
	if server == "" {
		server = c.domain
	}
	// Reset ticker for periodic pings if configured.
	if c.periodicPings {
		c.periodicPingTicker.Reset(c.periodicPingPeriod)
	}
	_, err := fmt.Fprintf(c.stanzaWriter, "<iq from='%s' to='%s' id='%s' type='get'>"+
		"<ping xmlns='urn:xmpp:ping'/>"+
		"</iq>\n",
		xmlEscape(jid), xmlEscape(server), getUUIDv4())
	return err
}

func (c *Client) PingS2S(fromServer, toServer string) error {
	// Reset ticker for periodic pings if configured.
	if c.periodicPings {
		c.periodicPingTicker.Reset(c.periodicPingPeriod)
	}
	_, err := fmt.Fprintf(c.stanzaWriter, "<iq from='%s' to='%s' id='%s' type='get'>"+
		"<ping xmlns='urn:xmpp:ping'/>"+
		"</iq>\n",
		xmlEscape(fromServer), xmlEscape(toServer), getUUIDv4())
	return err
}

func (c *Client) SendResultPing(id, toServer string) error {
	// Reset ticker for periodic pings if configured.
	if c.periodicPings {
		c.periodicPingTicker.Reset(c.periodicPingPeriod)
	}
	_, err := fmt.Fprintf(c.stanzaWriter, "<iq type='result' to='%s' id='%s'/>\n",
		xmlEscape(toServer), xmlEscape(id))
	return err
}

func (c *Client) sendPeriodicPings() {
	for t := range c.periodicPingTicker.C {
		fmt.Println("Tick at", t)
		_ = c.PingC2S(c.jid, c.domain)
	}
}
