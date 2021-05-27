package factory

import (
	"log"

	"github.com/ferama/bruco/pkg/conf"
	"github.com/ferama/bruco/pkg/processor"
	"github.com/ferama/bruco/pkg/sink"
	kafkasink "github.com/ferama/bruco/pkg/sink/kafka"
	natssink "github.com/ferama/bruco/pkg/sink/nats"
	"github.com/ferama/bruco/pkg/source"
	kafkasource "github.com/ferama/bruco/pkg/source/kafka"
	natssource "github.com/ferama/bruco/pkg/source/nats"
)

// GetSourceInstance builds up a source instance
func GetSourceInstance(cfg *conf.Config) source.Source {
	var eventSource source.Source

	sourceKind := cfg.Source.GetKind()
	switch sourceKind {
	case "kafka":
		eventSource = kafkasource.NewKafkaSource(cfg.Source.(*kafkasource.KafkaSourceConf))
		return eventSource
	case "nats":
		eventSource = natssource.NewNatsSource(cfg.Sink.(*natssource.NatsSourceConf))
		return eventSource
	default:
		log.Fatalf("[ROOT] invalid source kind: %s", sourceKind)
		return nil
	}
}

// GetSinkInstance builds up a sink instance
func GetSinkInstance(cfg *conf.Config) sink.Sink {
	var eventSink sink.Sink

	sinkKind := cfg.Sink.GetKind()
	switch sinkKind {
	case "kafka":
		eventSink = kafkasink.NewKafkaSink(cfg.Sink.(*kafkasink.KafkaSinkConf))
		return eventSink
	case "nats":
		eventSink = natssink.NewKNatsSink(cfg.Sink.(*natssink.NatsSinkConf))
		return eventSink
	default:
		log.Fatalf("[ROOT] invalid sink kind: %s", sinkKind)
		return nil
	}
}

// GetProcessorWorkerPool build up a processor worker pool
func GetProcessorWorkerPool(cfg *conf.Config) *processor.Pool {
	return processor.NewPool(cfg.Processor)
}
