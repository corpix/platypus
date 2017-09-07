package endpoint

import (
	"io"
	"net/http"

	"github.com/corpix/formats"
	"github.com/corpix/loggers"
	"github.com/corpix/queues"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	uuid "github.com/satori/go.uuid"

	"github.com/cryptounicorns/platypus/consumer"
	"github.com/cryptounicorns/platypus/errors"
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

func (h *Handler) pump() {
	var (
		data []byte
		err  error
	)

	for m := range h.Consumer.Consume() {
		err = h.Store.Set(
			uuid.NewV1().String(),
			m,
		)
		if err != nil {
			// XXX: If we can't set into store
			// then it is not critical, just log and go on.
			h.log.Error(err)
		}

		data, err = h.Format.Marshal(m)
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

func NewHandler(c Config, q queues.Queue, l loggers.Logger) (*Handler, error) {
	if q == nil {
		return nil, errors.NewErrNilArgument(q)
	}
	if l == nil {
		return nil, errors.NewErrNilArgument(l)
	}

	var (
		ws  = writerpool.New()
		f   formats.Format
		cr  *consumer.Consumer
		s   store.Store
		t   transmitters.Transmitter
		h   *Handler
		err error
	)

	f, err = formats.New(c.Consumer.Format)
	if err != nil {
		return nil, err
	}

	s, err = stores.New(
		c.Store,
		l,
	)
	if err != nil {
		return nil, err
	}

	t, err = transmitters.New(
		c.Transmitter,
		ws,
		wsutil.WriteServerText,
		endpointsTransmitter.WriterPoolCleanerErrorHandler(ws, l),
		l,
	)
	if err != nil {
		return nil, err
	}

	cr, err = consumer.New(
		q,
		new(interface{}),
		f,
		l,
	)
	if err != nil {
		return nil, err
	}

	h = &Handler{
		Consumer:    cr,
		Store:       s,
		Format:      formats.NewJSON(),
		Transmitter: t,
		WriterPool:  ws,
		log:         l,
	}

	go h.pump()

	return h, nil
}
