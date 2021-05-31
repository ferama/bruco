package conf

import (
	"fmt"
	"io/ioutil"
	"os"
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

type Config struct {
	Processor *processor.ProcessorConf
	Source    source.SourceConf
	Sink      sink.SinkConf

	WorkingDir string

	loader *loader.Loader
}

// LoadConfig parses the [config].yaml file and loads its values
// into the Config struct
func LoadConfig(fileURL string) (*Config, error) {
	config := &Config{
		loader: loader.NewLoader(),
	}

	fileHandler, err := config.findConfig(fileURL)
	if err != nil {
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
		return nil, err
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

func (c *Config) findConfig(fileURL string) (*os.File, error) {
	var fileHandler *os.File
	var err error

	filePath, err := c.loader.Load(fileURL)
	if err != nil {
		return nil, err
	}

	fileHandler, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}
	fi, err := fileHandler.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		// it's a directory
		path := filepath.Join(filePath, "config.yaml")
		fileHandler.Close()
		fileHandler, err = os.Open(path)

		if err != nil {
			entries, _ := ioutil.ReadDir(filePath)
			if len(entries) > 0 {
				path := filepath.Join(filePath, entries[0].Name(), "config.yaml")
				fileHandler.Close()
				fileHandler, err = os.Open(path)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return fileHandler, nil
}

func (c *Config) Cleanup() {
	c.loader.Cleanup()
}
