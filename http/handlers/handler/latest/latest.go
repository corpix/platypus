package latest

import (
	"net/http"

	"github.com/corpix/formats"
	"github.com/corpix/loggers"
	"github.com/cryptounicorns/queues"
	queueConsumer "github.com/cryptounicorns/queues/consumer"
	"github.com/cryptounicorns/queues/result"

	cacheHandler "github.com/cryptounicorns/platypus/http/handlers/cache"
)

type Latest struct {
	config Config
	log    loggers.Logger
	done   chan struct{}

	responseFormat formats.Format
	consumerFormat formats.Format

	Cache    *cacheHandler.Cache
	Queue    queues.Queue
	Consumer queueConsumer.Consumer
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
	if l.returnError(err, rw) {
		return
	}

	rw.Header().Set("Content-Type", "application/"+l.responseFormat.Name())
	rw.WriteHeader(http.StatusOK)

	err = l.Cache.Iter(
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

	_, err = rw.Write(buf)
	if l.returnError(err, rw) {
		return
	}
}

func (l *Latest) consumerWorker() {
	var (
		v      *interface{}
		key    string
		stream <-chan result.Result
		err    error
	)
	stream, err = l.Consumer.Consume()
	if err != nil {
		l.log.Error(err)
		return
	}

	for {
		select {
		case <-l.done:
			return
		case r := <-stream:
			if r.Err != nil {
				l.log.Error(r.Err)
				continue
			}

			v = new(interface{})
			err = l.consumerFormat.Unmarshal(
				r.Value,
				v,
			)
			if err != nil {
				l.log.Error(err)
				continue
			}

			key, err = l.Cache.Set(*v)
			if err != nil {
				l.log.Error(err)
				continue
			}
			l.log.Debugf("Key '%s' set to %+v", key, v)
		}
	}
}

func (l *Latest) Close() error {
	var (
		err error
	)

	err = l.Queue.Close()
	if err != nil {
		return err
	}

	err = l.Consumer.Close()
	if err != nil {
		return err
	}

	err = l.Cache.Close()
	if err != nil {
		return err
	}

	return nil
}

func New(c Config, l loggers.Logger) (*Latest, error) {
	var (
		responseFormat formats.Format
		consumerFormat formats.Format
		ce             *cacheHandler.Cache
		queue          queues.Queue
		consumer       queueConsumer.Consumer
		latest         *Latest
		err            error
	)

	responseFormat, err = formats.New(c.Format)
	if err != nil {
		return nil, err
	}

	consumerFormat, err = formats.New(c.Consumer.Format)
	if err != nil {
		return nil, err
	}

	ce, err = cacheHandler.New(c.Cache, l)
	if err != nil {
		return nil, err
	}

	queue, err = queues.New(c.Consumer.Queue, l)
	if err != nil {
		return nil, err
	}

	consumer, err = queue.Consumer()
	if err != nil {
		return nil, err
	}

	latest = &Latest{
		config: c,
		log:    l,

		responseFormat: responseFormat,
		consumerFormat: consumerFormat,

		Cache:    ce,
		Queue:    queue,
		Consumer: consumer,
	}

	go latest.consumerWorker()

	return latest, nil
}
