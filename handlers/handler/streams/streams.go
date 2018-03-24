package streams

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/corpix/effects/closer"
	"github.com/corpix/effects/writer"
	"github.com/corpix/loggers"
	"github.com/corpix/template"
	"github.com/cryptounicorns/queues"
	"github.com/cryptounicorns/queues/consumer"
	"github.com/cryptounicorns/queues/result"
	"github.com/cryptounicorns/websocket"

	"github.com/cryptounicorns/platypus/handlers/handler/stream"
)

const (
	Name = "streams"
)

type Streams struct {
	*websocket.HTTPUpgradeHandler

	config Config
	wrap   *template.Template
	writer *writer.ConcurrentMultiWriter
	log    loggers.Logger
}

func (s *Streams) Run(ctx context.Context) error {
	var (
		errs = make(chan error)
	)

	for _, config := range s.config.Inputs {
		go func(config stream.Config) {
			err := s.runStream(config, ctx)
			if err != nil {
				errs <- err
			}
		}(config)
	}

	return <-errs
}

func (s *Streams) runStream(config stream.Config, ctx context.Context) error {
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

	q, err = queues.FromConfig(config.Consumer.Queue, s.log)
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

		var (
			buf []byte
			err error
		)

		if s.wrap != nil {
			buf, err = s.wrap.Apply(struct {
				Config  stream.Config
				Message []byte
			}{
				Config:  config,
				Message: r.Value,
			})
			if err != nil {
				return err
			}
		} else {
			buf = r.Value
		}

		_, err = s.writer.Write(buf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Streams) ServeWebsocket(rwc io.ReadWriteCloser, r *http.Request) {
	s.writer.Add(rwc)
}

func (s *Streams) Close() error {
	return s.writer.Close()
}

func New(c Config, l loggers.Logger) (*Streams, error) {
	var (
		t       *template.Template
		w       *writer.ConcurrentMultiWriter
		streams *Streams
		err     error
	)

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
		wrap:   t,
		writer: w,
		log:    l,
	}

	streams.HTTPUpgradeHandler = websocket.NewHTTPUpgradeHandler(
		streams,
		l,
	)

	return streams, nil
}
