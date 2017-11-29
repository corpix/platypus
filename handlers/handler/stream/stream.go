package stream

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/corpix/effects/writer"
	"github.com/corpix/formats"
	"github.com/corpix/loggers"
	"github.com/cryptounicorns/queues"
	"github.com/cryptounicorns/websocket"
	websocketWriter "github.com/cryptounicorns/websocket/writer"
)

const (
	Name = "stream"
)

type Stream struct {
	*websocket.UpgradeHandler

	config Config
	format formats.Format
	writer *writer.ConcurrentMultiWriter
	log    loggers.Logger
}

func (s *Stream) Run(ctx context.Context) error {
	return queues.PipeConsumerToWriterWith(
		s.config.Consumer,
		ctx,
		s.format.Marshal,
		s.writer,
		s.log,
	)
}

func (s *Stream) ServeWebsocket(c io.WriteCloser, r *http.Request) {
	s.writer.Add(websocketWriter.NewServerText(c))
}

func (s *Stream) Close() error {
	// FIXME: Close every connection(this require suport from effects)
	return s.writer.Close()
}

func New(c Config, l loggers.Logger) (*Stream, error) {
	var (
		f      formats.Format
		w      *writer.ConcurrentMultiWriter
		stream *Stream
		err    error
	)

	f, err = formats.New(c.Format)
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

	stream = &Stream{
		config: c,
		format: f,
		writer: w,
		log:    l,
	}

	stream.UpgradeHandler = websocket.NewUpgradeHandler(
		stream,
		l,
	)

	return stream, nil
}
