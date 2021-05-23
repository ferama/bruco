package source

// SourceConf source configuration interface
type SourceConf interface {
	// IsFireAndForget defines if the function that
	// handles the source message should be async or not
	IsFireAndForget() bool
}

type SourceConfBase struct {
	FireAndForget bool `yaml:"fireAndForget"`
}

func (s *SourceConfBase) IsFireAndForget() bool {
	return s.FireAndForget
}
