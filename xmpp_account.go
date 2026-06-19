package xmpp

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

const (
	nsRegister     = "jabber:iq:register"
	nsSearch       = "jabber:iq:search"
	nsLastActivity = "jabber:iq:last"
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

type lastActivity struct {
	XMLName xml.Name `xml:"query"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Seconds string   `xml:"seconds,attr"`
}

// LastActivityResult represent last activity response
// when current entity session is authorized to view the user's presence information
type LastActivityResult struct {
	From              string
	Text              string
	LastActiveSeconds int
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

func handleLastActivityResult(from string, lastActivity lastActivity) (LastActivityResult, error) {
	lastActiveSeconds, err := strconv.Atoi(lastActivity.Seconds)
	if err != nil {
		return LastActivityResult{}, err
	}

	return LastActivityResult{
		From:              from,
		Text:              lastActivity.Text,
		LastActiveSeconds: lastActiveSeconds,
	}, nil
}

func buildAccountAttr(username string, password string, attributes map[string]string) string {
	var attrBuilder strings.Builder
	attrBuilder.WriteString(fmt.Sprintf("<username>%s</username>", xmlEscape(username)))
	attrBuilder.WriteString(fmt.Sprintf("<password>%s</password>", xmlEscape(password)))

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

// ChangePassword enables a user to change his or her password with a server or service.
// With the user's username parameter which will change the password and new password
func (c *Client) ChangePassword(username string, newPassword string) error {
	attrBuilder := buildAccountAttr(username, newPassword, nil)
	const xmlIQ = "<iq type='set' to='%s' id='changePassword1'><query xmlns='%s'>%s</query></iq>"
	_, err := fmt.Fprintf(c.stanzaWriter, xmlIQ, xmlEscape(c.domain), nsRegister, attrBuilder)
	return err
}

// RemoveAccount cancel or delete the current session registration
// with a host by sending a <remove/> element in an IQ set.
func (c *Client) RemoveAccount() error {
	from := c.jid
	const xmlIQ = "<iq type='set' from='%s' id='removeAccount1'><query xmlns='%s'><remove/></query></iq>"
	_, err := fmt.Fprintf(c.stanzaWriter, xmlIQ, xmlEscape(from), nsRegister)
	return err
}

// SearchAccount search information repositories on the Jabber network.
// searchServiceName is the Search Service Properties Name from your server.
func (c *Client) SearchAccount(searchServiceName, fieldName, fieldValue string) error {
	from := c.jid
	searchService := fmt.Sprintf("%s.%s", searchServiceName, c.domain)
	searchQuery := fmt.Sprintf("<%s>%s</%s>", fieldName, xmlEscape(fieldValue), fieldName)
	const xmlIQ = "<iq type='set' from='%s' to='%s' id='searchAccount1' xml:lang='en'><query xmlns='%s'>%s</query></iq>\n"
	_, err := fmt.Fprintf(c.stanzaWriter, xmlIQ, xmlEscape(from), searchService, nsSearch, searchQuery)
	return err
}

// RequestLastActivity request last activity information regarding another entity
func (c *Client) RequestLastActivity(to string) error {
	targetEntity := strings.SplitN(to, "@", 2)
	if len(targetEntity) < 2 {
		to = fmt.Sprintf("%s@%s", to, c.domain)
	}

	from := c.jid
	const xmlIQ = "<iq from='%s' id='requestLastActivity1' to='%s' type='get'><query xmlns='%s'/></iq>"
	_, err := fmt.Fprintf(c.stanzaWriter, xmlIQ, xmlEscape(from), xmlEscape(to), nsLastActivity)
	return err
}
