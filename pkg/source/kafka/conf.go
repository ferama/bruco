package kafkasource

type KafkaSourceConf struct {
	AsyncHandler    bool     `yaml:"asyncHandler"`
	BalanceStrategy string   `yaml:"balanceStrategy"`
	Brokers         []string `yaml:"brokers"`
	Topics          []string `yaml:"topics"`
	ConsumerGroup   string   `yaml:"consumerGroup"`
}

func (c *KafkaSourceConf) IsAsyncHandler() bool {
	return c.AsyncHandler
}
