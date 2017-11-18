package stream

import (
	"io"
	"net/http"

	"github.com/corpix/formats"
	"github.com/corpix/loggers"
	"github.com/cryptounicorns/websocket"
	"github.com/cryptounicorns/websocket/writer"

	"github.com/cryptounicorns/platypus/http/handlers/consumer"
	"github.com/cryptounicorns/platypus/http/handlers/routers"
	"github.com/cryptounicorns/platypus/iopool"
)

type Stream struct {
	*websocket.UpgradeHandler
	config          Config
	log             loggers.Logger
	done            chan struct{}
	websocketFormat formats.Format
	Connections     *iopool.Writer
	Router          routers.Router
	Consumer        *consumer.Consumer
}

func (s *Stream) consumerWorker() {
	for {
		select {
		case <-s.done:
			return
		case r := <-s.Consumer.Stream():
			// FIXME: I don't like this error handling
			var (
				v   = r.Value
				buf []byte
				err error
			)
			if r.Err != nil {
				// XXX: Consumer always closes after error, so return here.
				s.log.Error(r.Err)
				return
			}

			buf, err = s.websocketFormat.Marshal(v)
			if err != nil {
				s.log.Error(r.Err)
				continue
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
	s.Connections.Add(writer.NewServerText(c))
}

func (s *Stream) Close() error {
	var (
		err error
	)

	close(s.done)

	s.Router.Close()
	err = s.Consumer.Close()
	if err != nil {
		return err
	}

	return nil
}

func New(c Config, l loggers.Logger) (*Stream, error) {
	var (
		wp     = iopool.NewWriter()
		wf     formats.Format
		r      routers.Router
		cr     *consumer.Consumer
		stream *Stream
		err    error
	)

	wf, err = formats.New(c.Format)
	if err != nil {
		return nil, err
	}

	r, err = routers.New(
		c.Router,
		wp,
		routers.NewWriterPoolCleaner(wp, l),
		l,
	)
	if err != nil {
		return nil, err
	}

	cr, err = consumer.New(
		c.Consumer,
		l,
	)
	if err != nil {
		return nil, err
	}

	stream = &Stream{
		config:          c,
		log:             l,
		done:            make(chan struct{}),
		websocketFormat: wf,
		Connections:     wp,
		Router:          r,
		Consumer:        cr,
	}

	stream.UpgradeHandler = websocket.NewUpgradeHandler(
		stream,
		l,
	)

	go stream.consumerWorker()

	return stream, nil
}
