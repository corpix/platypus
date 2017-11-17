package nonce

import (
	"crypto/rand"
	"io"
)

var (
	DefaultClientGenerator = NewClientGenerator(rand.Reader, Size)
)

type ClientGenerator struct {
	random io.Reader
	size   int
}

func (g *ClientGenerator) Generate() (Nonce, error) {
	var (
		buf = make([]byte, g.size)
		err error
	)

	_, err = g.random.Read(buf)
	if err != nil {
		return nil, err
	}

	return Nonce(base64Encode(buf)), nil
}

func NewClientGenerator(r io.Reader, s int) *ClientGenerator {
	return &ClientGenerator{
		random: r,
		size:   s,
	}
}
