package source

// SourceConf source configuration interface
type SourceConf interface {
	// IsAsyncHandler defines if the function that
	// handles the source message should be async or not
	IsAsyncHandler() bool
}
