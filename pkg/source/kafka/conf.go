package kafka

type KafkaSourceConf struct {
	BalanceStrategy string   `yaml:"balanceStrategy"`
	Brokers         []string `yaml:"brokers"`
	Topics          []string `yaml:"topics"`
	ConsumerGroup   string   `yaml:"consumerGroup"`
	AsyncHandler    bool     `yaml:"asyncHandler"`
}
