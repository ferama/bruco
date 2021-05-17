package source

type Source interface {
	SetMessageHandler(handler MessageHandler)
}
