package processor

type ProcessorConf struct {
	LambdaPath string `yaml:"lambdaPath"`
	Workers    int    `yaml:"workers"`
}
