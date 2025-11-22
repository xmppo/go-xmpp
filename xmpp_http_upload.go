package xmpp

import (
	"encoding/xml"
)

const (
	XMPPNS_HTTP_UPLOAD = "urn:xmpp:http:upload:0"
)

type Slot struct {
	// TODO: Maybe this doesn't belong here
	ID      string
	XMLName xml.Name `xml:"slot"`
	Put     Put
	Get     Get
}

type Put struct {
	XMLName xml.Name `xml:"put"`
	Url     string   `xml:"url,attr"`
	Headers []Header `xml:"header"`
}

type Get struct {
	XMLName xml.Name `xml:"get"`
	Url     string   `xml:"url,attr"`
}

type Header struct {
	XMLName xml.Name `xml:"header"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",innerxml"`
}
