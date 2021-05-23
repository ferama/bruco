package kafkasource

// KafkaSourceConf holds the kafka source configuration
type KafkaSourceConf struct {
	FireAndForget bool `yaml:"fireAndForget"`

	FetchDefaultBytes string   `yaml:"fetchDefaultBytes"`
	ChannelBufferSize string   `yaml:"channelBufferSize"`
	BalanceStrategy   string   `yaml:"balanceStrategy"`
	RebalanceTimeout  int      `yaml:"rebalanceTimeout"`
	Brokers           []string `yaml:"brokers"`
	Topics            []string `yaml:"topics"`
	ConsumerGroup     string   `yaml:"consumerGroup"`
	Offset            string   `yaml:"offset"`
}

func (c *KafkaSourceConf) IsFireAndForget() bool {
	return c.FireAndForget
}
