package latests

import (
	"context"
	"net/http"
	"strings"

	"github.com/corpix/formats"
	"github.com/corpix/loggers"
	"github.com/corpix/stores"
	"github.com/corpix/template"
	"github.com/cryptounicorns/queues"

	"github.com/cryptounicorns/platypus/handlers/handler/latest"
	httpHelper "github.com/cryptounicorns/platypus/http"
)

const (
	Name = "latests"
)

type wrapContext struct {
	Input  latest.Config
	Events struct {
		JSON   []byte
		Values []interface{}
	}
}

type Latests struct {
	config Config
	format formats.Format
	stores []stores.Store
	wrap   *template.Template
	log    loggers.Logger
}

func (l *Latests) Run(ctx context.Context) error {
	var (
		errs chan error
	)

	for k, config := range l.config.Inputs {
		go func(config latest.Config, store stores.Store) {
			var (
				tpl *template.Template
				err error
			)

			tpl, err = template.Parse(config.Key)
			if err != nil {
				errs <- err
				return
			}

			errs <- queues.PipeConsumerToStoreWith(
				config.Consumer,
				ctx,
				func(v interface{}) (string, interface{}, error) {
					var (
						key []byte
						err error
					)

					key, err = tpl.Apply(v)
					if err != nil {
						return "", nil, err
					}

					return string(key), v, nil
				},
				store,
				l.log,
			)
		}(config, l.stores[k])
	}

	return <-errs
}

func (l *Latests) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx    = make([]wrapContext, len(l.stores))
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

		ctx[k].Input = l.config.Inputs[k]
		ctx[k].Events.JSON = buf
		ctx[k].Events.Values = values
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
