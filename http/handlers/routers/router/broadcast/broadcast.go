package broadcast

import (
	"context"
	"io"
	"sync"

	"github.com/corpix/loggers"
	"github.com/corpix/pool"

	"github.com/cryptounicorns/platypus/http/handlers/routers/router"
)

type Broadcast struct {
	log          loggers.Logger
	ErrorHandler router.ErrorHandler

	*pool.Pool
	router.Writers
	router.Writer
	Config
}

func (b *Broadcast) worker(buf []byte, wg *sync.WaitGroup, c io.Writer, cancel context.CancelFunc) pool.Executor {
	return func(ctx context.Context) {
		select {
		case <-ctx.Done():
		default:
			err := b.Writer(c, buf)
			if err != nil {
				b.ErrorHandler(c, err)
			}
			cancel()
		}
		wg.Done()
	}
}

func (b *Broadcast) iterator(buf []byte, wg *sync.WaitGroup) func(io.Writer) {
	return func(c io.Writer) {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			b.Config.WriteTimeout,
		)

		wg.Add(1)
		b.Pool.Feed <- pool.NewWork(
			ctx,
			b.worker(buf, wg, c, cancel),
		)
	}
}

// Write writes to a pool of writers.
// Assumes every Writer in pool is thread-safe.
func (b *Broadcast) Write(buf []byte) (int, error) {
	var (
		wg = &sync.WaitGroup{}
	)

	b.Writers.Iter(b.iterator(buf, wg))

	wg.Wait()
	return len(buf), nil
}

func (b *Broadcast) Close() error {
	b.Pool.Close()

	return nil
}

func New(c Config, ws router.Writers, w router.Writer, e router.ErrorHandler, l loggers.Logger) (*Broadcast, error) {
	return &Broadcast{
		log:          l,
		Pool:         pool.NewFromConfig(c.Pool),
		Writers:      ws,
		Writer:       w,
		ErrorHandler: e,
		Config:       c,
	}, nil
}
