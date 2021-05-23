package nats

import (
	"log"

	"github.com/ferama/bruco/pkg/source"
	"github.com/nats-io/nats.go"
)

type NatsSource struct {
	source.SourceBase
	conn *nats.Conn
}

func NewNatsSource(cfg *NatsSourceConf) *NatsSource {
	source := &NatsSource{}
	nc, err := nats.Connect(cfg.ServerUrl)
	if err != nil {
		log.Printf("[NATS-SOURCE] unable to connect to the server: %s", err)
	}
	source.conn = nc
	return source
}
