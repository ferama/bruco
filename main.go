package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ferama/bruco/pkg/processor"
	"github.com/ferama/bruco/pkg/source"
	"github.com/ferama/bruco/pkg/source/kafka"
)

func main() {
	kafkaSource := kafka.NewKafkaSource(&kafka.KafkaConf{
		Brokers:       []string{"localhost:9092"},
		Topics:        []string{"test"},
		ConsumerGroup: "my-consumer-group",
	})

	workers := processor.NewPool(4, "./hack/lambda")
	// callback := func(response *processor.Response) {
	// 	log.Println(response.Data)
	// }

	kafkaSource.SetMessageHandler(func(msg *source.Message) {
		// workers.HandleEvent(msg.Value)
		// NOTE: the async handler version will not guarantee
		// messages handling order between same partition
		workers.HandleEventAsync(msg.Value, nil)
	})

	c := make(chan os.Signal, 10)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
