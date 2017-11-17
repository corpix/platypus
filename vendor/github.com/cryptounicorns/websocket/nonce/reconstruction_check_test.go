package nonce

import (
	"testing"
)

func TestReconstructionCheck(t *testing.T) {
	var (
		samples = []struct {
			name        string
			nonce       Nonce
			clientNonce Nonce
			generator   *ServerGenerator
		}{
			{
				name: "",
			},
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {

			},
		)
	}
}
