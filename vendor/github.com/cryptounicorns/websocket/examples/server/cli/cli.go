package cli

import (
	"flag"
	"net/http"
	"net/url"
	"time"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/logrus"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/gorilla/mux"
	logrusLogger "github.com/sirupsen/logrus"

	"github.com/cryptounicorns/websocket"
	"github.com/cryptounicorns/websocket/examples/server/handler"
	"github.com/cryptounicorns/websocket/writer"
)

var (
	Addr  string
	Debug bool
)

func init() {
	flag.StringVar(
		&Addr,
		"addr",
		"ws://127.0.0.1:9999",
		"Websocket server address",
	)
	flag.BoolVar(
		&Debug,
		"debug",
		false,
		"Debug mode",
	)

	flag.Parse()
}

func newRouter(l loggers.Logger) *mux.Router {
	var (
		r = mux.NewRouter()
	)

	r.Handle(
		"/",
		websocket.NewUpgradeHandler(
			websocket.NewWriterWrapHandler(
				handler.New(l),
				writer.NewServerText,
			),
			prefixwrapper.New("/: ", l),
		),
	)

	return r
}

func newLogger() loggers.Logger {
	var (
		logger = logrusLogger.New()
	)

	if Debug {
		logger.Level = logrusLogger.DebugLevel
	}

	return logrus.New(logger)
}

func ListenAndServe(addr string, handler http.Handler, log loggers.Logger) error {
	var (
		l = prefixwrapper.New("ListenAndServe: ", log)

		u          *url.URL
		listenAddr string
		err        error
	)

	u, err = url.Parse(addr)
	if err != nil {
		return err
	}

	// FIXME: ws and wss has no difference at this time.
	// But this difference is very important.
	listenAddr = u.Hostname() + ":" + u.Port()

	for {
		l.Printf(
			"Starting server on %s",
			listenAddr,
		)

		err = http.ListenAndServe(
			listenAddr,
			handler,
		)
		if err != nil {
			l.Error(err)
			time.Sleep(1 * time.Second)
		}
	}
}

func Run() error {
	var (
		l = newLogger()
	)

	l.Debug("Running in debug mode")

	return ListenAndServe(
		Addr,
		newRouter(l),
		l,
	)
}
