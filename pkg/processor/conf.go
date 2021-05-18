package processor

type ProcessorConf struct {
	WorkDir    string `yaml:"workDir"`
	ModuleName string `yaml:"moduleName"`
	Workers    int    `yaml:"workers"`
}
