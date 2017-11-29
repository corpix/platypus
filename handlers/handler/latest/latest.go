package latest

import (
	"context"
	"net/http"
	"strings"

	"github.com/corpix/formats"
	"github.com/corpix/loggers"
	"github.com/corpix/stores"
	"github.com/corpix/template"
	"github.com/cryptounicorns/queues"

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
		tpl *template.Template
		err error
	)

	tpl, err = template.Parse(strings.TrimSpace(l.config.Key))
	if err != nil {
		return err
	}

	return queues.PipeConsumerToStoreWith(
		l.config.Consumer,
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
		l.store,
		l.log,
	)
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
