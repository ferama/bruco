package source

// MessageHandler is a type for handler callback function. The handler
// will be invoked each time a source gets a message
// type MessageHandler func(msg *Message, resolve chan error)
type MessageHandler func(msg *Message, resolve chan error)

// Source interface that needs to be implemented from each source
type Source interface {
	// SetMessageHandler sets the callback function that will be invoked on each
	// message received from the kafka source
	SetMessageHandler(handler MessageHandler)
}

type SourceBase struct {
	MessageHandler MessageHandler
}

// SetMessageHandler sets the callback function that will be invoked on each
// message received
func (s *SourceBase) SetMessageHandler(handler MessageHandler) {
	s.MessageHandler = handler
}
