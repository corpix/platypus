package v1

import (
	"io"
	"net/http"

	"github.com/corpix/formats"
	"github.com/corpix/logger"
	"github.com/corpix/queues"
	"github.com/corpix/queues/consumer"
	"github.com/corpix/trade/market"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/mux"

	"github.com/cryptounicorns/market-fetcher-http/errors"
	"github.com/cryptounicorns/market-fetcher-http/http/api/config"
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
	log         logger.Logger
	done        chan bool
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

	close(t.done)

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
		feed = t.Consumer.GetFeed()
		data []byte
		key  string
		err  error
	)

	for {
		select {
		case <-t.done:
			break
		case m := <-feed:
			key = m.(*market.Ticker).CurrencyPair.String()
			t.Store.Set(key, m)

			data, err = t.Format.Marshal(m)
			if err != nil {
				t.log.Error(err)
				continue
			}

			_, err = t.Transmitter.Write(data)
			if err != nil {
				// FIXME: Delete writer if broken pipe
				t.log.Error(err)
				continue
			}
		}
	}
}

func MountTickers(r *mux.Router, t *Tickers) {
	r.HandleFunc("/stream", t.Handle)
}

func NewTickers(c config.Config, r *mux.Router, q queues.Queue, l logger.Logger) (*Tickers, error) {
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

	cr, err = consumer.New(
		market.Ticker{},
		f,
		l,
	)
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
		l,
	)
	if err != nil {
		return nil, err
	}

	err = q.Consume(cr.Handler)
	if err != nil {
		cr.Close()
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
		done:        make(chan bool),
	}

	MountTickers(
		tickers.Router,
		tickers,
	)

	go tickers.pump()

	return tickers, nil
}
