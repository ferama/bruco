package processor

// ProcessorConf holds the processor configuration
type ProcessorConf struct {
	WorkDir    string            `yaml:"workDir"`
	ModuleName string            `yaml:"moduleName"`
	Workers    int               `yaml:"workers"`
	Env        map[string]string `yaml:"env"`
}
