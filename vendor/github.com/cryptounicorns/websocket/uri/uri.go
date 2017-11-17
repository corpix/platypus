package uri

import (
	"net/url"

	"github.com/cryptounicorns/websocket/uri/scheme"
)

func defaultValue(v string, d string) string {
	if v == "" {
		return d
	}
	return v
}

func New(s string, h string, p string) (*url.URL, error) {
	var (
		meta scheme.Scheme
		err  error
	)

	meta, err = scheme.Get(s)
	if err != nil {
		return nil, err
	}

	return &url.URL{
		Scheme: meta.Scheme,
		Host:   h + ":" + defaultValue(p, meta.Port),
	}, nil
}
