package kafkasink

import "github.com/ferama/bruco/pkg/sink"

// KafkaSinkConf ...
type KafkaSinkConf struct {
	sink.SinkConfCommon `json:",inline" yaml:",inline"`

	Brokers     []string `json:"brokers" yaml:"brokers"`
	Topic       string   `json:"topic" yaml:"topic"`
	Partitioner string   `json:"partitioner" yaml:"partitioner"` // hash, random, manual
	Partition   string   `json:"partition" yaml:"partition"`     // -1
}
