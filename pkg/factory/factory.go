package factory

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

// GetSourceInstance builds up a source instance
func GetSourceInstance(cfg *conf.Config) (source.Source, source.SourceConf) {
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
		log.Fatalf("[ROOT] invalid source kind: %s", sourceKind)
		return nil, nil
	}
}

// GetSinkInstance builds up a sink instance
func GetSinkInstance(cfg *conf.Config) (sink.Sink, sink.SinkConf) {
	var eventSink sink.Sink

	sinkKind := cfg.Sink["kind"]
	switch sinkKind {
	case "kafka":
		m, _ := yaml.Marshal(cfg.Sink)
		conf := &kafkasink.KafkaSinkConf{}
		yaml.Unmarshal(m, conf)
		eventSink = kafkasink.NewKafkaSink(conf)
		return eventSink, conf
	default:
		log.Printf("[ROOT] sink kind: %s", sinkKind)
		return nil, nil
	}
}

// GetProcessorWorkerPool build up a processor worker pool
func GetProcessorWorkerPool(cfg *conf.Config) *processor.Pool {
	m, _ := yaml.Marshal(cfg.Processor)
	conf := &processor.ProcessorConf{}
	yaml.Unmarshal(m, conf)
	return processor.NewPool(conf)
}
