package xmpp

import (
    "fmt"
)

func (c* Client) PingC2S(jid, server string) {
    fmt.Fprintf(c.conn, "<iq from='%s' to='%s' id='c2s1' type='get'>\n" +
            "<ping xmlns='urn:xmpp:ping'/>\n" +
            "</iq>",
            xmlEscape(jid), xmlEscape(server))
}

func (c* Client) PingS2S(fromServer, toServer string) {
    fmt.Fprintf(c.conn, "<iq from='%s' to='%s' id='s2s1' type='get'>\n" +
            "<ping xmlns='urn:xmpp:ping'/>\n" +
            "</iq>",
            xmlEscape(fromServer), xmlEscape(toServer))
}
