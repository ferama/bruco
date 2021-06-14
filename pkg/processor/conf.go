package processor

// EnvVar env vars name => value map
type EnvVar struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

// ProcessorConf holds the processor configuration
type ProcessorConf struct {
	HandlerPath string   `json:"handlerPath,omitempty" yaml:"handlerPath,omitempty"`
	ModuleName  string   `json:"moduleName,omitempty" yaml:"moduleName,omitempty"`
	Workers     int      `json:"workers" yaml:"workers"`
	Env         []EnvVar `json:"env,omitempty" yaml:"env,omitempty"`
}
