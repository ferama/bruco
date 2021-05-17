package conf

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	LambdaPath string                 `yaml:"lambdaPath"`
	Workers    int                    `yaml:"workers"`
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
