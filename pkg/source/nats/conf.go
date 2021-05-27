package natssource

import "github.com/ferama/bruco/pkg/source"

type NatsSourceConf struct {
	source.SourceConfCommon `yaml:",inline"`

	ServerUrl  string `yaml:"serverUrl"`
	QueueGroup string `yaml:"queueGroup"`
	Subject    string `yaml:"subject"`
}

func (s *NatsSourceConf) IsFireAndForget() bool {
	return true
}
