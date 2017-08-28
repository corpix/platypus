package endpoint

import (
	"io"
	"net/http"

	"github.com/corpix/formats"
	"github.com/corpix/loggers"
	"github.com/corpix/queues"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"github.com/cryptounicorns/market-fetcher-http/consumer"
	"github.com/cryptounicorns/market-fetcher-http/errors"
	endpointsTransmitter "github.com/cryptounicorns/market-fetcher-http/http/endpoints/transmitter"
	"github.com/cryptounicorns/market-fetcher-http/stores"
	"github.com/cryptounicorns/market-fetcher-http/transmitters"
	"github.com/cryptounicorns/market-fetcher-http/writerpool"
)

type Handler struct {
	Consumer    *consumer.Consumer
	Store       stores.Store
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

	return nil
}

func (h *Handler) pump() {
	var (
		// ticker *market.Ticker
		data []byte
		// key    string
		err error
	)

	for m := range h.Consumer.Consume() {
		// FIXME: Use some sort of ring buffer for messages
		// ticker = m.(*market.Ticker)
		// key = ticker.CurrencyPair.String() + "|" + ticker.Market
		// err = t.Store.Set(key, m)
		// if err != nil {
		// 	// XXX: If we can't set into store
		// 	// then it is not critical, just log and go on.
		// 	t.log.Error(err)
		// }

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
		s   stores.Store
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
