package kafkasink

type KafkaSinkConf struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}
