package kafkasource

type KafkaSourceConf struct {
	FireAndForget bool `yaml:"fireAndForget"`

	BalanceStrategy string   `yaml:"balanceStrategy"`
	Brokers         []string `yaml:"brokers"`
	Topics          []string `yaml:"topics"`
	ConsumerGroup   string   `yaml:"consumerGroup"`
	Offset          string   `yaml:"offset"`
}

func (c *KafkaSourceConf) IsFireAndForget() bool {
	return c.FireAndForget
}
