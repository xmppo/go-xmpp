package xmpp

import (
	"encoding/xml"
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

// Oob is an out-of-band url/description, used in file uploads.
// See https://xmpp.org/extensions/xep-0066.html
type Oob struct {
	XMLName xml.Name `xml:"x,xmlns:jabber:x:oob"`
	Url     string   `xml:"url"`
	Desc    string   `xml:"desc"`
}
