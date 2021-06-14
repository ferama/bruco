package natssink

import "github.com/ferama/bruco/pkg/sink"

type NatsSinkConf struct {
	sink.SinkConfCommon `json:",inline" yaml:",inline"`

	ServerUrl string `json:"serverUrl" yaml:"serverUrl"`
	Subject   string `json:"subject" yaml:"subject"`
}
