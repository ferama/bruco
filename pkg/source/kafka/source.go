package kafkasource

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ferama/bruco/pkg/source"
)

type KafkaSource struct {
	messageHandler source.MessageHandler
}

func NewKafkaSource(kconf *KafkaSourceConf) *KafkaSource {
	kafkaSource := &KafkaSource{
		messageHandler: nil,
	}

	config := sarama.NewConfig()
	config.ClientID = "bruco"
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.Initial = kafkaSource.resolveOffset(kconf.Offset)
	balanceStrategy := kafkaSource.resolveBalanceStrategy(kconf.BalanceStrategy)
	config.Consumer.Group.Rebalance.Strategy = balanceStrategy

	consumerGroup, err := sarama.NewConsumerGroup(kconf.Brokers, kconf.ConsumerGroup, config)
	if err != nil {
		log.Panicf("[KAFKA-SOURCE] error creating consumer group client: %v", err)
	}
	go func() {
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := consumerGroup.Consume(context.Background(), kconf.Topics, kafkaSource); err != nil {
				log.Printf("[KAFKA-SOURCE] error from consumer: %v", err)
				time.Sleep(time.Second * 1)
			}
		}
	}()
	log.Printf("[KAFKA-SOURCE] consuming topics: %s", kconf.Topics)

	return kafkaSource
}

// SetMessageHandler sets the callback function that will be invoked on each
// message received from the kafka source
func (k *KafkaSource) SetMessageHandler(handler source.MessageHandler) {
	k.messageHandler = handler
}

func (k *KafkaSource) resolveOffset(offset string) int64 {
	if offset == "" {
		return sarama.OffsetNewest
	}
	switch lower := strings.ToLower(offset); lower {
	case "latest":
		return sarama.OffsetNewest
	case "earliest":
		return sarama.OffsetOldest
	default:
		log.Fatalf("[KAFKA-SOURCE] invalid offset spec. %s", offset)
		return 0
	}
}

func (k *KafkaSource) resolveBalanceStrategy(strategy string) sarama.BalanceStrategy {
	if strategy == "" {
		return sarama.BalanceStrategyRange
	}
	switch lower := strings.ToLower(strategy); lower {
	case "sticky":
		return sarama.BalanceStrategySticky
	case "roundrobin":
		return sarama.BalanceStrategyRoundRobin
	case "range":
		return sarama.BalanceStrategyRange
	default:
		log.Fatalf("[KAFKA-SOURCE] unrecognized balance strategy: %s", strategy)
		return nil
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (k *KafkaSource) Setup(session sarama.ConsumerGroupSession) error {
	log.Printf("[KAFKA-SOURCE] starting consumer session. claims %v", session.Claims())
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (k *KafkaSource) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (k *KafkaSource) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	claimedMessage := make(chan sarama.ConsumerMessage)
	go func() {
		for msg := range claimedMessage {
			if k.messageHandler != nil {
				outMsg := &source.Message{
					Timestamp: msg.Timestamp,
					Value:     msg.Value,
				}
				resolveChan := make(chan error)
				go func(ch chan error) {
					err := <-ch
					if err == nil {
						session.MarkMessage(&msg, "")
					} else {
						log.Printf("[KAFKA-SOURCE] processor error: %s", err)
					}
					close(ch)
				}(resolveChan)
				k.messageHandler(outMsg, resolveChan)
			}
		}
		log.Printf("[KAFKA-SOURCE] message handler stopped for partition %d", claim.Partition())
	}()

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		// lag := claim.HighWaterMarkOffset() - message.Offset
		// log.Printf("### partition: %d lag: %d", claim.Partition(), lag)
		// log.Printf("value = %s, timestamp = %v, topic = %s, partition = %d", string(message.Value), message.Timestamp, message.Topic, claim.Partition())
		claimedMessage <- *message
	}

	close(claimedMessage)
	return nil
}
