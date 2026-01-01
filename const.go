package xmpp

const (
	// Version is current go-xmpp software version.
	Version = "0.3.2"

	// HT_SHA_256_ENDP used in XEP-0484: Fast Authentication Streamlining Tokens, https://xmpp.org/extensions/xep-0484.html
	HT_SHA_256_ENDP = "HT-SHA-256-ENDP"
	// HT_SHA_256_EXPR used in XEP-0484: Fast Authentication Streamlining Tokens, https://xmpp.org/extensions/xep-0484.html
	HT_SHA_256_EXPR = "HT-SHA-256-EXPR"
	// HT_SHA_256_NONE used in XEP-0484: Fast Authentication Streamlining Tokens, https://xmpp.org/extensions/xep-0484.html
	HT_SHA_256_NONE = "HT-SHA-256-NONE"
	// HT_SHA_256_UNIQ used in XEP-0484: Fast Authentication Streamlining Tokens, https://xmpp.org/extensions/xep-0484.html
	HT_SHA_256_UNIQ = "HT-SHA-256-UNIQ"

	// IQTypeError represents iq response type error.
	IQTypeError = "error"
	// IQTypeGet represents iq request type get.
	IQTypeGet = "get"
	// IQTypeResult represents iq response type result.
	IQTypeResult = "result"
	// IQTypeSet represents iq request type set.
	IQTypeSet = "set"

	// SCRAM_SHA_1_PLUS used in SASL auth process, as described in XEP-0388: Extensible SASL Profile, https://xmpp.org/extensions/xep-0388.html
	SCRAM_SHA_1_PLUS = "SCRAM-SHA-1-PLUS"
	// SCRAM_SHA_1 used in SASL auth process, as described in XEP-0388: Extensible SASL Profile, https://xmpp.org/extensions/xep-0388.html
	SCRAM_SHA_1 = "SCRAM-SHA-1"
	// SCRAM_SHA_256_PLUS used in SASL auth process, as described in XEP-0388: Extensible SASL Profile, https://xmpp.org/extensions/xep-0388.html
	SCRAM_SHA_256_PLUS = "SCRAM-SHA-256-PLUS"
	// SCRAM_SHA_256 used in SASL auth process, as described in XEP-0388: Extensible SASL Profile, https://xmpp.org/extensions/xep-0388.html
	SCRAM_SHA_256 = "SCRAM-SHA-256"
	// SCRAM_SHA_512_PLUS used in SASL auth process, as described in XEP-0388: Extensible SASL Profile, https://xmpp.org/extensions/xep-0388.html
	SCRAM_SHA_512_PLUS = "SCRAM-SHA-512-PLUS"
	// SCRAM_SHA_512 used in SASL auth process, as described in XEP-0388: Extensible SASL Profile, https://xmpp.org/extensions/xep-0388.html
	SCRAM_SHA_512 = "SCRAM-SHA-512"
	// UPGR_SCRAM_SHA_256 used in SASL auth process, as described in XEP-0388: Extensible SASL Profile, https://xmpp.org/extensions/xep-0388.html
	UPGR_SCRAM_SHA_256 = "UPGR-SCRAM-SHA-256"
	// UPGR_SCRAM_SHA_512 used in SASL auth process, as described in XEP-0388: Extensible SASL Profile, https://xmpp.org/extensions/xep-0388.html
	UPGR_SCRAM_SHA_512 = "UPGR-SCRAM-SHA-512"

	// XMPPNS_AVATAR_PEP_DATA represents xml namespace for https://xmpp.org/extensions/xep-0084.html#process-pubdata
	// (User Publishes Data).
	XMPPNS_AVATAR_PEP_DATA = "urn:xmpp:avatar:data"
	// XMPPNS_AVATAR_PEP_METADATA represents xml namespace for https://xmpp.org/extensions/xep-0084.html#process-pubmeta
	// (User Publishes Metadata Notification).
	XMPPNS_AVATAR_PEP_METADATA = "urn:xmpp:avatar:metadata"
	// XMPPNS_BIND_0 used in bunding resource identifier to a session as described in https://xmpp.org/extensions/xep-0386.html
	XMPPNS_BIND_0 = "urn:xmpp:bind:0"
	// XMPPNS_CLIENT namespace is a foundational XML namespace used in the Extensible Messaging and Presence Protocol
	// (XMPP) to scope the core client-to-server (C2S) communication stanzas.
	XMPPNS_CLIENT = "jabber:client"
	// XMPPNS_DISCO_INFO namespace used in Service Discovery protocol, https://xmpp.org/extensions/xep-0030.html
	XMPPNS_DISCO_INFO = "http://jabber.org/protocol/disco#info"
	// XMPPNS_DISCO_ITEMS namespace used in item discover queries, described https://xmpp.org/extensions/xep-0030.html#items
	XMPPNS_DISCO_ITEMS = "http://jabber.org/protocol/disco#items"
	// XMPPNS_FAST_0 namespace used in XEP-0484: Fast Authentication Streamlining Tokens, https://xmpp.org/extensions/xep-0484.html
	XMPPNS_FAST_0 = "urn:xmpp:fast:0"
	// XMPPNS_HTTP_UPLOAD_0 namespace used in XEP-0363: HTTP File Upload, https://xmpp.org/extensions/xep-0363.html
	XMPPNS_HTTP_UPLOAD_0 = "urn:xmpp:http:upload:0"
	// XMPPNS_IQ_VERSION namespace used in XEP-0092: Software Version, https://xmpp.org/extensions/xep-0092.html
	XMPPNS_IQ_VERSION = "jabber:iq:version"
	// XMPPNS_MUC namespace used in XEP-0045: Multi-User Chat, https://xmpp.org/extensions/xep-0045.html
	XMPPNS_MUC = "http://jabber.org/protocol/muc"
	// XMPPNS_MUC_USER namespace used in XEP-0045: Multi-User Chat, https://xmpp.org/extensions/xep-0045.html
	XMPPNS_MUC_USER = "http://jabber.org/protocol/muc#user"
	// XMPPNS_PING namespace used in XEP-0199: XMPP Ping, https://xmpp.org/extensions/xep-0199.html
	XMPPNS_PING = "urn:xmpp:ping"
	// XMPPNS_PUBSUB_EVENT namespace used in XEP-0060: Publish-Subscribe, https://xmpp.org/extensions/xep-0060.html
	XMPPNS_PUBSUB_EVENT = "http://jabber.org/protocol/pubsub#event"
	// XMPPNS_PUBSUB namespace used in XEP-0060: Publish-Subscribe, https://xmpp.org/extensions/xep-0060.html
	XMPPNS_PUBSUB = "http://jabber.org/protocol/pubsub"
	// XMPPNS_SASL_2 namespace used during SASL auth, as described in XEP-0388: Extensible SASL Profile, https://xmpp.org/extensions/xep-0388.html
	XMPPNS_SASL_2 = "urn:xmpp:sasl:2"
	// XMPPNS_SASL_CB_0 namespace used during SASL auth, as described in XEP-0388: Extensible SASL Profile, https://xmpp.org/extensions/xep-0388.html
	// and XEP-0440: SASL Channel-Binding Type Capability, https://xmpp.org/extensions/xep-0440.html
	XMPPNS_SASL_CB_0 = "urn:xmpp:sasl-cb:0"
	// XMPPNS_SASL_UPGRADE_0 namespace used during SASL auth, as described in XEP-0388: Extensible SASL Profile, https://xmpp.org/extensions/xep-0388.html
	// and XEP-0480: SASL Upgrade Tasks, https://xmpp.org/extensions/xep-0480.html
	XMPPNS_SASL_UPGRADE_0 = "urn:xmpp:sasl:upgrade:0"
	// XMPPNS_XMPP_SASL namespace used during SASL auth, as described in XEP-0388: Extensible SASL Profile, https://xmpp.org/extensions/xep-0388.html
	XMPPNS_XMPP_SASL = "urn:ietf:params:xml:ns:xmpp-sasl"
	// XMPPNS_SCRAM_UPGRADE_0 namespace used during SASL auth, as described in XEP-0388: Extensible SASL Profile, https://xmpp.org/extensions/xep-0388.html
	// and XEP-0480: SASL Upgrade Tasks, https://xmpp.org/extensions/xep-0480.html
	XMPPNS_SCRAM_UPGRADE_0 = "urn:xmpp:scram-upgrade:0"
	// XMPPNS_SID_0 namespace used in XEP-0359: Unique and Stable Stanza IDs, https://xmpp.org/extensions/xep-0359.html
	// to implement unique and relible id.
	XMPPNS_SID_0 = "urn:xmpp:sid:0"
	// XMPPNS_STREAM namespace used in description of clients xml data stream as described in XEP-0044: Full Namespace
	// Support for XML Streams, https://xmpp.org/extensions/xep-0044.html
	XMPPNS_STREAM = "http://etherx.jabber.org/streams"
	// XMPPNS_STREAM_LIMITS_0 namespace used during stream advertisement limits (on early connection init stage) as
	// described in XEP-0478: Stream Limits Advertisement, https://xmpp.org/extensions/xep-0478.html
	XMPPNS_STREAM_LIMITS_0 = "urn:xmpp:stream-limits:0"
	// XMPPNS_TIME namespace used in response to information query for client local time, as described in XEP-0202: Entity Time, https://xmpp.org/extensions/xep-0202.html
	XMPPNS_TIME = "urn:xmpp:time"
	// XMPPNS_XMPP_TLS namespace used during session initialization to start tls session, as described in https://www.ietf.org/rfc/rfc6120.txt
	XMPPNS_XMPP_TLS = "urn:ietf:params:xml:ns:xmpp-tls"
	// XMPPNS_XMPP_BIND namespace used during session initialization to start tls session, as described in https://www.ietf.org/rfc/rfc6120.txt
	XMPPNS_XMPP_BIND = "urn:ietf:params:xml:ns:xmpp-bind"
	// XMPPNS_XMPP_SESSION namespace used during xmpp session establisment process, as described in https://www.ietf.org/rfc/rfc6121.txt
	XMPPNS_XMPP_SESSION = "urn:ietf:params:xml:ns:xmpp-session"
)
