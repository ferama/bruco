package source

// SourceConf source configuration interface
type SourceConf interface {
	// IsFireAndForget defines if the function that
	// handles the source message should be async or not
	GetKind() string
	IsFireAndForget() bool
}

type SourceConfCommon struct {
	Kind          string `yaml:"kind"`
	FireAndForget bool   `yaml:"fireAndForget"`
}

func (s *SourceConfCommon) IsFireAndForget() bool {
	return s.FireAndForget
}

func (s *SourceConfCommon) GetKind() string {
	return s.Kind
}
