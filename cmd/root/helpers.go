package rootcmd

import (
	"log"

	"github.com/ferama/bruco/pkg/conf"
	"github.com/ferama/bruco/pkg/processor"
	"github.com/ferama/bruco/pkg/sink"
	kafkasink "github.com/ferama/bruco/pkg/sink/kafka"
	"github.com/ferama/bruco/pkg/source"
	kafkasource "github.com/ferama/bruco/pkg/source/kafka"
	"gopkg.in/yaml.v2"
)

func GerProcessorConf(cfg *conf.Config) *processor.ProcessorConf {
	m, _ := yaml.Marshal(cfg.Processor)
	conf := &processor.ProcessorConf{}
	yaml.Unmarshal(m, conf)
	return conf
}

func GetEventSource(cfg *conf.Config) (source.Source, source.SourceConf) {
	var eventSource source.Source

	sourceKind := cfg.Source["kind"]
	switch sourceKind {
	case "kafka":
		conf := cfg.GetSourceConf().(*kafkasource.KafkaSourceConf)
		eventSource = kafkasource.NewKafkaSource(conf)
		return eventSource, conf
	default:
		log.Fatalf("Invalid source kind: %s", sourceKind)
		return nil, nil
	}
}

func GetEventSink(cfg *conf.Config) (sink.Sink, sink.SinkConf) {
	var eventSink sink.Sink

	sinkKind := cfg.Sink["kind"]
	switch sinkKind {
	case "kafka":
		conf := cfg.GetSinkConf().(*kafkasink.KafkaSinkConf)
		eventSink = kafkasink.NewKafkaSink(conf)
		return eventSink, conf
	default:
		log.Fatalf("Invalid sink kind: %s", sinkKind)
		return nil, nil
	}
}
