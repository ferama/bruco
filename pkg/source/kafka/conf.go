package kafka

type KafkaConf struct {
	BalanceStrategy string
	Brokers         []string
	Topics          []string
	ConsumerGroup   string
}
