package config

import (
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/corpix/formats"
	"github.com/corpix/pool"
	"github.com/corpix/queues"
	"github.com/corpix/queues/queue/nsq"
	"github.com/imdario/mergo"

	"github.com/cryptounicorns/platypus/consumer"
	"github.com/cryptounicorns/platypus/http"
	"github.com/cryptounicorns/platypus/http/endpoints"
	"github.com/cryptounicorns/platypus/logger"
	"github.com/cryptounicorns/platypus/stores"
	"github.com/cryptounicorns/platypus/stores/store/memoryttl"
	"github.com/cryptounicorns/platypus/transmitters"
	"github.com/cryptounicorns/platypus/transmitters/transmitter/broadcast"
)

var (
	// LoggerConfig represents default logger config.
	LoggerConfig = logger.Config{
		Level: "info",
	}

	// HTTPConfig represents default http server config.
	HTTPConfig = http.Config{
		Addr: ":8080",
		Endpoints: endpoints.Config{
			{
				Path:   "/api/v1/tickers/stream",
				Method: "get",

				Queue: queues.Config{
					Type: queues.NsqQueueType,
					Nsq: nsq.Config{
						Addr:     "127.0.0.1:4150",
						Topic:    "ticker",
						Channel:  "platypus",
						LogLevel: nsq.LogLevelInfo,
					},
				},
				Consumer: consumer.Config{
					Format: "json",
				},
				Store: stores.Config{
					Type: stores.MemoryTTLStoreType,
					MemoryTTL: memoryttl.Config{
						TTL:        2 * time.Second,
						Resolution: 1 * time.Second,
					},
				},
				Transmitter: transmitters.Config{
					Type: transmitters.BroadcastTransmitterType,
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

// FromReader fills Config structure `c` passed by reference with
// parsed config data in some `f` from reader `r`.
// It copies `Default` into the target structure before unmarshaling
// config, so it will have default values.
func FromReader(f formats.Format, r io.Reader, c *Config) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	err = mergo.Merge(c, Default)
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
