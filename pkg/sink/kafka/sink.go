package kafkasink

import (
	"log"

	"github.com/Shopify/sarama"
)

type KafkaSink struct {
	topic    string
	producer sarama.SyncProducer
}

func NewKafkaSink(kconf *KafkaSinkConf) *KafkaSink {

	// Example here: https://github.com/Shopify/sarama/blob/master/tools/kafka-console-producer/kafka-console-producer.go
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewSyncProducer(kconf.Brokers, config)
	if err != nil {
		log.Fatalln(err)
	}

	sink := &KafkaSink{
		topic:    kconf.Topic,
		producer: producer,
	}
	return sink
}

func (s *KafkaSink) Publish(key string, msg []byte) {
	// log.Printf("Publishing: %s %s", key, string(msg))
	message := &sarama.ProducerMessage{Topic: s.topic, Partition: int32(-1)}
	message.Value = sarama.ByteEncoder(msg)
	s.producer.SendMessage(message)
}
