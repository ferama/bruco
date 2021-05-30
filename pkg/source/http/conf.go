package httpsource

import "github.com/ferama/bruco/pkg/source"

type HttpSourceConf struct {
	source.SourceConfCommon `yaml:",inline"`
}

func (s *HttpSourceConf) IsFireAndForget() bool {
	return true
}