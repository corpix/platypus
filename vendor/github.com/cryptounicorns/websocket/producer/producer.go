package producer

import (
	"io"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
)

type Producer struct {
	writer io.Writer
	log    loggers.Logger
}

func (c *Producer) Produce(buf []byte) error {
	var (
		err error
	)

	_, err = c.writer.Write(buf)
	return err
}

func (c *Producer) Close() error {
	return nil
}

func NewProducer(w io.Writer, l loggers.Logger) *Producer {
	return &Producer{
		writer: w,
		log: prefixwrapper.New(
			"Producer: ",
			l,
		),
	}
}
