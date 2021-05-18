package kafkasink

// KafkaSinkConf ...
type KafkaSinkConf struct {
	Brokers     []string `yaml:"brokers"`
	Topic       string   `yaml:"topic"`
	Partitioner string   `yaml:"partitioner"` // hash, random, manual
	Partition   string   `yaml:"partition"`   // -1
}
