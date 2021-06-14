package source

// SourceConf source configuration interface
type SourceConf interface {
	// IsFireAndForget defines if the function that
	// handles the source message should be async or not
	GetKind() string
	IsFireAndForget() bool
}

type SourceConfCommon struct {
	Kind          string `json:"kind" yaml:"kind"`
	FireAndForget bool   `json:"fireAndForget" yaml:"fireAndForget"`
}

func (s *SourceConfCommon) IsFireAndForget() bool {
	return s.FireAndForget
}

func (s *SourceConfCommon) GetKind() string {
	return s.Kind
}
