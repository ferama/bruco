package kafkasource

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ferama/bruco/pkg/source"
)

type KafkaSource struct {
	source.SourceBase
}

func NewKafkaSource(kconf *KafkaSourceConf) *KafkaSource {
	kafkaSource := &KafkaSource{}

	config := sarama.NewConfig()
	config.ClientID = "bruco"
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true

	channelBufferSize := 256
	if kconf.ChannelBufferSize != "" {
		var err error
		channelBufferSize, err = strconv.Atoi(kconf.ChannelBufferSize)
		if err != nil {
			log.Fatalf("[KAFKA-SOURCE] invalid channelBufferSize conf %s", kconf.ChannelBufferSize)
		}
	}
	config.ChannelBufferSize = channelBufferSize

	fetchDefaultBytes := 1024 * 1024
	if kconf.FetchDefaultBytes != "" {
		var err error
		fetchDefaultBytes, err = strconv.Atoi(kconf.FetchDefaultBytes)
		if err != nil {
			log.Fatalf("[KAFKA-SOURCE] invalid fetchDefaultBytes conf %s", kconf.FetchDefaultBytes)
		}
	}
	config.Consumer.Fetch.Default = int32(fetchDefaultBytes)

	config.Consumer.MaxProcessingTime = time.Hour * 24
	// config.Version = sarama.V2_4_0_0
	rebalanceTimeout := 60
	if kconf.RebalanceTimeout != 0 {
		rebalanceTimeout = kconf.RebalanceTimeout
	}
	config.Consumer.Group.Rebalance.Timeout = time.Second * time.Duration(rebalanceTimeout)
	// log.Printf("rebalance timeout %d", config.Consumer.Group.Rebalance.Timeout)

	config.Consumer.Offsets.Initial = kafkaSource.resolveOffset(kconf.Offset)
	config.Consumer.Group.Rebalance.Strategy = kafkaSource.resolveBalanceStrategy(kconf.BalanceStrategy)

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
	log.Printf("[KAFKA-SOURCE] ending consumer session. claims %v", session.Claims())
	// log.Println("[KAFKA-SOURCE] cleanup")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (k *KafkaSource) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	claimedMessage := make(chan sarama.ConsumerMessage)
	go func() {
		// log.Printf("[KAFKA-SOURCE] starting message handler for partition %d", claim.Partition())
		for msg := range claimedMessage {
			if k.MessageHandler != nil {
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
				k.MessageHandler(outMsg, resolveChan)
			}
		}
		log.Printf("[KAFKA-SOURCE] message handler stopped for partition %d", claim.Partition())
	}()

	for {
		select {
		case message := <-claim.Messages():
			claimedMessage <- *message
			// log.Printf("value = %s, timestamp = %v, topic = %s, partition = %d", string(message.Value), message.Timestamp, message.Topic, claim.Partition())
			// log.Printf("value = %s, partition = %d", string(message.Value), claim.Partition())
		case <-session.Context().Done():
			close(claimedMessage)
			return nil
		}
	}
}
