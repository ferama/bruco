package sink

type Sink interface {
	Publish(msg *Message)
}
