package xmpp

import (
	"bytes"
	"encoding/xml"
	"io"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"
)

type localAddr struct{}

func (a *localAddr) Network() string {
	return "tcp"
}

func (addr *localAddr) String() string {
	return "localhost:5222"
}

type testConn struct {
	*bytes.Buffer
}

func tConnect(s string) net.Conn {
	var conn testConn
	conn.Buffer = bytes.NewBufferString(s)
	return &conn
}

func (*testConn) Close() error {
	return nil
}

func (*testConn) LocalAddr() net.Addr {
	return &localAddr{}
}

func (*testConn) RemoteAddr() net.Addr {
	return &localAddr{}
}

func (*testConn) SetDeadline(time.Time) error {
	return nil
}

func (*testConn) SetReadDeadline(time.Time) error {
	return nil
}

func (*testConn) SetWriteDeadline(time.Time) error {
	return nil
}

var text = strings.TrimSpace(`
<message xmlns="jabber:client" id="3" type="error" to="123456789@gcm.googleapis.com/ABC">
	<gcm xmlns="google:mobile:data">
		{"random": "&lt;text&gt;"}
	</gcm>
	<error code="400" type="modify">
		<bad-request xmlns="urn:ietf:params:xml:ns:xmpp-stanzas"/>
		<text xmlns="urn:ietf:params:xml:ns:xmpp-stanzas">
			InvalidJson: JSON_PARSING_ERROR : Missing Required Field: message_id\n
		</text>
	</error>
</message>
`)

func TestStanzaError(t *testing.T) {
	var c Client
	c.conn = tConnect(text)
	c.p = xml.NewDecoder(c.conn)
	v, err := c.Recv()
	if err != nil {
		t.Fatalf("Recv() = %v", err)
	}

	chat := Chat{
		Type: "error",
		Other: []string{
			"\n\t\t{\"random\": \"<text>\"}\n\t",
			"\n\t\t\n\t\t\n\t",
		},
		OtherElem: []XMLElement{
			{
				XMLName:  xml.Name{Space: "google:mobile:data", Local: "gcm"},
				Attr:     []xml.Attr{{Name: xml.Name{Space: "", Local: "xmlns"}, Value: "google:mobile:data"}},
				InnerXML: "\n\t\t{\"random\": \"&lt;text&gt;\"}\n\t",
			},
			{
				XMLName: xml.Name{Space: "jabber:client", Local: "error"},
				Attr:    []xml.Attr{{Name: xml.Name{Space: "", Local: "code"}, Value: "400"}, {Name: xml.Name{Space: "", Local: "type"}, Value: "modify"}},
				InnerXML: `
		<bad-request xmlns="urn:ietf:params:xml:ns:xmpp-stanzas"/>
		<text xmlns="urn:ietf:params:xml:ns:xmpp-stanzas">
			InvalidJson: JSON_PARSING_ERROR : Missing Required Field: message_id\n
		</text>
	`,
			},
		},
	}

	if !reflect.DeepEqual(v, chat) {
		t.Errorf("Recv() = %#v; want %#v", v, chat)
	}
}

func TestEOFError(t *testing.T) {
	var c Client
	c.conn = tConnect("")
	c.p = xml.NewDecoder(c.conn)
	_, err := c.Recv()

	if err != io.EOF {
		t.Errorf("Recv() did not return io.EOF on end of input stream")
	}
}

var emptyPubSub = strings.TrimSpace(`
<iq xmlns="jabber:client" type='result' from='juliet@capulet.lit' id='items3'>
  <pubsub xmlns='http://jabber.org/protocol/pubsub'>
    <items node='urn:xmpp:avatar:data'></items>
  </pubsub>
</iq>
`)

func TestEmptyPubsub(t *testing.T) {
	var c Client
	c.itemsIDs = append(c.itemsIDs, "items3")
	c.conn = tConnect(emptyPubSub)
	c.p = xml.NewDecoder(c.conn)
	m, err := c.Recv()

	switch m.(type) {
	case AvatarData:
		if err == nil {
			t.Errorf("Expected an error to be returned")
		}

	default:
		t.Errorf("Recv() = %v", m)
		t.Errorf("Expected a return value of AvatarData")
	}
}

// https://xmpp.org/extensions/xep-0363.html#example-6
var exampleSlot = strings.TrimSpace(`
<iq from='upload.montague.tld'
    xmlns="jabber:client"
    id='abcdef'
    to='romeo@montague.tld/garden'
    type='result'>
  <slot xmlns='urn:xmpp:http:upload:0'>
    <put url='https://upload.montague.tld/4a771ac1-f0b2-4a4a-9700-f2a26fa2bb67/tr%C3%A8s%20cool.jpg'>
      <header name='Authorization'>Basic Base64String==</header>
      <header name='Cookie'>foo=bar; user=romeo</header>
    </put>
    <get url='https://download.montague.tld/4a771ac1-f0b2-4a4a-9700-f2a26fa2bb67/tr%C3%A8s%20cool.jpg' />
  </slot>
</iq>
`)

