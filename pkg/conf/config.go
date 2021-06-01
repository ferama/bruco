package conf

import (
	"fmt"
	"path/filepath"

	"github.com/ferama/bruco/pkg/loader"
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

// Config keeps bruco config information and helps to load them
// It also get the processor code if needed and does all the preparation
// required for the application to run properly
type Config struct {
	Processor *processor.ProcessorConf
	Source    source.SourceConf
	Sink      sink.SinkConf

	// WorkingDir is the directory where the config file resides
	WorkingDir string

	loader *loader.Loader
}

// LoadConfig search for a config.yaml file and parses it
// The config.yaml file could be the first bruco command line argument
// Another option is if bruco is called with an url scheme.
// In that case if the path is a zip archive for example, this function
// will loads and extracts its contents and will search for a config.yaml
// there.
func LoadConfig(fileURL string) (*Config, error) {
	config := &Config{
		loader: loader.NewLoader(),
	}

	fileHandler, err := config.loader.LoadFunction(fileURL)
	if err != nil {
		config.loader.Cleanup()
		return nil, err
	}
	defer fileHandler.Close()

	config.WorkingDir = filepath.Dir(fileHandler.Name())

	var cfgFile struct {
		Processor map[string]interface{} `yaml:"processor"`
		Source    map[string]interface{} `yaml:"source"`
		Sink      map[string]interface{} `yaml:"sink"`
	}

	decoder := yaml.NewDecoder(fileHandler)
	err = decoder.Decode(&cfgFile)
	if err != nil {
		return config, err
	}

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
		return config, fmt.Errorf("invalid source kind: %s", sourceKind)
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

// Cleanup removes temprary files
func (c *Config) Cleanup() {
	c.loader.Cleanup()
}
