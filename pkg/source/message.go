package source

import "time"

// Message is the source message struct
type Message struct {
	Value     []byte
	Timestamp time.Time
}
