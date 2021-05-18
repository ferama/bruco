package source

type SourceConf interface {
	IsAsyncHandler() bool
}
