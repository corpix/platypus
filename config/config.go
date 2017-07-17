package config

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/corpix/formats"
	"github.com/corpix/queues"
	"github.com/corpix/queues/queue/nsq"
	"github.com/jinzhu/copier"

	"github.com/cryptounicorns/market-fetcher-http/feeds"
	"github.com/cryptounicorns/market-fetcher-http/http"
	"github.com/cryptounicorns/market-fetcher-http/logger"
)

var (
	// TickerFeedConfig represents default ticker feed config.
	TickerFeedConfig = queues.Config{
		Type: queues.NsqQueueType,
		Nsq: nsq.Config{
			Addr:     "127.0.0.1:4150",
			Topic:    "ticker",
			Channel:  "market-fetcher-http",
			LogLevel: nsq.LogLevelInfo,
		},
	}

	// FeedsConfig represents default feeds config.
	FeedsConfig = feeds.Config{
		Format: "json",
		Ticker: TickerFeedConfig,
	}

	// LoggerConfig represents default logger config.
	LoggerConfig = logger.Config{
		Level: "info",
	}

	// HTTPConfig represents default http server config.
	HTTPConfig = http.Config{
		Addr: ":8080",
	}

	// Default represents default application config.
	Default = Config{
		Logger: LoggerConfig,
		Feeds:  FeedsConfig,
		HTTP:   HTTPConfig,
	}
)

// Config represents application configuration structure.
type Config struct {
	Logger logger.Config
	Feeds  feeds.Config
	HTTP   http.Config
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
