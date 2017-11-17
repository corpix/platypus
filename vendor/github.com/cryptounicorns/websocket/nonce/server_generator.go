package nonce

import (
	"crypto/sha1"
)

var (
	DefaultServerGenerator = NewServerGenerator(ProtocolUUID)
)

type ServerGenerator struct {
	protocolUUID []byte
}

func (g *ServerGenerator) Generate(clientNonce Nonce) (Nonce, error) {
	var (
		hash = sha1.New()
		err  error
	)

	_, err = hash.Write(clientNonce)
	if err != nil {
		return nil, err
	}

	_, err = hash.Write(g.protocolUUID)
	if err != nil {
		return nil, err
	}

	return Nonce(
		base64Encode(
			hash.Sum(nil),
		),
	), nil
}

func NewServerGenerator(protocolUUID string) *ServerGenerator {
	return &ServerGenerator{
		protocolUUID: []byte(protocolUUID),
	}
}
