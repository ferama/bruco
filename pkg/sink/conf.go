package sink

// SinkConf all sinks should implement this interfaces
type SinkConf interface {
	GetKind() string
}

type SinkConfCommon struct {
	Kind string `json:"kind" yaml:"kind"`
}

func (s *SinkConfCommon) GetKind() string {
	return s.Kind
}
