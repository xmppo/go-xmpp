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
			XMLElement{
				XMLName:  xml.Name{Space: "google:mobile:data", Local: "gcm"},
				InnerXML: "\n\t\t{\"random\": \"&lt;text&gt;\"}\n\t",
			},
			XMLElement{
				XMLName: xml.Name{Space: "jabber:client", Local: "error"},
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
