package kafkasink

type KafkaSinkConf struct {
	Brokers []string `yaml:"brokers"`
	Topics  []string `yaml:"topics"`
}
