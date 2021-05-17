package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

type MessageHandler func(msg *sarama.ConsumerMessage)

type KafkaSource struct {
	consumerGroupSession sarama.ConsumerGroupSession
	messageHandler       MessageHandler
}

func NewKafkaSource(kconf *KafkaConf) *KafkaSource {
	kafkaSource := &KafkaSource{
		messageHandler: nil,
	}

	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	balanceStrategy, err := kafkaSource.resolveBalanceStrategy(kconf.BalanceStrategy)
	if err != nil {
		log.Panicln(err)
	}
	config.Consumer.Group.Rebalance.Strategy = balanceStrategy

	consumerGroup, err := sarama.NewConsumerGroup(kconf.Brokers, kconf.ConsumerGroup, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}
	go func() {
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := consumerGroup.Consume(context.Background(), kconf.Topics, kafkaSource); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
		}
	}()

	// log.Println("Sarama consumer up and running!...")

	return kafkaSource
}

func (k *KafkaSource) SetMessageHandler(handler MessageHandler) {
	k.messageHandler = handler
}

func (k *KafkaSource) resolveBalanceStrategy(strategy string) (sarama.BalanceStrategy, error) {
	if strategy == "" {
		return sarama.BalanceStrategyRange, nil
	}
	switch strategy {
	case "sticky":
		return sarama.BalanceStrategySticky, nil
	case "roundrobin":
		return sarama.BalanceStrategyRoundRobin, nil
	case "range":
		return sarama.BalanceStrategyRange, nil
	default:
		return nil, fmt.Errorf("unrecognized balance strategy: %s", strategy)
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (k *KafkaSource) Setup(session sarama.ConsumerGroupSession) error {
	k.consumerGroupSession = session
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (k *KafkaSource) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (k *KafkaSource) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Printf("Consume claim called for %d", claim.Partition())

	claimedMessage := make(chan sarama.ConsumerMessage)
	go func() {
		for {
			msg := <-claimedMessage
			// log.Printf("value = %s, partition = %d", string(msg.Value), claim.Partition())
			if k.messageHandler != nil {
				k.messageHandler(&msg)
			}
			k.MarkMessage(&msg)
		}
	}()

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		// log.Printf("value = %s, timestamp = %v, topic = %s, partition = %d", string(message.Value), message.Timestamp, message.Topic, claim.Partition())
		claimedMessage <- *message
	}

	close(claimedMessage)

	return nil
}

func (k *KafkaSource) MarkMessage(msg *sarama.ConsumerMessage) {
	k.consumerGroupSession.MarkMessage(msg, "")
}
