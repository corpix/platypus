package nonce

const (
	Size         = 16
	ProtocolUUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

type Nonce []byte

func (n Nonce) String() string {
	return string(n)
}
