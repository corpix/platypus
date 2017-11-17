package nonce

import (
	"encoding/base64"
	"github.com/davecgh/go-spew/spew"
)

func base64Encode(buf []byte) []byte {
	var (
		res = make(
			[]byte,
			base64.StdEncoding.EncodedLen(
				len(buf),
			),
		)
	)

	base64.StdEncoding.Encode(res, buf)

	return res
}

func base64Decode(buf []byte) ([]byte, error) {
	var (
		res = make(
			[]byte,
			base64.StdEncoding.DecodedLen(
				len(buf),
			),
		)

		n   int
		err error
	)

	n, err = base64.StdEncoding.Decode(res, buf)
	if err != nil {
		return nil, err
	}

	spew.Dump(buf, res, "---")

	return res
}

func base64DecodedLen(buf []byte) (int, error) {
	var (
		n   int
		err error
	)

	n, err = base64Decode(buf)
	if err != nil {
		return nil, err
	}

	return len(n), nil
}
