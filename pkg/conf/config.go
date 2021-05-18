package conf

import (
	"log"
	"os"

	"github.com/ferama/bruco/pkg/processor"
	"github.com/ferama/bruco/pkg/sink"
	kafkasink "github.com/ferama/bruco/pkg/sink/kafka"
	"github.com/ferama/bruco/pkg/source"
	kafkasource "github.com/ferama/bruco/pkg/source/kafka"
	"gopkg.in/yaml.v2"
)

type Config struct {
	LambdaPath string                 `yaml:"lambdaPath"`
	Workers    int                    `yaml:"workers"`
	Processor  map[string]interface{} `yaml:"processor"`
	Source     map[string]interface{} `yaml:"source"`
	Sink       map[string]interface{} `yaml:"sink"`
}

// LoadConfig parses the [config].yaml file and loads its values
// into the Config struct
func LoadConfig(filePath string) (*Config, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// set some reasonable defaults
	cfg := Config{}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// GerProcessorConf extract the processor configuration from yaml
func (c *Config) GerProcessorConf() *processor.ProcessorConf {
	m, _ := yaml.Marshal(c.Processor)
	conf := &processor.ProcessorConf{}
	yaml.Unmarshal(m, conf)
	return conf
}

// GetSourceConf resoslve the source conf
func (c *Config) GetSourceConf() source.SourceConf {
	sourceKind := c.Source["kind"]
	switch sourceKind {
	case "kafka":
		m, _ := yaml.Marshal(c.Source)
		conf := &kafkasource.KafkaSourceConf{}
		yaml.Unmarshal(m, conf)
		return conf
	default:
		log.Fatalf("Invalid source kind: %s", sourceKind)
		return nil
	}
}

// GetSinkConf resolve the sink conf
func (c *Config) GetSinkConf() sink.SinkConf {
	sinkKind := c.Sink["kind"]
	switch sinkKind {
	case "kafka":
		m, _ := yaml.Marshal(c.Sink)
		conf := &kafkasink.KafkaSinkConf{}
		yaml.Unmarshal(m, conf)
		return conf
	default:
		log.Fatalf("Invalid sink kind: %s", sinkKind)
		return nil
	}
}
