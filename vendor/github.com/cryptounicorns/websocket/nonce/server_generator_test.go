package nonce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultServerGenerator(t *testing.T) {
	assert.NotNil(
		t,
		DefaultServerGenerator,
	)
}

func TestServerGenerator(t *testing.T) {
	var (
		samples = []struct {
			name        string
			clientNonce Nonce
			generator   *ServerGenerator
			result      []byte
		}{
			{
				name:        "create a nonce from rfc example",
				clientNonce: Nonce("dGhlIHNhbXBsZSBub25jZQ=="),
				generator:   NewServerGenerator(ProtocolUUID),
				result:      []byte("s3pPLMBiTxaQ9kYGzzhZRbK+xOo="),
			},
			{
				name:        "create a nonce from a custom client nonce",
				clientNonce: Nonce("SSd2ZSBzZWVuIHRoaW5ncw=="),
				generator:   NewServerGenerator(ProtocolUUID),
				result:      []byte("Ex+lPt0Oh9Re6ZVfbK3u4sqR8zg="),
			},
		}

		nonce []byte
		err   error
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				nonce, err = sample.generator.Generate(sample.clientNonce)
				if err != nil {
					t.Error(err)
					return
				}

				assert.EqualValues(
					t,
					sample.result,
					nonce,
				)
			},
		)
	}
}
