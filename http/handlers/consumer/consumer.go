package consumer

type Consumer interface {
	Consume() (<-chan Result, error)
	Close() error
}

type Consumers = []Consumer
