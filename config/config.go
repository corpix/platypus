package config

import (
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/corpix/formats"
	"github.com/corpix/pool"
	"github.com/cryptounicorns/queues"
	"github.com/cryptounicorns/queues/queue/nsq"
	"github.com/imdario/mergo"

	"github.com/cryptounicorns/platypus/http"
	"github.com/cryptounicorns/platypus/http/handlers"
	"github.com/cryptounicorns/platypus/http/handlers/cache"
	"github.com/cryptounicorns/platypus/http/handlers/consumer"
	"github.com/cryptounicorns/platypus/http/handlers/handler/latest"
	"github.com/cryptounicorns/platypus/http/handlers/handler/stream"
	"github.com/cryptounicorns/platypus/http/handlers/routers"
	"github.com/cryptounicorns/platypus/http/handlers/routers/router/broadcast"
	"github.com/cryptounicorns/platypus/http/handlers/stores"
	"github.com/cryptounicorns/platypus/http/handlers/stores/store/memoryttl"
	"github.com/cryptounicorns/platypus/logger"
)

var (
	// LoggerConfig represents default logger config.
	LoggerConfig = logger.Config{
		Level: "info",
	}

	// HTTPConfig represents default http server config.
	HTTPConfig = http.Config{
		Addr: ":8080",
		Handlers: handlers.Configs{
			{
				Path:   "/api/v1/tickers/stream",
				Method: "get",
				Type:   handlers.StreamType,
				Stream: stream.Config{
					Format: formats.JSON,
					Consumer: consumer.Config{
						Format: formats.JSON,
						Queue: queues.Config{
							Type: queues.NsqQueueType,
							Nsq: nsq.Config{
								Addr:     "127.0.0.1:4150",
								Topic:    "ticker",
								Channel:  "platypus-stream",
								LogLevel: nsq.LogLevelInfo,
							},
						},
					},
					// FIXME: Rename to Router?
					Transmitter: routers.Config{
						Type: routers.BroadcastTransmitterType,
						Broadcast: broadcast.Config{
							WriteTimeout: 10 * time.Second,
							Pool: pool.Config{
								Workers:   128,
								QueueSize: 1024,
							},
						},
					},
				},
			},

			{
				Path:   "/api/v1/tickers",
				Method: "get",
				Type:   handlers.LatestType,
				Latest: latest.Config{
					Format: formats.JSON,
					Consumer: consumer.Config{
						Format: formats.JSON,
						Queue: queues.Config{
							Type: queues.NsqQueueType,
							Nsq: nsq.Config{
								Addr:     "127.0.0.1:4150",
								Topic:    "ticker",
								Channel:  "platypus-latest",
								LogLevel: nsq.LogLevelInfo,
							},
						},
					},
					Cache: cache.Config{
						Key: "",
						Store: stores.Config{
							Type: stores.MemoryTTLStoreType,
							MemoryTTL: memoryttl.Config{
								TTL:        30 * time.Second,
								Resolution: 1 * time.Second,
							},
						},
					},
				},
			},
		},
	}

	// Default represents default application config.
	Default = Config{
		Logger: LoggerConfig,
		HTTP:   HTTPConfig,
	}
)

// Config represents application configuration structure.
type Config struct {
	Logger logger.Config
	HTTP   http.Config
}

// FromReader returns parsed config data in some `f` from reader `r`.
// It copies `Default` into the target structure before unmarshaling
// config, so it will have default values.
func FromReader(f formats.Format, r io.Reader) (Config, error) {
	var (
		c   Config
		buf []byte
		err error
	)

	buf, err = ioutil.ReadAll(r)
	if err != nil {
		return c, err
	}

	err = mergo.Merge(&c, Default)
	if err != nil {
		return c, err
	}

	err = f.Unmarshal(buf, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

// FromFile returns parsed config data from file in `path`.
func FromFile(path string) (Config, error) {
	var (
		c   Config
		f   formats.Format
		r   io.ReadWriteCloser
		err error
	)
	f, err = formats.NewFromPath(path)
	if err != nil {
		return c, err
	}

	r, err = os.Open(path)
	if err != nil {
		return c, err
	}
	defer r.Close()

	c, err = FromReader(f, r)
	if err != nil {
		return c, err
	}

	return c, nil
}
