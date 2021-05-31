package conf

import (
	"fmt"
	"os"

	"github.com/ferama/bruco/pkg/processor"
	"github.com/ferama/bruco/pkg/sink"
	kafkasink "github.com/ferama/bruco/pkg/sink/kafka"
	natssink "github.com/ferama/bruco/pkg/sink/nats"
	"github.com/ferama/bruco/pkg/source"
	httpsource "github.com/ferama/bruco/pkg/source/http"
	kafkasource "github.com/ferama/bruco/pkg/source/kafka"
	natssource "github.com/ferama/bruco/pkg/source/nats"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Processor *processor.ProcessorConf
	Source    source.SourceConf
	Sink      sink.SinkConf
}

// LoadConfig parses the [config].yaml file and loads its values
// into the Config struct
func LoadConfig(filePath string) (*Config, error) {
	var cfgFile struct {
		Processor map[string]interface{} `yaml:"processor"`
		Source    map[string]interface{} `yaml:"source"`
		Sink      map[string]interface{} `yaml:"sink"`
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfgFile)
	if err != nil {
		return nil, err
	}
	config := &Config{}

	// source
	sourceKind := cfgFile.Source["kind"]
	switch sourceKind {
	case "kafka":
		m, _ := yaml.Marshal(cfgFile.Source)
		c := &kafkasource.KafkaSourceConf{}
		yaml.Unmarshal(m, c)
		config.Source = c
	case "nats":
		m, _ := yaml.Marshal(cfgFile.Source)
		c := &natssource.NatsSourceConf{}
		yaml.Unmarshal(m, c)
		config.Source = c
	case "http":
		m, _ := yaml.Marshal(cfgFile.Source)
		c := &httpsource.HttpSourceConf{}
		yaml.Unmarshal(m, c)
		config.Source = c
	default:
		return nil, fmt.Errorf("invalid source kind: %s", sourceKind)
	}

	// sink
	// sink is not mandatory so no default in this switch/case
	sinkKind := cfgFile.Sink["kind"]
	switch sinkKind {
	case "kafka":
		m, _ := yaml.Marshal(cfgFile.Sink)
		c := &kafkasink.KafkaSinkConf{}
		yaml.Unmarshal(m, c)
		config.Sink = c
	case "nats":
		m, _ := yaml.Marshal(cfgFile.Sink)
		c := &natssink.NatsSinkConf{}
		yaml.Unmarshal(m, c)
		config.Sink = c
	}

	// processor
	if len(cfgFile.Processor) != 0 {
		m, _ := yaml.Marshal(cfgFile.Processor)
		c := &processor.ProcessorConf{}
		yaml.Unmarshal(m, c)
		config.Processor = c
	}

	return config, nil
}
