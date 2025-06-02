package xmpp

import (
	"fmt"
)

// ErrorServiceUnavailable implements error response about feature that is not available. Currently implemented for
// xep-0030.
// QueryXmlns is about incoming xmlns attribute in query tag.
// Node is about incoming node attribute in query tag (looks like it used only in disco#commands).
//
// If queried feature is not here on purpose, standards suggest answer with this stanza.
func (c *Client) ErrorServiceUnavailable(v IQ, queryXmlns, node string) (string, error) {
	query := fmt.Sprintf("<query xmlns=\"%s\" ", queryXmlns)

	if node != "" {
		query += fmt.Sprintf("node=\"%s\" />", node)
	} else {
		query += "/>"
	}

	query += "<error type=\"cancel\">"
	query += "<service-unavailable xmlns=\"urn:ietf:params:xml:ns:xmpp-stanzas\" />"
	query += "</error>"

	return c.RawInformation(
		v.To,
		v.From,
		v.ID,
		IQTypeError,
		query,
	)
}

// ErrorNotImplemented implements error response about feature that is not (yet?) implemented.
// Xmlns is about not implemented feature.
//
// If queried feature is not here because of it under development or for similar reasons, standards suggest answer with
// this stanza.
func (c *Client) ErrorNotImplemented(v IQ, xmlns, feature string) (string, error) {
	query := "<error type=\"cancel\">"
	query += "<feature-not-implemented xmlns=\"urn:ietf:params:xml:ns:xmpp-stanzas\" />"
	query += fmt.Sprintf(
		"<unsupported xmlns=\"%s\" feature=\"%s\" />",
		xmlns,
		feature,
	)
	query += "</error>"

	return c.RawInformation(
		v.To,
		v.From,
		v.ID,
		IQTypeError,
		query,
	)
}
