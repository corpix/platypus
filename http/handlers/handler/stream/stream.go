package stream

import (
	"net/http"

	"github.com/corpix/loggers"
)

type Stream struct{}

func (s *Stream) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

}

func (s *Stream) Close() error {
	return nil
}

func New(c Config, l loggers.Logger) (*Stream, error) {
	return &Stream{}, nil
}
