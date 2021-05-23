package nats

type NatsSourceConf struct {
	ServerUrl string `yaml:"serverUrl"`
}

func (c *NatsSourceConf) IsFireAndForget() bool {
	return true
}