func TestUploadSlot(t *testing.T) {
	var c Client
	// c.itemsIDs = append(c.itemsIDs, "step_03")
	c.conn = tConnect(exampleSlot)
	c.p = xml.NewDecoder(c.conn)
	m, err := c.Recv()
	if err != nil {
		panic(err)
	}

	t.Logf("Recv() = %v", m)

	switch m.(type) {
	case Slot:
		v, _ := m.(Slot)

		if v.ID != "abcdef" {
			t.Errorf("Invalid ID: %s", v.ID)
		}

		if v.Put.Url != "https://upload.montague.tld/4a771ac1-f0b2-4a4a-9700-f2a26fa2bb67/tr%C3%A8s%20cool.jpg" {
			t.Errorf("Invalid PUT URL: %s", v.Put.Url)
		}
		if v.Get.Url != "https://download.montague.tld/4a771ac1-f0b2-4a4a-9700-f2a26fa2bb67/tr%C3%A8s%20cool.jpg" {
			t.Errorf("Invalid GET URL: %s", v.Get.Url)
		}

		foundAuthorization := false
		foundCookie := false

		for _, header := range v.Put.Headers {
			if header.Name == "Authorization" && header.Value == "Basic Base64String==" {
				foundAuthorization = true
				continue
			}

			if header.Name == "Cookie" && header.Value == "foo=bar; user=romeo" {
				foundCookie = true
				continue
			}

			t.Errorf("Unknown header: %s: %s", header.Name, header.Value)
		}

		if !foundAuthorization {
			t.Errorf("Authorization header not found")
		}
		if !foundCookie {
			t.Errorf("Cookie header not found")
		}

	default:
		t.Errorf("Recv() = %V", m)
		t.Errorf("Expected a return value of Slot")
	}
}

// https://xmpp.org/extensions/xep-0066.html#example-5
var exampleOOB = strings.TrimSpace(`
<message from='stpeter@jabber.org/work'
         to='MaineBoy@jabber.org/home'
		 xmlns='jabber:client'>
  <body>Yeah, but do you have a license to Jabber?</body>
  <x xmlns='jabber:x:oob'>
    <url>http://www.jabber.org/images/psa-license.jpg</url>
  </x>
</message>
`)

func TestChatOOB(t *testing.T) {
	var c Client
	c.conn = tConnect(exampleOOB)
	c.p = xml.NewDecoder(c.conn)
	m, err := c.Recv()
	if err != nil {
		panic(err)
	}

	t.Logf("Recv() = %v", m)

	switch m.(type) {
	case Chat:
		v, _ := m.(Chat)

		if v.Oob.Url != "http://www.jabber.org/images/psa-license.jpg" {
			t.Errorf("Wrong URL, found: `%s`", v.Oob.Url)
		}

		if v.Oob.Desc != "" {
			t.Errorf("Should not find Desc: `%s`", v.Oob.Desc)
		}

	default:
		t.Errorf("Recv() = %v", m)
		t.Errorf("Expected a return value of AvatarData")
	}
}

var exampleNoOOB = strings.TrimSpace(`
<message from='stpeter@jabber.org/work'
         to='MaineBoy@jabber.org/home'
		 xmlns='jabber:client'>
  <body>Yeah, but do you have a license to Jabber?</body>
</message>
`)

func TestChatNoOOB(t *testing.T) {
	var c Client
	c.conn = tConnect(exampleNoOOB)
	c.p = xml.NewDecoder(c.conn)
	m, err := c.Recv()
	if err != nil {
		panic(err)
	}

	t.Logf("Recv() = %v", m)

	switch m.(type) {
	case Chat:
		v, _ := m.(Chat)

		if v.Oob.Url != "" {
			t.Errorf("Should not find URL: `%s`", v.Oob.Url)
		}

		if v.Oob.Desc != "" {
			t.Errorf("Should not find Desc: `%s`", v.Oob.Desc)
		}

	default:
		t.Errorf("Recv() = %v", m)
		t.Errorf("Expected a return value of AvatarData")
	}
}

var rawOob = strings.TrimSpace(`
<x xmlns='jabber:x:oob'>
<url>http://www.jabber.org/images/psa-license.jpg</url>
</x>
`)

func TestRawOob(t *testing.T) {
	var s Oob
	err := xml.Unmarshal([]byte(rawOob), &s)
	if err != nil {
		t.Errorf("%v", err)
	}

	if s.Url != "http://www.jabber.org/images/psa-license.jpg" {
		t.Errorf("Wrong URL: %s", s.Url)
	}
}
