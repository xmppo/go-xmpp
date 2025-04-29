package xmpp

import (
	"encoding/xml"
)

const (
	XMPPNS_DISCO_ITEMS = "http://jabber.org/protocol/disco#items"
	XMPPNS_DISCO_INFO  = "http://jabber.org/protocol/disco#info"
)

type clientDiscoFeature struct {
	XMLName xml.Name `xml:"feature"`
	Var     string   `xml:"var,attr"`
}

type clientDiscoIdentity struct {
	XMLName  xml.Name `xml:"identity"`
	Category string   `xml:"category,attr"`
	Type     string   `xml:"type,attr"`
	Name     string   `xml:"name,attr"`
}

type clientDiscoQuery struct {
	XMLName    xml.Name              `xml:"query"`
	Features   []clientDiscoFeature  `xml:"feature"`
	Identities []clientDiscoIdentity `xml:"identity"`
	X          []DiscoX              `xml:"x"`
}

type clientDiscoItem struct {
	XMLName xml.Name `xml:"item"`
	Jid     string   `xml:"jid,attr"`
	Node    string   `xml:"node,attr"`
	Name    string   `xml:"name,attr"`
}

type clientDiscoItemsQuery struct {
	XMLName xml.Name          `xml:"query"`
	Items   []clientDiscoItem `xml:"item"`
}

type DiscoIdentity struct {
	Category string
	Type     string
	Name     string
}

type DiscoItem struct {
	Jid  string
	Name string
	Node string
}

type DiscoResult struct {
	Features   []string
	Identities []DiscoIdentity
	X          []DiscoX
}

type DiscoX struct {
	XMLName xml.Name `xml:"x"`
	Field   []DiscoXField
}

type DiscoXField struct {
	XMLName xml.Name `xml:"field"`
	Type    string   `xml:"type,attr"`
	Var     string   `xml:"var,attr"`
	Value   []string `xml:"value"`
}

type DiscoItems struct {
	Jid   string
	Items []DiscoItem
}

func clientFeaturesToReturn(features []clientDiscoFeature) []string {
	var ret []string

	for _, feature := range features {
		ret = append(ret, feature.Var)
	}

	return ret
}

func clientIdentitiesToReturn(identities []clientDiscoIdentity) []DiscoIdentity {
	var ret []DiscoIdentity

	for _, id := range identities {
		ret = append(ret, DiscoIdentity{
			Category: id.Category,
			Type:     id.Type,
			Name:     id.Name,
		})
	}

	return ret
}

func clientDiscoItemsToReturn(items []clientDiscoItem) []DiscoItem {
	var ret []DiscoItem
	for _, item := range items {
		ret = append(ret, DiscoItem{
			Jid:  item.Jid,
			Name: item.Name,
			Node: item.Node,
		})
	}

	return ret
}
