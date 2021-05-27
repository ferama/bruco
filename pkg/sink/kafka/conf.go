package kafkasink

import "github.com/ferama/bruco/pkg/sink"

// KafkaSinkConf ...
type KafkaSinkConf struct {
	sink.SinkConfCommon `yaml:",inline"`

	Brokers     []string `yaml:"brokers"`
	Topic       string   `yaml:"topic"`
	Partitioner string   `yaml:"partitioner"` // hash, random, manual
	Partition   string   `yaml:"partition"`   // -1
}
