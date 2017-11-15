package endpoint

import (
	"fmt"
	"io"
	"net/http"

	"github.com/corpix/formats"
	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/cryptounicorns/queues"
	queuesConsumer "github.com/cryptounicorns/queues/consumer"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	uuid "github.com/satori/go.uuid"

	"github.com/cryptounicorns/platypus/consumer"
	endpointsTransmitter "github.com/cryptounicorns/platypus/http/endpoints/transmitter"
	"github.com/cryptounicorns/platypus/stores"
	"github.com/cryptounicorns/platypus/stores/store"
	"github.com/cryptounicorns/platypus/transmitters"
	"github.com/cryptounicorns/platypus/writerpool"
)

type Handler struct {
	Consumer    *consumer.Consumer
	Store       store.Store
	Format      formats.Format
	Transmitter transmitters.Transmitter
	WriterPool  *writerpool.WriterPool
	log         loggers.Logger
}

func (h *Handler) Handle(rw http.ResponseWriter, r *http.Request) {
	var (
		conn io.WriteCloser
		err  error
	)

	conn, _, _, err = ws.UpgradeHTTP(r, rw, nil)
	if err != nil {
		h.log.Error(err)
		return
	}

	err = h.Store.Iter(
		func(k string, v interface{}) {
			var (
				data []byte
				err  error
			)

			data, err = h.Format.Marshal(v)
			if err != nil {
				h.log.Error(err)
				return
			}

			err = wsutil.WriteServerText(
				conn,
				data,
			)
			if err != nil {
				h.log.Error(err)
				return
			}
		},
	)
	if err != nil {
		h.log.Error(err)
		conn.Close()
		return
	}

	h.WriterPool.Add(conn)
}

func (h *Handler) Close() error {
	var (
		err error
	)

	err = h.Consumer.Close()
	if err != nil {
		return err
	}

	err = h.Transmitter.Close()
	if err != nil {
		return err
	}

	err = h.Store.Close()
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) pump(stream <-chan consumer.Result) {
	var (
		data []byte
		err  error
	)

	for r := range stream {
		if r.Err != nil {
			if r.Err == io.EOF {
				break
			}
			h.log.Error(err)
			continue
		}

		err = h.Store.Set(
			uuid.NewV1().String(),
			r.Value,
		)
		if err != nil {
			// XXX: If we can't set into store
			// then it is not critical, just log and go on.
			h.log.Error(err)
		}

		data, err = h.Format.Marshal(r.Value)
		if err != nil {
			h.log.Error(err)
			continue
		}

		_, err = h.Transmitter.Write(data)
		if err != nil {
			h.log.Error(err)
			continue
		}
	}
}

func NewHandler(c Config, queue queues.Queue, l loggers.Logger) (*Handler, error) {
	var (
		ws  = writerpool.New()
		log = prefixwrapper.New(
			fmt.Sprintf(
				"Handler(%s, %s): ",
				c.Method,
				c.Path,
			),
			l,
		)

		cr     *consumer.Consumer
		qc     queuesConsumer.Consumer
		s      store.Store
		t      transmitters.Transmitter
		h      *Handler
		stream <-chan consumer.Result
		err    error
	)

	s, err = stores.New(
		c.Store,
		log,
	)
	if err != nil {
		return nil, err
	}

	t, err = transmitters.New(
		c.Transmitter,
		ws,
		wsutil.WriteServerText,
		endpointsTransmitter.WriterPoolCleanerErrorHandler(ws, log),
		log,
	)
	if err != nil {
		return nil, err
	}

	qc, err = queue.Consumer()
	if err != nil {
		t.Close()
		return nil, err
	}

	cr, err = consumer.New(qc, c.Consumer)
	if err != nil {
		t.Close()
		qc.Close()
		return nil, err
	}

	stream, err = cr.Consume()
	if err != nil {
		t.Close()
		qc.Close()
		cr.Close()
		return nil, err
	}

	h = &Handler{
		Consumer:    cr,
		Store:       s,
		Format:      formats.NewJSON(),
		Transmitter: t,
		WriterPool:  ws,
		log:         log,
	}

	go h.pump(stream)

	return h, nil
}
