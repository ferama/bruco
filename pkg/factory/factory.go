package factory

import (
	"log"

	"github.com/ferama/bruco/pkg/conf"
	"github.com/ferama/bruco/pkg/processor"
	"github.com/ferama/bruco/pkg/sink"
	kafkasink "github.com/ferama/bruco/pkg/sink/kafka"
	natssink "github.com/ferama/bruco/pkg/sink/nats"
	"github.com/ferama/bruco/pkg/source"
	httpsource "github.com/ferama/bruco/pkg/source/http"
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
		eventSource = natssource.NewNatsSource(cfg.Source.(*natssource.NatsSourceConf))
		return eventSource
	case "http":
		eventSource = httpsource.NewHttpSource(cfg.Source.(*httpsource.HttpSourceConf))
		return eventSource
	default:
		log.Fatalf("invalid source kind: %s", sourceKind)
		return nil
	}
}

// GetSinkInstance builds up a sink instance
func GetSinkInstance(cfg *conf.Config) sink.Sink {
	// sink is not mandatory
	if cfg.Sink == nil {
		return nil
	}

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
		log.Fatalf("invalid sink kind: %s", sinkKind)
		return nil
	}
}

// GetProcessorWorkerPoolInstance build up a processor worker pool
func GetProcessorWorkerPoolInstance(cfg *conf.Config) *processor.Pool {
	if cfg.Processor == nil {
		return nil
	}
	return processor.NewPool(cfg.Processor, cfg.WorkingDir)
}
