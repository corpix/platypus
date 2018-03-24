package stream

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/corpix/effects/closer"
	"github.com/corpix/effects/writer"
	"github.com/corpix/loggers"
	"github.com/cryptounicorns/queues"
	"github.com/cryptounicorns/queues/consumer"
	"github.com/cryptounicorns/queues/result"
	"github.com/cryptounicorns/websocket"
)

const (
	Name = "stream"
)

type Stream struct {
	*websocket.HTTPUpgradeHandler

	config Config
	writer *writer.ConcurrentMultiWriter
	log    loggers.Logger
}

func (s *Stream) Run(ctx context.Context) error {
	var (
		closers = closer.Closers{}
		q       queues.Queue
		cr      consumer.Consumer
		st      <-chan result.Result
		err     error
	)

	go func() {
		select {
		case <-ctx.Done():
			err := closers.Close()
			if err != nil {
				s.log.Error(err)
			}
			return
		}
	}()

	q, err = queues.FromConfig(s.config.Consumer.Queue, s.log)
	if err != nil {
		return err
	}
	closers = append(closers, q)

	cr, err = q.Consumer()
	if err != nil {
		return err
	}
	closers = append(closers, cr)

	st, err = cr.Consume()
	if err != nil {
		return err
	}

	for r := range st {
		if r.Err != nil {
			return r.Err
		}

		_, err = s.writer.Write(r.Value)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (s *Stream) ServeWebsocket(rwc io.ReadWriteCloser, r *http.Request) {
	s.writer.Add(rwc)
}

func (s *Stream) Close() error {
	// FIXME: Close every connection(this require suport from effects)
	return s.writer.Close()
}

func New(c Config, l loggers.Logger) (*Stream, error) {
	var (
		w      *writer.ConcurrentMultiWriter
		stream *Stream
	)

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
		writer: w,
		log:    l,
	}

	stream.HTTPUpgradeHandler = websocket.NewHTTPUpgradeHandler(
		stream,
		l,
	)

	return stream, nil
}
