package httpsource

import "github.com/ferama/bruco/pkg/source"

type HttpSourceConf struct {
	source.SourceConfCommon `json:",inline" yaml:",inline"`

	Port int `json:"port" yaml:"port"`
	// If set the http source will not return the processor response
	// to the caller
	IgnoreProcessorResponse bool `json:"ignoreProcessorResponse" yaml:"ignoreProcessorResponse"`
}

func (s *HttpSourceConf) IsFireAndForget() bool {
	return true
}
