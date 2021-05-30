package httpsource

import "github.com/ferama/bruco/pkg/source"

type HttpSourceConf struct {
	source.SourceConfCommon `yaml:",inline"`

	Port int `yaml:"port"`
	// If set the http source will not return the processor response
	// to the caller
	IgnoreProcessorResponse bool `yaml:"ignoreProcessorResponse"`
}

func (s *HttpSourceConf) IsFireAndForget() bool {
	return true
}
