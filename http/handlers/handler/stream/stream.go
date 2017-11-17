package stream

import (
	"io"
	"net/http"

	"github.com/corpix/formats"
	"github.com/corpix/loggers"
	"github.com/cryptounicorns/queues"
	"github.com/cryptounicorns/queues/consumer"
	"github.com/cryptounicorns/queues/result"
	"github.com/cryptounicorns/websocket"

	"github.com/cryptounicorns/platypus/http/handlers/routers"
	"github.com/cryptounicorns/platypus/http/handlers/routers/router"
	"github.com/cryptounicorns/platypus/writerpool"
)

type Stream struct {
	*websocket.UpgradeHandler

	config Config
	log    loggers.Logger
	done   chan struct{}

	websocketFormat formats.Format
	consumerFormat  formats.Format

	Conns    *writerpool.WriterPool
	Router   routers.Router
	Queue    queues.Queue
	Consumer consumer.Consumer
	events   <-chan result.Result
}

func (s *Stream) websocketWorker() {
	for {
		select {
		case <-s.done:
			return
		case r := <-s.events:
			// FIXME: I don't like this error handling
			var (
				v   interface{}
				buf []byte
				err error
			)
			if r.Err != nil {
				// XXX: Consumer always closes after error, so return here.
				s.log.Error(r.Err)
				return
			}

			if s.websocketFormat.Name() != s.consumerFormat.Name() {
				err = s.websocketFormat.Unmarshal(r.Value, v)
				if err != nil {
					s.log.Error(r.Err)
					continue
				}

				buf, err = s.consumerFormat.Marshal(v)
				if err != nil {
					s.log.Error(r.Err)
					continue
				}
			} else {
				buf = r.Value
			}

			_, err = s.Router.Write(buf)
			if err != nil {
				s.log.Error(err)
				continue
			}
		}
	}
}

func (s *Stream) ServeWebsocket(c io.WriteCloser, r *http.Request) {
	s.log.Print("websocket handler!")
	s.Conns.Add(c)
}

func (s *Stream) Close() error {
	var (
		err error
	)

	close(s.done)

	err = s.Consumer.Close()
	if err != nil {
		s.Queue.Close()
		return err
	}
	err = s.Queue.Close()
	if err != nil {
		return err
	}

	return nil
}

func New(c Config, l loggers.Logger) (*Stream, error) {
	var (
		writerPool      = writerpool.New()
		websocketFormat formats.Format
		consumerFormat  formats.Format
		r               routers.Router
		queue           queues.Queue
		consumer        consumer.Consumer
		events          <-chan result.Result
		stream          *Stream
		err             error
	)

	websocketFormat, err = formats.New(c.Format)
	if err != nil {
		return nil, err
	}

	consumerFormat, err = formats.New(c.Consumer.Format)
	if err != nil {
		return nil, err
	}

	r, err = routers.New(
		c.Router,
		writerPool,
		router.WebsocketWriter,
		router.WriterPoolCleanerErrorHandler(writerPool, l),
		l,
	)
	if err != nil {
		return nil, err
	}

	queue, err = queues.New(c.Consumer.Queue, l)
	if err != nil {
		return nil, err
	}

	consumer, err = queue.Consumer()
	if err != nil {
		queue.Close()
		return nil, err
	}

	events, err = consumer.Consume()
	if err != nil {
		consumer.Close()
		queue.Close()
		return nil, err
	}

	stream = &Stream{
		config: c,
		log:    l,

		websocketFormat: websocketFormat,
		consumerFormat:  consumerFormat,

		Conns:    writerPool,
		Router:   r,
		Queue:    queue,
		Consumer: consumer,
		events:   events,
	}
	stream.UpgradeHandler = websocket.NewUpgradeHandler(stream, l)

	go stream.websocketWorker()

	return stream, nil
}
