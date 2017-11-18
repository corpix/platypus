package latest

import (
	"net/http"

	"github.com/corpix/formats"
	"github.com/corpix/loggers"

	"github.com/cryptounicorns/platypus/http/handlers/cache"
	"github.com/cryptounicorns/platypus/http/handlers/consumer"
	"github.com/cryptounicorns/platypus/http/handlers/memoize"
)

type Latest struct {
	config         Config
	log            loggers.Logger
	done           chan struct{}
	responseFormat formats.Format
	Memoize        memoize.Memoize
}

func (l *Latest) returnError(err error, rw http.ResponseWriter) bool {
	if err != nil {
		l.log.Error(err)
		rw.WriteHeader(500)
		return true
	}

	return false
}

func (l *Latest) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var (
		res []interface{}
		buf []byte
		err error
	)

	rw.Header().Set("Content-Type", "application/"+l.responseFormat.Name())

	err = l.Memoize.Cache.Iter(
		func(key string, value interface{}) {
			res = append(
				res,
				value,
			)
		},
	)
	if l.returnError(err, rw) {
		return
	}

	buf, err = l.responseFormat.Marshal(&res)
	if l.returnError(err, rw) {
		return
	}

	rw.WriteHeader(http.StatusOK)

	_, err = rw.Write(buf)
	if l.returnError(err, rw) {
		return
	}
}

func (l *Latest) consumerWorker() {
	var (
		key string
		err error
	)

	for {
		select {
		case <-l.done:
			return
		case r := <-l.Memoize.Consumer.Stream():
			if r.Err != nil {
				l.log.Error(r.Err)
				continue
			}

			if err != nil {
				l.log.Error(err)
				continue
			}

			key, err = l.Memoize.Cache.Set(r.Value)
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

func (l *Latest) Close() error {
	close(l.done)
	return l.Memoize.Close()
}

func New(c Config, l loggers.Logger) (*Latest, error) {
	var (
		rf     formats.Format
		ce     *cache.Cache
		cr     *consumer.Consumer
		latest *Latest
		err    error
	)

	rf, err = formats.New(c.Format)
	if err != nil {
		return nil, err
	}

	ce, err = cache.New(c.Memoize.Cache, l)
	if err != nil {
		return nil, err
	}

	cr, err = consumer.New(c.Memoize.Consumer, l)
	if err != nil {
		return nil, err
	}

	latest = &Latest{
		config:         c,
		log:            l,
		done:           make(chan struct{}),
		responseFormat: rf,
		Memoize: memoize.Memoize{
			Cache:    ce,
			Consumer: cr,
		},
	}

	go latest.consumerWorker()

	return latest, nil
}
