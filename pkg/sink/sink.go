package sink

// Sink the sink interface that needs to be implemented from each
// sink
type Sink interface {
	// Publish publishes a message to the sink
	Publish(msg *Message)
}
