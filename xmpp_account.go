package xmpp

import (
	"encoding/xml"
	"fmt"
	"strings"
)

const (
	nsRegister = "jabber:iq:register"
	nsSearch   = "jabber:iq:search"
)

type clientSearchAccountItem struct {
	Jid   string `xml:"jid,attr"`
	First string `xml:"first"`
	Last  string `xml:"last"`
	Nick  string `xml:"nick"`
	Email string `xml:"email"`
}

type clientSearchAccountQuery struct {
	XMLName xml.Name                  `xml:"query"`
	Xmlns   string                    `xml:"xmlns,attr"`
	Items   []clientSearchAccountItem `xml:"item"`
}

// SearchAccountResultItem represent search account item result
type SearchAccountResultItem struct {
	Jid       string
	FirstName string
	LastName  string
	NickName  string
	Email     string
}

// SearchAccountResult represent search account item result
type SearchAccountResult struct {
	Jid      string
	Accounts []SearchAccountResultItem
}

func clientSearchAccountItemToReturn(accounts []clientSearchAccountItem) []SearchAccountResultItem {
	var ret []SearchAccountResultItem
	for _, account := range accounts {
		ret = append(ret, SearchAccountResultItem{
			Jid:       account.Jid,
			FirstName: account.First,
			LastName:  account.Last,
			NickName:  account.Nick,
			Email:     account.Email,
		})
	}

	return ret
}

func buildAccountAttr(username string, password string, attributes map[string]string) string {
	var attrBuilder strings.Builder
	attrBuilder.WriteString(fmt.Sprintf("<username>%s</username>", username))
	attrBuilder.WriteString(fmt.Sprintf("<password>%s</password>", password))

	if attributes != nil {
		for k, v := range attributes {
			attrBuilder.WriteString(fmt.Sprintf("<%s>%s</%s>", k, v, k))
		}
	}

	return attrBuilder.String()
}

// CreateAccount Creates a new account using the specified username, password and extra attributes
func (c *Client) CreateAccount(username string, password string, attributes map[string]string) error {
	attrBuilder := buildAccountAttr(username, password, attributes)
	const xmlIQ = "<iq type='set' id='createAccount1'><query xmlns='%s'>%s</query></iq>"
	_, err := fmt.Fprintf(c.stanzaWriter, xmlIQ, nsRegister, attrBuilder)
	return err
}

// ChangePassword enables a user to change his or her password with a server or service
func (c *Client) ChangePassword(username string, password string) error {
	attrBuilder := buildAccountAttr(username, password, nil)
	const xmlIQ = "<iq type='set' to='%s' id='changePassword1'><query xmlns='%s'>%s</query></iq>"
	_, err := fmt.Fprintf(c.stanzaWriter, xmlIQ, c.domain, nsRegister, attrBuilder)
	return err
}

// RemoveAccount cancel a registration with a host by sending a <remove/> element in an IQ set
func (c *Client) RemoveAccount(username string) error {
	from := c.jid
	const xmlIQ = "<iq type='set' from='%s' id='removeAccount1'><query xmlns='%s'><remove/></query></iq>"
	_, err := fmt.Fprintf(c.stanzaWriter, xmlIQ, from, nsRegister)
	return err
}

// SearchAccount search information repositories on the Jabber network
// searchServiceName is the Search Service Properties Name from your server
func (c *Client) SearchAccount(searchServiceName, username, fieldName, fieldValue string) error {
	from := c.jid
	searchService := fmt.Sprintf("%s.%s", searchServiceName, c.domain)
	searchQuery := fmt.Sprintf("<%s>%s</%s>", fieldName, fieldValue, fieldName)
	const xmlIQ = "<iq type='set' from='%s' to='%s' id='searchAccount1' xml:lang='en'><query xmlns='%s'>%s</query></iq>\n"
	_, err := fmt.Fprintf(c.stanzaWriter, xmlIQ, from, searchService, nsSearch, searchQuery)
	return err
}
