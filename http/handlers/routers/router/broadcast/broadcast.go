package broadcast

import (
	"context"
	"io"
	"sync"

	"github.com/corpix/loggers"
	"github.com/corpix/pool"

	"github.com/cryptounicorns/platypus/http/handlers/routers/errors"
	"github.com/cryptounicorns/platypus/http/handlers/routers/writer"
)

type Broadcast struct {
	ErrorHandler errors.Handler
	Pool         *pool.Pool
	Iterator     writer.Iterator
	Reader       io.Reader
	Config       Config
	log          loggers.Logger
}

func (b *Broadcast) worker(buf []byte, wg *sync.WaitGroup, c io.Writer, cancel context.CancelFunc) pool.Executor {
	return func(ctx context.Context) {
		defer wg.Done()

		var (
			err error
		)

		select {
		case <-ctx.Done():
			deadline, _ := ctx.Deadline()
			b.log.Error("Canceled after ", deadline)
		default:
			defer cancel()

			_, err = c.Write(buf)
			if err != nil {
				b.ErrorHandler(c, err)
				return
			}
		}
	}
}

// Write writes to a pool of writers.
// Assumes every Writer in pool is thread-safe.
func (b *Broadcast) Write(buf []byte) (int, error) {
	var (
		wg = &sync.WaitGroup{}
	)

	b.Iterator.Iter(
		func(c io.Writer) {
			var (
				ctx     context.Context
				cancel  context.CancelFunc
				timeout = b.Config.WriteTimeout
			)

			ctx, cancel = context.WithTimeout(
				context.Background(),
				timeout.Duration(),
			)

			wg.Add(1)
			b.Pool.Feed <- pool.NewWork(
				ctx,
				b.worker(buf, wg, c, cancel),
			)
		},
	)

	wg.Wait()
	return len(buf), nil
}

func (b *Broadcast) Close() error {
	b.Pool.Close()
	return nil
}

func New(c Config, w writer.Iterator, e errors.Handler, l loggers.Logger) (*Broadcast, error) {
	return &Broadcast{
		Pool:         pool.NewFromConfig(c.Pool),
		Iterator:     w,
		ErrorHandler: e,
		Config:       c,
		log:          l,
	}, nil
}
