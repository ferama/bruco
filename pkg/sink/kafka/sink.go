package kafkasink

import (
	"log"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/ferama/bruco/pkg/sink"
)

// KafkaSink ...
type KafkaSink struct {
	producer sarama.SyncProducer

	topic string

	// resolved field
	partition int32
}

// NewKafkaSink creates a new sink object given the conf
func NewKafkaSink(kconf *KafkaSinkConf) *KafkaSink {
	sink := &KafkaSink{
		topic: kconf.Topic,
	}
	// Example here: https://github.com/Shopify/sarama/blob/master/tools/kafka-console-producer/kafka-console-producer.go
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sink.resolvePartitioner(kconf)

	producer, err := sarama.NewSyncProducer(kconf.Brokers, config)
	if err != nil {
		log.Fatalln(err)
	}

	sink.partition = sink.resolvePartition(kconf)
	sink.producer = producer

	log.Printf("[KAFKA-SINK] started sink brokers: %s, topic: %s", kconf.Brokers, kconf.Topic)
	return sink
}

func (s *KafkaSink) resolvePartition(cfg *KafkaSinkConf) int32 {
	if cfg.Partition == "" {
		return -1
	}
	i, err := strconv.Atoi(cfg.Partition)
	if err != nil {
		log.Fatalf("[KAFKA-SINK] invalid conf. Sink partition %s", cfg.Partition)
	}
	return int32(i)
}

func (s *KafkaSink) resolvePartitioner(cfg *KafkaSinkConf) func(string) sarama.Partitioner {
	switch lower := strings.ToLower(cfg.Partitioner); lower {
	case "":
		if s.partition >= 0 {
			return sarama.NewManualPartitioner
		} else {
			return sarama.NewHashPartitioner
		}
	case "hash":
		return sarama.NewHashPartitioner
	case "random":
		return sarama.NewRandomPartitioner
	case "manual":
		if s.partition == -1 {
			log.Fatalf("[KAFKA-SINK] partition is required while using manual partitioner")
		}
		return sarama.NewManualPartitioner
	default:
		log.Fatalf("[KAFKA-SINK] invalid partitioner %s", cfg.Partitioner)
		return nil
	}
}

// Publish send a message through the sink
func (s *KafkaSink) Publish(msg *sink.Message) error {
	// log.Printf("Publishing: %s %s", msg.Key, msg.Value)
	message := &sarama.ProducerMessage{
		Topic:     s.topic,
		Partition: s.partition,
	}
	if msg.Key != "" {
		message.Key = sarama.StringEncoder(msg.Key)
	}
	message.Value = sarama.ByteEncoder(msg.Value)
	_, _, err := s.producer.SendMessage(message)
	return err
}
