package natssource

import "github.com/ferama/bruco/pkg/source"

type NatsSourceConf struct {
	source.SourceConfCommon `json:",inline" yaml:",inline"`

	ServerUrl  string `json:"serverUrl" yaml:"serverUrl"`
	QueueGroup string `json:"queueGroup" yaml:"queueGroup"`
	Subject    string `json:"subject" yaml:"subject"`
}

func (s *NatsSourceConf) IsFireAndForget() bool {
	return true
}
