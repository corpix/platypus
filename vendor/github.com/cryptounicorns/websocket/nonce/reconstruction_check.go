package nonce

import (
	"bytes"
)

func ReconstructionCheck(nonce Nonce, clientNonce Nonce, g *ServerGenerator) error {
	var (
		serverNonce Nonce
		err         error
	)

	serverNonce, err = g.Generate(clientNonce)
	if err != nil {
		return err
	}

	if !bytes.Equal(nonce, serverNonce) {
		return NewErrNonceMismatch(
			serverNonce,
			nonce,
		)
	}

	return nil
}
