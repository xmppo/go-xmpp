package xmpp

import (
	"fmt"
)

const (
	IQTypeGet    = "get"
	IQTypeSet    = "set"
	IQTypeResult = "result"
	IQTypeError  = "error"
)

func (c *Client) Discovery() (string, error) {
	// use UUIDv4 for a pseudo random id.
	reqID := getUUIDv4()
	// Reset ticker for periodic pings if configured.
	if c.periodicPings {
		c.periodicPingTicker.Reset(c.periodicPingPeriod)
	}
	return c.RawInformationQuery(c.jid, c.domain, reqID, IQTypeGet, XMPPNS_DISCO_ITEMS, "")
}

// Discover information about a node. Empty node queries info about server itself.
func (c *Client) DiscoverNodeInfo(node string) (string, error) {
	query := fmt.Sprintf("<query xmlns='%s' node='%s'/>", XMPPNS_DISCO_INFO, node)
	// Reset ticker for periodic pings if configured.
	if c.periodicPings {
		c.periodicPingTicker.Reset(c.periodicPingPeriod)
	}
	return c.RawInformation(c.jid, c.domain, getUUIDv4(), IQTypeGet, query)
}

// Discover information about given item from given jid.
func (c *Client) DiscoverInfo(to string) (string, error) {
	query := fmt.Sprintf("<query xmlns='%s'/>", XMPPNS_DISCO_INFO)
	// Reset ticker for periodic pings if configured.
	if c.periodicPings {
		c.periodicPingTicker.Reset(c.periodicPingPeriod)
	}
	return c.RawInformation(c.jid, to, getUUIDv4(), IQTypeGet, query)
}

// Discover items that the server exposes
func (c *Client) DiscoverServerItems() (string, error) {
	// Reset ticker for periodic pings if configured.
	if c.periodicPings {
		c.periodicPingTicker.Reset(c.periodicPingPeriod)
	}
	return c.DiscoverEntityItems(c.domain)
}

// Discover items that an entity exposes
func (c *Client) DiscoverEntityItems(jid string) (string, error) {
	query := fmt.Sprintf("<query xmlns='%s'/>", XMPPNS_DISCO_ITEMS)
	// Reset ticker for periodic pings if configured.
	if c.periodicPings {
		c.periodicPingTicker.Reset(c.periodicPingPeriod)
	}
	return c.RawInformation(c.jid, jid, getUUIDv4(), IQTypeGet, query)
}

// RawInformationQuery sends an information query request to the server.
func (c *Client) RawInformationQuery(from, to, id, iqType, requestNamespace, body string) (string, error) {
	const xmlIQ = "<iq from='%s' to='%s' id='%s' type='%s'><query xmlns='%s'>%s</query></iq>\n"
	// Reset ticker for periodic pings if configured.
	if c.periodicPings {
		c.periodicPingTicker.Reset(c.periodicPingPeriod)
	}
	_, err := fmt.Fprintf(c.stanzaWriter, xmlIQ, xmlEscape(from), xmlEscape(to), id, iqType, requestNamespace, body)
	return id, err
}

// rawInformation send a IQ request with the payload body to the server
func (c *Client) RawInformation(from, to, id, iqType, body string) (string, error) {
	const xmlIQ = "<iq from='%s' to='%s' id='%s' type='%s'>%s</iq>\n"
	// Reset ticker for periodic pings if configured.
	if c.periodicPings {
		c.periodicPingTicker.Reset(c.periodicPingPeriod)
	}
	_, err := fmt.Fprintf(c.stanzaWriter, xmlIQ, xmlEscape(from), xmlEscape(to), id, iqType, body)
	return id, err
}
