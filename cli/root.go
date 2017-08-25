package cli

import (
	"time"

	"github.com/urfave/cli"

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
		f   *feeds.Feeds
		s   *http.Server
		err error
	)

	f, err = feeds.New(
		Config.Feeds,
		log,
	)
	if err != nil {
		return err
	}
	defer f.Close()

	s, err = http.New(
		Config.HTTP,
		f,
		log,
	)
	if err != nil {
		return err
	}

	for {
		err = s.Serve()
		if err != nil {
			log.Error(err)
			time.Sleep(1 * time.Second)
		}
	}
}
