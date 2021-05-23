package kafkasource

import "github.com/ferama/bruco/pkg/source"

// KafkaSourceConf holds the kafka source configuration
type KafkaSourceConf struct {
	source.SourceConfBase

	FetchDefaultBytes string   `yaml:"fetchDefaultBytes"`
	ChannelBufferSize string   `yaml:"channelBufferSize"`
	BalanceStrategy   string   `yaml:"balanceStrategy"`
	RebalanceTimeout  int      `yaml:"rebalanceTimeout"`
	Brokers           []string `yaml:"brokers"`
	Topics            []string `yaml:"topics"`
	ConsumerGroup     string   `yaml:"consumerGroup"`
	Offset            string   `yaml:"offset"`
}
