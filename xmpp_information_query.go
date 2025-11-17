package xmpp

import (
	"fmt"
	"time"
)

const (
	IQTypeGet    = "get"
	IQTypeSet    = "set"
	IQTypeResult = "result"
	IQTypeError  = "error"
)

func (c *Client) Discovery() (string, error) {
	// use UUIDv4 for a pseudo random id.
	reqID := getUUID()
	return c.RawInformationQuery(c.jid, c.domain, reqID, IQTypeGet, XMPPNS_DISCO_ITEMS, "")
}

// Discover information about a node. Empty node queries info about server itself.
func (c *Client) DiscoverNodeInfo(node string) (string, error) {
	query := fmt.Sprintf("<query xmlns='%s' node='%s'/>", XMPPNS_DISCO_INFO, node)
	return c.RawInformation(c.jid, c.domain, getUUID(), IQTypeGet, query)
}

// Discover information about given item from given jid.
func (c *Client) DiscoverInfo(to string) (string, error) {
	query := fmt.Sprintf("<query xmlns='%s'/>", XMPPNS_DISCO_INFO)
	return c.RawInformation(c.jid, to, getUUID(), IQTypeGet, query)
}

// Discover items that the server exposes
func (c *Client) DiscoverServerItems() (string, error) {
	return c.DiscoverEntityItems(c.domain)
}

// Discover items that an entity exposes
func (c *Client) DiscoverEntityItems(jid string) (string, error) {
	query := fmt.Sprintf("<query xmlns='%s'/>", XMPPNS_DISCO_ITEMS)
	return c.RawInformation(c.jid, jid, getUUID(), IQTypeGet, query)
}

// RawInformationQuery sends an information query request to the server.
func (c *Client) RawInformationQuery(from, to, id, iqType, requestNamespace, body string) (string, error) {
	const xmlIQ = "<iq from='%s' to='%s' id='%s' type='%s'><query xmlns='%s'>%s</query></iq>\n"
	_, err := fmt.Fprintf(c.stanzaWriter, xmlIQ, xmlEscape(from), xmlEscape(to), id, iqType, requestNamespace, body)
	return id, err
}

// rawInformation send a IQ request with the payload body to the server
func (c *Client) RawInformation(from, to, id, iqType, body string) (string, error) {
	const xmlIQ = "<iq from='%s' to='%s' id='%s' type='%s'>%s</iq>\n"
	_, err := fmt.Fprintf(c.stanzaWriter, xmlIQ, xmlEscape(from), xmlEscape(to), id, iqType, body)
	return id, err
}

// UrnXMPPTimeResponse implements response to query entity's current time (xep-0202).
func (c *Client) UrnXMPPTimeResponse(v IQ, timezoneOffset string) (string, error) {
	query := fmt.Sprintf(
		"<time xmlns=\"%s\"><tzo>%s</tzo><utc>%s</utc></time>",
		nsTime,
		timezoneOffset,
		time.Now().UTC().Format(time.RFC3339),
	)

	return c.RawInformation(
		v.To,
		v.From,
		v.ID,
		IQTypeResult,
		query,
	)
}
