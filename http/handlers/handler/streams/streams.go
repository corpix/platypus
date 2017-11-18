package streams

import (
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/corpix/formats"
	"github.com/corpix/loggers"
	"github.com/cryptounicorns/websocket"
	"github.com/cryptounicorns/websocket/writer"

	"github.com/cryptounicorns/platypus/http/handlers/consumer"
	"github.com/cryptounicorns/platypus/http/handlers/routers"
	"github.com/cryptounicorns/platypus/iopool"
)

type templateEventContext struct {
	JSON  []byte
	Value interface{}
}
type templateContext struct {
	Consumer consumer.Config
	Event    templateEventContext
}

type Streams struct {
	*websocket.UpgradeHandler
	config          Config
	log             loggers.Logger
	done            chan struct{}
	websocketFormat formats.Format
	Connections     *iopool.Writer
	Router          routers.Router
	Consumers       []*consumer.Consumer
}

func (s *Streams) consumerWorker(wrap *template.Template, consumer *consumer.Consumer) {
	var (
		buf []byte
		err error
	)

	for {
		select {
		case <-s.done:
			return
		case r := <-consumer.Stream():
			// FIXME: I don't like this error handling
			if r.Err != nil {
				// XXX: Consumer always closes after error, so return here.
				s.log.Error(r.Err)
				return
			}

			buf, err = s.websocketFormat.Marshal(r.Value)
			if err != nil {
				s.log.Error(r.Err)
				continue
			}

			err = wrap.Execute(
				s.Router,
				templateContext{
					Consumer: consumer.Meta.Config,
					Event: templateEventContext{
						JSON:  buf,
						Value: r.Value,
					},
				},
			)
			if err != nil {
				s.log.Error(err)
				continue
			}
		}
	}
}

func (s *Streams) consumersWorkers(wrap *template.Template) {
	for _, cr := range s.Consumers {
		go s.consumerWorker(wrap, cr)
	}
}

func (s *Streams) ServeWebsocket(c io.WriteCloser, r *http.Request) {
	s.Connections.Add(writer.NewServerText(c))
}

func (s *Streams) Close() error {
	var (
		err error
	)

	close(s.done)

	// FIXME: Refactor writer pool to writecloser
	// So we could close all connections.
	for _, cr := range s.Consumers {
		err = cr.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func New(c Config, l loggers.Logger) (*Streams, error) {
	var (
		wp  = iopool.NewWriter()
		wf  formats.Format
		cs  []*consumer.Consumer
		t   *template.Template
		r   routers.Router
		s   *Streams
		err error
	)

	wf, err = formats.New(c.Format)
	if err != nil {
		return nil, err
	}

	t, err = template.New("wrap").Parse(
		strings.TrimSpace(c.Wrap),
	)
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

	cs, err = consumer.NewConsumers(c.Consumers, l)
	if err != nil {
		return nil, err
	}

	s = &Streams{
		config:          c,
		log:             l,
		done:            make(chan struct{}),
		websocketFormat: wf,
		Connections:     wp,
		Router:          r,
		Consumers:       cs,
	}

	s.UpgradeHandler = websocket.NewUpgradeHandler(s, l)

	s.consumersWorkers(t)

	return s, nil
}
