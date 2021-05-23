package nats

import "github.com/ferama/bruco/pkg/source"

type NatsSourceConf struct {
	source.SourceConfBase

	ServerUrl string `yaml:"serverUrl"`
}
