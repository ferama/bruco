package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ferama/bruco/pkg/conf"
	"github.com/ferama/bruco/pkg/processor"
	"github.com/ferama/bruco/pkg/source"
	"github.com/ferama/bruco/pkg/source/kafka"
	"gopkg.in/yaml.v2"
)

func handleEventResponse(response *processor.Response) {
	log.Println(response.Data)
}

func main() {
	cfg, err := conf.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	var eventSource source.Source
	sourceKind := cfg.Source["kind"]
	asyncHandler := true

	switch sourceKind {
	case "kafka":
		m, _ := yaml.Marshal(cfg.Source)
		conf := &kafka.KafkaSourceConf{}
		yaml.Unmarshal(m, conf)
		eventSource = kafka.NewKafkaSource(conf)
		asyncHandler = conf.AsyncHandler
	default:
		log.Fatalf("Invalid source kind: %s", sourceKind)
	}

	workers := processor.NewPool(cfg.Workers, cfg.LambdaPath)

	eventSource.SetMessageHandler(func(msg *source.Message) {
		if asyncHandler {
			// NOTE: the async handler version will not guarantee
			// messages handling order between same partition
			workers.HandleEventAsync(msg.Value, handleEventResponse)
		} else {
			response, err := workers.HandleEvent(msg.Value)
			if err != nil {
				log.Println(err)
				return
			}
			handleEventResponse(response)
		}
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
