package latests

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/corpix/formats"
	"github.com/corpix/loggers"

	"github.com/cryptounicorns/platypus/http/handlers/cache"
	"github.com/cryptounicorns/platypus/http/handlers/consumer"
	"github.com/cryptounicorns/platypus/http/handlers/memoize"
)

type templateEventContext struct {
	JSON  []byte
	Value []interface{}
}
type templateContext struct {
	Consumer consumer.Config
	Events   templateEventContext
}

type Latests struct {
	config         Config
	log            loggers.Logger
	done           chan struct{}
	wrap           *template.Template
	responseFormat formats.Format
	Memoize        []memoize.Memoize
}

func (l *Latests) returnError(err error, rw http.ResponseWriter, sendHeader bool) bool {
	if err != nil {
		l.log.Error(err)
		if sendHeader {
			rw.WriteHeader(500)
		}
		return true
	}

	return false
}

func (l *Latests) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx = make(
			[]templateContext,
			len(l.Memoize),
		)
		err error
	)

	rw.Header().Set("Content-Type", "application/"+l.responseFormat.Name())

	for k, mz := range l.Memoize {
		ctx[k] = templateContext{
			Consumer: mz.Consumer.Meta.Config,
		}
		// FIXME: Could be parallel
		err = mz.Cache.Iter(
			func(key string, value interface{}) {
				ctx[k].Events.Value = append(
					ctx[k].Events.Value,
					value,
				)
			},
		)
		if l.returnError(err, rw, true) {
			return
		}

		ctx[k].Events.JSON, err = l.responseFormat.Marshal(&ctx[k].Events.Value)
		if l.returnError(err, rw, true) {
			return
		}
	}

	rw.WriteHeader(http.StatusOK)

	err = l.wrap.Execute(rw, &ctx)
	if l.returnError(err, rw, false) {
		return
	}
}

func (l *Latests) consumersWorkers() {
	for _, mz := range l.Memoize {
		go l.consumerWorker(mz)
	}
}

func (l *Latests) consumerWorker(memoize memoize.Memoize) {
	var (
		key string
		err error
	)

	for {
		select {
		case <-l.done:
			return
		case r := <-memoize.Consumer.Stream():
			if r.Err != nil {
				l.log.Error(r.Err)
				return
			}

			if err != nil {
				l.log.Error(err)
				continue
			}

			key, err = memoize.Cache.Set(r.Value)
			if err != nil {
				l.log.Error(err)
				continue
			}

			l.log.Debugf(
				"Key '%s' set to %+v",
				key,
				r.Value,
			)
		}
	}
}

func (l *Latests) Close() error {
	var (
		err error
	)

	close(l.done)

	for _, v := range l.Memoize {
		err = v.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func New(c Config, l loggers.Logger) (*Latests, error) {
	var (
		mz = make(
			[]memoize.Memoize,
			len(c.Memoize),
		)
		rf      formats.Format
		ce      *cache.Cache
		cr      *consumer.Consumer
		t       *template.Template
		latests *Latests
		err     error
	)

	rf, err = formats.New(c.Format)
	if err != nil {
		return nil, err
	}

	t, err = template.New("wrap").Parse(
		strings.TrimSpace(c.Wrap),
	)
	if err != nil {
		return nil, err
	}

	for k, v := range c.Memoize {
		ce, err = cache.New(v.Cache, l)
		if err != nil {
			return nil, err
		}

		cr, err = consumer.New(v.Consumer, l)
		if err != nil {
			return nil, err
		}

		mz[k] = memoize.Memoize{
			Cache:    ce,
			Consumer: cr,
		}
	}

	latests = &Latests{
		config:         c,
		log:            l,
		done:           make(chan struct{}),
		wrap:           t,
		responseFormat: rf,
		Memoize:        mz,
	}

	latests.consumersWorkers()

	return latests, nil
}
