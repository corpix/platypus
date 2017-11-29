package cli

import (
	"context"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/cli"

	"github.com/cryptounicorns/platypus/handlers"
	"github.com/cryptounicorns/platypus/http"
)

var (
	// RootCommands is a list of subcommands for the application.
	RootCommands = []cli.Command{}

	// RootFlags is a list of flags for the application.
	RootFlags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "application configuration file",
			EnvVar: "CONFIG",
			Value:  "config.json",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "add this flag to enable debug mode",
		},
	}
)

// RootAction is executing when program called without any subcommand.
func RootAction(c *cli.Context) error {
	var (
		r   = mux.NewRouter()
		s   *http.Server
		h   handlers.Handler
		err error
	)

	for _, handler := range Config.Handlers {
		h, err = handlers.New(handler, log)
		if err != nil {
			return err
		}

		go func(h handlers.Handler) {
			var (
				ctx    context.Context
				cancel context.CancelFunc
				err    error
			)

			for {
				ctx, cancel = context.WithCancel(context.Background())
				err = h.Run(ctx)
				cancel()
				if err != nil {
					log.Error(err)
					// FIXME: Implement exponential grow component for timer
					// because in case of error where will be too much
					// shit in logs.
					time.Sleep(1 * time.Second)
				} else {
					log.Print("Handler finished without error, looks like we want to exit, breaking loop...")
					break
				}
			}
		}(h)

		r.
			Methods(handler.Method).
			Path(handler.Path).
			Handler(h)
	}

	s = http.New(
		Config.HTTP,
		r,
		log,
	)

	for {
		err = s.Serve()
		if err != nil {
			log.Error(err)
			time.Sleep(1 * time.Second)
		}
	}
}
