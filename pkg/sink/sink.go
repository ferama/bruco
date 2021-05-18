package sink

type Sink interface {
	Publish(key string, msg []byte)
}
