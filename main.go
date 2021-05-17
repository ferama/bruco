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
		Brokers:       []string{"localhost:9092"},
		Topics:        []string{"test"},
		ConsumerGroup: "my-consumer-group",
	})

	workers := processor.NewPool(4, "./hack/lambda")
	callback := func(response *processor.Response) {
		log.Println(response.Data)
	}

	source.SetMessageHandler(func(msg *sarama.ConsumerMessage) {
		// workers.HandleEvent(msg.Value)
		// NOTE: the async handler version will not guarantee
		// messages handling order between same partition
		workers.HandleEventAsync(msg.Value, callback)
	})

	c := make(chan os.Signal, 10)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
