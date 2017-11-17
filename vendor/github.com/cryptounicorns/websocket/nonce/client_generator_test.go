package nonce

import (
	"testing"

	"github.com/corpix/effects"
	"github.com/stretchr/testify/assert"
)

func TestDefaultClientGenerator(t *testing.T) {
	assert.NotNil(
		t,
		DefaultClientGenerator,
	)
}

func TestClientGenerator(t *testing.T) {
	var (
		random = effects.NewRingReader(
			[]byte(
				"I've seen things you people wouldn't believe. Attack ships on fire off the shoulder of Orion. I watched C-beams glitter in the dark near the Tannhäuser Gate. All those moments will be lost in time, like tears in rain. Time to die.",
			),
		)

		samples = []struct {
			name      string
			generator *ClientGenerator
			result    []byte
		}{
			{
				name:      "create a nonce from a custom random source",
				generator: NewClientGenerator(random, Size),
				result:    []byte("SSd2ZSBzZWVuIHRoaW5ncw=="),
			},
		}

		nonce []byte
		err   error
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				nonce, err = sample.generator.Generate()
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
