package config

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/corpix/formats"
	"github.com/jinzhu/copier"

	"github.com/cryptounicorns/market-fetcher-http/feeds"
	"github.com/cryptounicorns/market-fetcher-http/feeds/feed/nsq"
	"github.com/cryptounicorns/market-fetcher-http/logger"
)

var (
	// NsqConfig represents default nsq config.
	NsqConfig = nsq.WrapConfig(
		nsq.Config{
			Addr:    "127.0.0.1:4150",
			Topic:   "ticker",
			Channel: "market-fetcher-http",
		},
	)

	// NsqFeed represents default nsq feed config.
	NsqFeed = feeds.Config{
		Type: feeds.NsqFeedType,
		Nsq:  NsqConfig,
	}

	// LoggerConfig represents default configuration for a logger.
	LoggerConfig = logger.Config{
		Level: "info",
	}

	// Default represents default application configuration.
	Default = Config{
		Logger: LoggerConfig,
		Feed:   NsqFeed,
	}
)

// Config represents application configuration structure.
type Config struct {
	Logger logger.Config
	Feed   feeds.Config
}

// FromReader fills Config structure `c` passed by reference with
// parsed config data in some `f` from reader `r`.
// It copies `Default` into the target structure before unmarshaling
// config, so it will have default values.
func FromReader(f formats.Format, r io.Reader, c *Config) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	err = copier.Copy(c, Default)
	if err != nil {
		return err
	}

	return f.Unmarshal(data, c)
}

// FromFile fills Config structure `c` passed by reference with
// parsed config data from file in `path`.
func FromFile(path string, c *Config) error {
	f, err := formats.NewFromPath(path)
	if err != nil {
		return err
	}

	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	return FromReader(f, r, c)
}
