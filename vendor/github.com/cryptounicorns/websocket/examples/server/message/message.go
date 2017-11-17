package message

import (
	"time"
)

type Message struct {
	Text string
	Date time.Time
}

func New(text string) *Message {
	return &Message{
		Text: text,
		Date: time.Now(),
	}
}
