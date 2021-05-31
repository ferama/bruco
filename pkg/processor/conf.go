package processor

// EnvVar env vars name => value map
type EnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// ProcessorConf holds the processor configuration
type ProcessorConf struct {
	HandlerPath string   `yaml:"handlerPath"`
	ModuleName  string   `yaml:"moduleName"`
	Workers     int      `yaml:"workers"`
	Env         []EnvVar `yaml:"env"`
}
