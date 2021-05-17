package kafkasink

import (
	"log"

	"github.com/Shopify/sarama"
)

type KafkaSink struct {
}

func NewKafkaSink(kconf *KafkaSinkConf) *KafkaSink {
	sink := &KafkaSink{}

	// Example here: https://github.com/Shopify/sarama/blob/master/tools/kafka-console-producer/kafka-console-producer.go
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	// message := &sarama.ProducerMessage{Topic: kconf.Topic, Partition: int32(-1)}

	return sink
}

func (s *KafkaSink) Publish(msg []byte) {
	log.Printf("Publishing: %s", string(msg))
}
