package natssink

import (
	"log"

	"github.com/ferama/bruco/pkg/sink"
	"github.com/nats-io/nats.go"
)

type NatsSink struct {
	conn    *nats.Conn
	subject string
}

// NewKNatsSink creates a new sink object given the conf
func NewKNatsSink(cfg *NatsSinkConf) *NatsSink {
	sink := &NatsSink{
		subject: cfg.Subject,
	}
	nc, err := nats.Connect(cfg.ServerUrl)
	if err != nil {
		log.Printf("[NATS-SOURCE] unable to connect to the server: %s", err)
	}
	sink.conn = nc
	return sink
}

// Publish send a message through the sink
func (s *NatsSink) Publish(msg *sink.Message) error {
	subject := s.subject
	if msg.Key != "" {
		subject = msg.Key
	}
	return s.conn.Publish(subject, msg.Value)
}
