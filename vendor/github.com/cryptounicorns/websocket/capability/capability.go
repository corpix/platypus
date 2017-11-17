package capability

// Capability represents capability of websocket endpoint.
// It used by some functions to be more strict when checking
// compatibility with RFC6455.
type Capability uint8

const (
	// ServerSide means that endpoint (caller) is a server.
	ServerSide Capability = 0x1 << iota
	// ServerSide means that endpoint (caller) is a client.
	ClientSide
	// Extended means that extension was negotiated during handshake.
	Extended
	// Fragmented means that endpoint (caller) has received fragmented
	// frame and waits for continuation parts.
	Fragmented
)

// Is checks whether the m has n enabled.
func Has(m Capability, n Capability) bool {
	return m&n != 0
}

// Set enables n capability on m.
func Set(m Capability, n Capability) Capability {
	return m | n
}

// Clear disables n capability on m.
func Clear(m Capability, n Capability) Capability {
	return m & (^n)
}
