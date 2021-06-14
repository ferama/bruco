package kafkasource

import "github.com/ferama/bruco/pkg/source"

// KafkaSourceConf holds the kafka source configuration
type KafkaSourceConf struct {
	source.SourceConfCommon `json:",inline" yaml:",inline"`

	FetchDefaultBytes string   `json:"fetchDefaultBytes" yaml:"fetchDefaultBytes"`
	ChannelBufferSize string   `json:"channelBufferSize" yaml:"channelBufferSize"`
	BalanceStrategy   string   `json:"balanceStrategy" yaml:"balanceStrategy"`
	RebalanceTimeout  int      `json:"rebalanceTimeout" yaml:"balanceStrategy"`
	Brokers           []string `json:"brokers" yaml:"brokers"`
	Topics            []string `json:"topics" yaml:"topics"`
	ConsumerGroup     string   `json:"consumerGroup" yaml:"consumerGroup"`
	Offset            string   `json:"offset" yaml:"offset"`
}
