package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/ferama/coreai/pkg/processor"
	"github.com/ferama/coreai/pkg/source/kafka"
)

func main() {
	source := kafka.NewKafkaSource(&kafka.KafkaConf{
		Brokers:         []string{"localhost:9092"},
		Topics:          []string{"test"},
		ConsumerGroup:   "my-consumer-group",
		AutoMarkMessage: false,
	})

	c := make(chan os.Signal, 10)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	workers := processor.NewPool(3, "./hack/lambda")
	// callback := func(response *pool.Response) {
	// 	log.Println(response.Data)
	// }
	getCallback := func(source *kafka.KafkaSource, msg *sarama.ConsumerMessage) processor.EvenCallback {
		return func(response *processor.Response) {
			// log.Println(response.Data)
			source.MarkMessage(msg)
		}
	}
	for {
		select {
		case <-c:
			log.Println("quit")
			workers.Destroy()
			return
		case msg := <-source.Stream:
			// log.Printf("value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)
			// workers.HandleEvent(msg.Value, callback)
			workers.HandleEvent(msg.Value, getCallback(source, &msg))
		}
	}
}
