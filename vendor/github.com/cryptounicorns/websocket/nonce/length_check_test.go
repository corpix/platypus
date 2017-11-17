package nonce

import (
	"testing"
)

func TestLengthCheck(t *testing.T) {
	var (
		samples = []struct {
			name          string
			data          string
			decodedLength int
		}{
			{
				name:          "empty",
				data:          "",
				decodedLength: 0,
			},
			{
				name:          "not empty",
				data:          "MTIz",
				decodedLength: 3,
			},
			{
				name:          "random not empty",
				data:          "uBRcDtmO805ixG5IRdQvug==",
				decodedLength: 16,
			},
		}

		err error
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				err = LengthCheck(
					Nonce(sample.data),
					sample.decodedLength,
				)

				if err != nil {
					t.Error(err)
				}
			},
		)
	}
}
