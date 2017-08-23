package v1

import (
	"io"
	"net/http"

	"github.com/corpix/formats"
	"github.com/corpix/loggers"
	"github.com/corpix/queues"
	"github.com/corpix/trade/markets/market"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/mux"

	"github.com/cryptounicorns/market-fetcher-http/consumer"
	"github.com/cryptounicorns/market-fetcher-http/errors"
	"github.com/cryptounicorns/market-fetcher-http/http/api/config"
	apiTransmitter "github.com/cryptounicorns/market-fetcher-http/http/api/transmitter"
	"github.com/cryptounicorns/market-fetcher-http/stores"
	"github.com/cryptounicorns/market-fetcher-http/transmitters"
	"github.com/cryptounicorns/market-fetcher-http/writerpool"
)

type Tickers struct {
	Router      *mux.Router
	Consumer    *consumer.Consumer
	Store       stores.Store
	Format      formats.Format
	Transmitter transmitters.Transmitter
	WriterPool  *writerpool.WriterPool
	log         loggers.Logger
}

func (t *Tickers) Handle(rw http.ResponseWriter, r *http.Request) {
	var (
		conn io.WriteCloser
		err  error
	)

	conn, _, _, err = ws.UpgradeHTTP(r, rw, nil)
	if err != nil {
		t.log.Error(err)
		return
	}

	err = t.Store.Iter(
		func(k string, v interface{}) {
			var (
				data []byte
				err  error
			)

			data, err = t.Format.Marshal(v)
			if err != nil {
				t.log.Error(err)
				return
			}

			err = wsutil.WriteServerText(
				conn,
				data,
			)
			if err != nil {
				t.log.Error(err)
				return
			}
		},
	)
	if err != nil {
		t.log.Error(err)
		conn.Close()
		return
	}

	t.WriterPool.Add(conn)
}

func (t *Tickers) Close() error {
	var (
		err error
	)

	err = t.Consumer.Close()
	if err != nil {
		return err
	}

	err = t.Transmitter.Close()
	if err != nil {
		return err
	}

	return nil
}

func (t *Tickers) pump() {
	var (
		ticker *market.Ticker
		data   []byte
		key    string
		err    error
	)

	for m := range t.Consumer.Consume() {
		ticker = m.(*market.Ticker)
		key = ticker.CurrencyPair.String() + "|" + ticker.Market
		err = t.Store.Set(key, m)
		if err != nil {
			// XXX: If we can't set into store
			// then it is not critical, just log and go on.
			t.log.Error(err)
		}

		data, err = t.Format.Marshal(m)
		if err != nil {
			t.log.Error(err)
			continue
		}

		_, err = t.Transmitter.Write(data)
		if err != nil {
			t.log.Error(err)
			continue
		}
	}
}

func MountTickers(r *mux.Router, t *Tickers) {
	r.HandleFunc("/stream", t.Handle)
}

func NewTickers(c config.Config, r *mux.Router, q queues.Queue, l loggers.Logger) (*Tickers, error) {
	if r == nil {
		return nil, errors.NewErrNilArgument(r)
	}
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
		apiTransmitter.WriterPoolCleanerErrorHandler(ws, l),
		l,
	)
	if err != nil {
		return nil, err
	}

	cr, err = consumer.New(
		q,
		market.Ticker{},
		f,
		l,
	)
	if err != nil {
		return nil, err
	}

	tickers := &Tickers{
		Router:      r.PathPrefix("/tickers").Subrouter(),
		Consumer:    cr,
		Store:       s,
		Format:      formats.NewJSON(),
		Transmitter: t,
		WriterPool:  ws,
		log:         l,
	}

	MountTickers(
		tickers.Router,
		tickers,
	)

	go tickers.pump()

	return tickers, nil
}
