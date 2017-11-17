package cli

import (
	"context"
	"flag"
	"io"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/logrus"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	logrusLogger "github.com/sirupsen/logrus"

	"github.com/cryptounicorns/websocket/consumer"
)

var (
	Addr  string
	Limit uint
	Debug bool
)

func init() {
	flag.StringVar(
		&Addr,
		"addr",
		"ws://127.0.0.1:9999",
		"Websocket server address",
	)
	flag.UintVar(
		&Limit,
		"limit",
		5,
		"Number of messages to receive before exit(use 0 to disable limit)",
	)
	flag.BoolVar(
		&Debug,
		"debug",
		false,
		"Debug mode",
	)

	flag.Parse()
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

func Run() error {
	var (
		l = newLogger()

		r          io.ReadCloser
		messages   *consumer.Consumer
		messageNum uint
		err        error
	)

	r, _, _, err = ws.DefaultDialer.Dial(
		context.Background(),
		Addr,
	)
	if err != nil {
		return err
	}
	defer r.Close()

	messages = consumer.New(
		wsutil.NewReader(
			r,
			ws.StateClientSide,
		),
	)
	defer messages.Close()

	for m := range messages.Consume() {
		l.Printf(
			"message=%s",
			m,
		)

		if m.Err != nil {
			return m.Err
		}

		if Limit > 0 {
			messageNum++
			if messageNum >= Limit {
				break
			}
		}
	}

	return nil
}
