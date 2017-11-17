package consumer

import (
	"io"
	"io/ioutil"
)

type Consumer struct {
	reader io.Reader

	messages chan Result
	done     chan struct{}
}

func (c *Consumer) Consume() <-chan Result {
	go c.consumeLoop()

	return c.messages
}

func (c *Consumer) consumeLoop() {
	var (
		buf []byte
		err error
	)

	for {
		buf, err = ioutil.ReadAll(c.reader)
		// FIXME: Send me an empty buffer and I will
		// return you EOF. You can't know where was a
		// network connection reset.
		if err == nil && len(buf) == 0 {
			err = io.EOF
		}

		select {
		case <-c.done:
			return
		default:
			c.messages <- Result{
				Value: buf,
				Err:   err,
			}
			if err != nil {
				return
			}
		}
	}
}

func (c *Consumer) Close() error {
	close(c.done)
	// Not closing, it will be GC'ed.
	// close(c.messages)

	return nil
}

func New(r io.Reader) *Consumer {
	return &Consumer{
		reader:   r,
		messages: make(chan Result),
		done:     make(chan struct{}),
	}
}
