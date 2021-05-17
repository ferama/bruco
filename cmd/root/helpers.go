package rootcmd

import (
	"log"

	"github.com/ferama/bruco/pkg/conf"
	"github.com/ferama/bruco/pkg/sink"
	kafkasink "github.com/ferama/bruco/pkg/sink/kafka"
	"github.com/ferama/bruco/pkg/source"
	kafkasource "github.com/ferama/bruco/pkg/source/kafka"
	"gopkg.in/yaml.v2"
)

func GetEventSource(cfg *conf.Config) (source.Source, *kafkasource.KafkaSourceConf) {
	var eventSource source.Source

	sourceKind := cfg.Source["kind"]
	switch sourceKind {
	case "kafka":
		m, _ := yaml.Marshal(cfg.Source)
		conf := &kafkasource.KafkaSourceConf{}
		yaml.Unmarshal(m, conf)
		eventSource = kafkasource.NewKafkaSource(conf)
		return eventSource, conf
	default:
		log.Fatalf("Invalid source kind: %s", sourceKind)
		return nil, nil
	}
}

func GetEventSink(cfg *conf.Config) (sink.Sink, *kafkasink.KafkaSinkConf) {
	var eventSink sink.Sink

	sinkKind := cfg.Sink["kind"]
	switch sinkKind {
	case "kafka":
		m, _ := yaml.Marshal(cfg.Source)
		conf := &kafkasink.KafkaSinkConf{}
		yaml.Unmarshal(m, conf)
		eventSink = kafkasink.NewKafkaSink(conf)
		return eventSink, conf
	default:
		log.Fatalf("Invalid sink kind: %s", sinkKind)
		return nil, nil
	}
}
