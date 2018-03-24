package latests

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

	"github.com/cryptounicorns/platypus/handlers/handler/latest"
	httpHelper "github.com/cryptounicorns/platypus/http"
)

const (
	Name = "latests"
)

type Latests struct {
	config Config
	format formats.Format
	stores []stores.Store
	wrap   *template.Template
	log    loggers.Logger
}

func (l *Latests) Run(ctx context.Context) error {
	var (
		errs = make(chan error)
	)

	for k, config := range l.config.Inputs {
		go func(config latest.Config, store stores.Store) {
			err := l.runLatest(config, store, ctx)
			if err != nil {
				errs <- err
			}
		}(config, l.stores[k])
	}

	return <-errs
}

func (l *Latests) runLatest(config latest.Config, store stores.Store, ctx context.Context) error {
	var (
		closers = closer.Closers{}
		tpl     *template.Template
		f       formats.Format
		q       queues.Queue
		cr      consumer.Consumer
		st      <-chan result.Result
		k       []byte
		v       interface{}
		err     error
	)

	tpl, err = template.Parse(strings.TrimSpace(config.Key))
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

	f, err = formats.New(config.Format)
	if err != nil {
		return err
	}

	q, err = queues.FromConfig(config.Consumer.Queue, l.log)
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

		v = new(interface{})
		err = f.Unmarshal(r.Value, v)
		if err != nil {
			return err
		}

		k, err = tpl.Apply(struct {
			Config  latest.Config
			Message interface{}
		}{
			Config:  config,
			Message: v,
		})
		if err != nil {
			return err
		}

		err = store.Set(string(k), v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Latests) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx = make(
			[]struct {
				Config latest.Config
				Data   []byte
			},
			len(l.stores),
		)
		values []interface{}
		buf    []byte
		err    error
	)

	rw.Header().Set("Content-Type", "application/"+l.format.Name())

	for k, store := range l.stores {
		values, err = store.Values()
		if httpHelper.HandleError(rw, err, l.log, true) {
			return
		}

		// FIXME: https://github.com/corpix/formats/issues/2
		buf, err = l.format.Marshal(values)
		if httpHelper.HandleError(rw, err, l.log, true) {
			return
		}

		ctx[k].Config = l.config.Inputs[k]
		ctx[k].Data = buf
	}

	buf, err = l.wrap.Apply(ctx)
	if httpHelper.HandleError(rw, err, l.log, true) {
		return
	}

	rw.WriteHeader(http.StatusOK)

	_, err = rw.Write(buf)
	if httpHelper.HandleError(rw, err, l.log, false) {
		return
	}
}

func (l *Latests) Close() error {
	var (
		err error
	)

	for _, store := range l.stores {
		err = store.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func New(c Config, l loggers.Logger) (*Latests, error) {
	var (
		s   = make([]stores.Store, len(c.Inputs))
		f   formats.Format
		t   *template.Template
		err error
	)

	f, err = formats.New(c.Format)
	if err != nil {
		return nil, err
	}

	t, err = template.Parse(strings.TrimSpace(c.Wrap))
	if err != nil {
		return nil, err
	}

	for k, input := range c.Inputs {
		s[k], err = stores.New(input.Store, l)
		if err != nil {
			// FIXME: We should close all previously opened stores
			// on error, could we do this better?
			for _, v := range s {
				v.Close()
			}

			return nil, err
		}
	}

	return &Latests{
		config: c,
		format: f,
		stores: s,
		wrap:   t,
		log:    l,
	}, nil
}
