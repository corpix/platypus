package cli

import (
	"time"

	"github.com/corpix/formats"
	"github.com/urfave/cli"

	"github.com/cryptounicorns/market-fetcher-http/datasources"
	"github.com/cryptounicorns/market-fetcher-http/feeds"
	"github.com/cryptounicorns/market-fetcher-http/http"
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
		fmts formats.Format
		f    *feeds.Feeds
		d    *datasources.Datasources
		s    *http.Server
		err  error
	)

	fmts, err = formats.New(Config.Feeds.Format)
	if err != nil {
		return err
	}

	f, err = feeds.NewFromConfig(
		Config.Feeds,
		log,
	)
	if err != nil {
		return err
	}

	d, err = datasources.New(f, fmts, log)
	if err != nil {
		return err
	}

	s = http.New(
		Config.HTTP,
		d,
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
