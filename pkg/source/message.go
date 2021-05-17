package source

import "time"

type MessageHandler func(msg *Message)

type Message struct {
	Value     []byte
	Timestamp time.Time
}
