package sink

type Sink interface {
	Publish(msg []byte)
}
