package nonce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonceString(t *testing.T) {
	var (
		samples = []struct {
			nonce  Nonce
			result string
		}{
			{
				nonce:  Nonce("Ex+lPt0Oh9Re6Z"),
				result: "Ex+lPt0Oh9Re6Z",
			},
		}

		str string
	)

	for _, sample := range samples {
		str = sample.nonce.String()

		assert.Equal(t, sample.result, str)
	}
}
