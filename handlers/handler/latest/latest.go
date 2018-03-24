package latest

import (
	"context"
	"net/http"
	"strings"

	"github.com/corpix/effects/closer"
	"github.com/corpix/formats"
	"github.com/corpix/loggers"
	"github.com/corpix/stores"
	"github.com/corpix/template"
	"github.com/cryptounicorns/queues"
	"github.com/cryptounicorns/queues/consumer"
	"github.com/cryptounicorns/queues/result"

	httpHelper "github.com/cryptounicorns/platypus/http"
)

const (
	Name = "latest"
)

type Latest struct {
	config Config
	format formats.Format
	store  stores.Store
	log    loggers.Logger
}

func (l *Latest) Run(ctx context.Context) error {
	var (
		closers = closer.Closers{}
		tpl     *template.Template
		f       formats.Format
		q       queues.Queue
		cr      consumer.Consumer
		mcr     consumer.Generic
		st      <-chan result.Generic
		k       []byte
		v       interface{}
		err     error
	)

	tpl, err = template.Parse(strings.TrimSpace(l.config.Key))
	if err != nil {
		return err
	}

	go func() {
		select {
		case <-ctx.Done():
			err := closers.Close()
			if err != nil {
				l.log.Error(err)
			}
			return
		}
	}()

	f, err = formats.New(l.config.Format)
	if err != nil {
		return err
	}

	q, err = queues.FromConfig(l.config.Consumer.Queue, l.log)
	if err != nil {
		return err
	}
	closers = append(closers, q)

	cr, err = q.Consumer()
	if err != nil {
		return err
	}
	closers = append(closers, cr)

	mcr = consumer.NewUnmarshal(cr, new(interface{}), f)
	closers = append(closers, mcr)

	st, err = mcr.Consume()
	if err != nil {
		return err
	}

	for r := range st {
		if r.Err != nil {
			return r.Err
		}

		k, err = tpl.Apply(struct {
			Config  Config
			Message interface{}
		}{
			Config:  l.config,
			Message: r.Value,
		})
		if err != nil {
			return err
		}

		err = l.store.Set(string(k), v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Latest) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var (
		values []interface{}
		buf    []byte
		err    error
	)

	rw.Header().Set("Content-Type", "application/"+l.format.Name())

	values, err = l.store.Values()
	if httpHelper.HandleError(rw, err, l.log, true) {
		return
	}

	// FIXME: https://github.com/corpix/formats/issues/2
	buf, err = l.format.Marshal(values)
	if httpHelper.HandleError(rw, err, l.log, true) {
		return
	}

	rw.WriteHeader(http.StatusOK)

	_, err = rw.Write(buf)
	if httpHelper.HandleError(rw, err, l.log, false) {
		return
	}
}

func (l *Latest) Close() error {
	return l.store.Close()
}

func New(c Config, l loggers.Logger) (*Latest, error) {
	var (
		f   formats.Format
		s   stores.Store
		err error
	)

	f, err = formats.New(c.Format)
	if err != nil {
		return nil, err
	}

	s, err = stores.New(c.Store, l)
	if err != nil {
		return nil, err
	}

	return &Latest{
		config: c,
		format: f,
		store:  s,
		log:    l,
	}, nil
}
