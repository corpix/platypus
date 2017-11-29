package streams

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/corpix/effects/writer"
	"github.com/corpix/formats"
	"github.com/corpix/loggers"
	"github.com/corpix/template"
	"github.com/cryptounicorns/queues"
	"github.com/cryptounicorns/websocket"
	websocketWriter "github.com/cryptounicorns/websocket/writer"

	"github.com/cryptounicorns/platypus/handlers/handler/stream"
)

const (
	Name = "streams"
)

type wrapContext struct {
	Input stream.Config
	Event struct {
		JSON  []byte
		Value interface{}
	}
}

type Streams struct {
	*websocket.UpgradeHandler

	config Config
	format formats.Format
	wrap   *template.Template
	writer *writer.ConcurrentMultiWriter
	log    loggers.Logger
}

func (s *Streams) Run(ctx context.Context) error {
	var (
		errs chan error
	)

	for _, config := range s.config.Inputs {
		go func(config stream.Config) {
			errs <- queues.PipeConsumerToWriterWith(
				config.Consumer,
				ctx,
				func(v interface{}) ([]byte, error) {
					var (
						ctx = wrapContext{}
						buf []byte
						err error
					)

					buf, err = s.format.Marshal(v)
					if err != nil {
						return nil, err
					}

					ctx.Input = config
					ctx.Event.JSON = buf
					ctx.Event.Value = v

					return s.wrap.Apply(ctx)
				},
				s.writer,
				s.log,
			)
		}(config)
	}

	return <-errs

}

func (s *Streams) ServeWebsocket(c io.WriteCloser, r *http.Request) {
	s.writer.Add(websocketWriter.NewServerText(c))
}

func (s *Streams) Close() error {
	return s.writer.Close()
}

func New(c Config, l loggers.Logger) (*Streams, error) {
	var (
		f       formats.Format
		t       *template.Template
		w       *writer.ConcurrentMultiWriter
		streams *Streams
		err     error
	)

	f, err = formats.New(c.Format)
	if err != nil {
		return nil, err
	}

	t, err = template.Parse(strings.TrimSpace(c.Wrap))
	if err != nil {
		return nil, err
	}

	w = writer.NewConcurrentMultiWriter(
		c.Writer,
		func(mw *writer.ConcurrentMultiWriter, w io.Writer, err error) {
			mw.Remove(w)
			w.(io.Closer).Close()

			var (
				ok bool
			)

			_, ok = err.(*net.OpError)
			if !ok {
				l.Error(err)
			}
		},
	)

	streams = &Streams{
		config: c,
		format: f,
		wrap:   t,
		writer: w,
		log:    l,
	}

	streams.UpgradeHandler = websocket.NewUpgradeHandler(
		streams,
		l,
	)

	return streams, nil
}
