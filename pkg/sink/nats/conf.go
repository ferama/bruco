package natssink

import "github.com/ferama/bruco/pkg/sink"

type NatsSinkConf struct {
	sink.SinkConfCommon `yaml:",inline"`

	ServerUrl string `yaml:"serverUrl"`
	Subject   string `yaml:"subject"`
}
