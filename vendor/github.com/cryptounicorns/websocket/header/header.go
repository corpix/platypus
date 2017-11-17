package header

import (
	"net/textproto"

	"github.com/cryptounicorns/websocket/capability"
	"github.com/cryptounicorns/websocket/frame/opcode"
)

var (
	Upgrade       = textproto.CanonicalMIMEHeaderKey("Upgrade")
	Connection    = textproto.CanonicalMIMEHeaderKey("Connection")
	Host          = textproto.CanonicalMIMEHeaderKey("Host")
	Origin        = textproto.CanonicalMIMEHeaderKey("Origin")
	SecVersion    = textproto.CanonicalMIMEHeaderKey("Sec-Websocket-Version")
	SecProtocol   = textproto.CanonicalMIMEHeaderKey("Sec-Websocket-Protocol")
	SecExtensions = textproto.CanonicalMIMEHeaderKey("Sec-Websocket-Extensions")
	SecKey        = textproto.CanonicalMIMEHeaderKey("Sec-Websocket-Key")
	SecAccept     = textproto.CanonicalMIMEHeaderKey("Sec-Websocket-Accept")
)

// Header represents websocket frame header.
// See https://tools.ietf.org/html/rfc6455#section-5.2
type Header struct {
	// FIXME: SHould they be there?
	Fin    bool
	Rsv    byte
	OpCode opcode.OpCode
	Length int64
	Masked bool
	Mask   [4]byte
}

// CheckHeader checks h to contain valid header data for given capability c.
//
// Note that zero capability (0) means that capability is clean,
// neither server or client side, nor fragmented, nor extended.
func (h Header) Check(c Capability) error {
	if h.OpCode.IsReserved() {
		return ErrProtocolOpCodeReserved
	}
	if h.OpCode.IsControl() {
		if h.Length > MaxControlFramePayloadSize {
			return ErrProtocolControlPayloadOverflow
		}
		if !h.Fin {
			return ErrProtocolControlNotFinal
		}
	}

	switch {
	// [RFC6455]: MUST be 0 unless an extension is negotiated that defines meanings for
	// non-zero values. If a nonzero value is received and none of the
	// negotiated extensions defines the meaning of such a nonzero value, the
	// receiving endpoint MUST _Fail the WebSocket Connection_.
	case h.Rsv != 0 && !capability.Has(c, Capability.Extended):
		return ErrProtocolNonZeroRsv

	// [RFC6455]: The server MUST close the connection upon receiving a frame that is not masked.
	// In this case, a server MAY send a Close frame with a status code of 1002 (protocol error)
	// as defined in Section 7.4.1. A server MUST NOT mask any frames that it sends to the client.
	// A client MUST close a connection if it detects a masked frame. In this case, it MAY use the
	// status code 1002 (protocol error) as defined in Section 7.4.1.
	case capability.Has(c, Capability.ServerSide) && !h.Masked:
		return ErrProtocolMaskRequired
	case capability.Has(c, Capability.ClientSide) && h.Masked:
		return ErrProtocolMaskUnexpected

		// [RFC6455]: See detailed explanation in 5.4 section.
	case capability.Has(c, Capability.Fragmented) && !h.OpCode.IsControl() &&
		h.OpCode != OpContinuation:
		return ErrProtocolContinuationExpected
	case !capability.Has(c, Capability.Fragmented) && h.OpCode == OpContinuation:
		return ErrProtocolContinuationUnexpected
	}

	return nil
}
