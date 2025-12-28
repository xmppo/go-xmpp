package xmpp

import (
	"fmt"
	"time"
)

// Discovery discovers items information according https://xmpp.org/extensions/xep-0030.html#items (Discovering the
// Items Associated with a Jabber Entity).
func (c *Client) Discovery() (string, error) {
	return c.RawInformationQuery(c.jid, c.domain, getUUID(), IQTypeGet, XMPPNS_DISCO_ITEMS, "")
}

// DiscoverNodeInfo discovers information about a node. Empty node queries info about server itself.
// Discovery query performed according to https://xmpp.org/extensions/xep-0030.html#info (Discovering Information About
// a Jabber Entity).
func (c *Client) DiscoverNodeInfo(node string) (string, error) {
	return c.RawInformation(
		c.jid,
		c.domain,
		getUUID(),
		IQTypeGet,
		fmt.Sprintf("<query xmlns=%q node=%q/>", XMPPNS_DISCO_INFO, node),
	)
}

// DiscoverInfo discovers information about given item from given jid.
// Discovery query performed according to https://xmpp.org/extensions/xep-0030.html#info (Discovering Information About
// a Jabber Entity).
// The only difference between DiscoverInfo() and DiscoverNodeInfo() is that DiscoverInfo() does not supply From field,
// which is useful in very limited amount use cases.
func (c *Client) DiscoverInfo(to string) (string, error) {
	query := fmt.Sprintf("<query xmlns=%q/>", XMPPNS_DISCO_INFO)

	return c.RawInformation(c.jid, to, getUUID(), IQTypeGet, query)
}

// DiscoverServerItems discovers items that the server exposes. It is actually thin wrapper for DiscoverEntityItems().
func (c *Client) DiscoverServerItems() (string, error) {
	return c.DiscoverEntityItems(c.domain)
}

// DiscoverEntityItems discovers items that an entity exposes.
func (c *Client) DiscoverEntityItems(jid string) (string, error) {
	query := fmt.Sprintf("<query xmlns=%q/>", XMPPNS_DISCO_ITEMS)

	return c.RawInformation(c.jid, jid, getUUID(), IQTypeGet, query)
}

// RawInformationQuery sends an information query request to the server.
func (c *Client) RawInformationQuery(from, to, id, iqType, requestNamespace, body string) (string, error) {
	_, err := fmt.Fprintf(
		c.stanzaWriter,
		"<iq from=%q to=%q id=%q type=%q><query xmlns=%q>%s</query></iq>\n",
		xmlEscape(from),
		xmlEscape(to),
		id,
		iqType,
		requestNamespace,
		body,
	)

	return id, err
}

// RawInformation send a IQ request with the payload body to the server.
func (c *Client) RawInformation(from, to, id, iqType, body string) (string, error) {
	_, err := fmt.Fprintf(
		c.stanzaWriter,
		"<iq from=%q to=%q id=%q type=%q>%s</iq>\n",
		xmlEscape(from),
		xmlEscape(to),
		id,
		iqType,
		body,
	)

	return id, err
}

// UrnXMPPTimeResponse implements response to query entity's current time accodring to
// https://xmpp.org/extensions/xep-0202.html#example-2 (A Response to the Query).
func (c *Client) UrnXMPPTimeResponse(v IQ, timezoneOffset string) (string, error) {
	query := fmt.Sprintf(
		"<time xmlns=%q><tzo>%s</tzo><utc>%s</utc></time>",
		XMPPNS_TIME,
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

// IqVersionResponse responding with software version, according to example described in
// https://xmpp.org/extensions/xep-0092.html#example-2 (Receiving a Reply Regarding Software Version).
func (c *Client) IqVersionResponse(v IQ, name string, version string, os string) (string, error) {
	if name == "" {
		name = "go-xmpp"
		version = Version
	}

	if version == "" {
		version = "undefined"
	}

	query := fmt.Sprintf("<query xmlns=%q>", XMPPNS_IQ_VERSION)

	query += fmt.Sprintf("<name>%s</name>", name)
	query += fmt.Sprintf("<version>%s</version>", version)

	if os != "" {
		query += fmt.Sprintf("<os>%s</os>", os)
	}

	query += "</query>"

	return c.RawInformation(
		v.To,
		v.From,
		v.ID,
		IQTypeResult,
		query,
	)
}
