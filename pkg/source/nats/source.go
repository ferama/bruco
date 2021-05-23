package natssource

import (
	"log"
	"time"

	"github.com/ferama/bruco/pkg/source"
	"github.com/nats-io/nats.go"
)

type NatsSource struct {
	source.SourceBase

	msgCh      chan *nats.Msg
	conn       *nats.Conn
	subject    string
	queueGroup string
}

func NewNatsSource(cfg *NatsSourceConf) *NatsSource {
	source := &NatsSource{
		subject:    cfg.Subject,
		queueGroup: cfg.QueueGroup,
		msgCh:      make(chan *nats.Msg, 16),
	}
	nc, err := nats.Connect(cfg.ServerUrl)
	if err != nil {
		log.Printf("[NATS-SOURCE] unable to connect to the server: %s", err)
	}
	source.conn = nc
	nc.ChanQueueSubscribe(source.subject, source.queueGroup, source.msgCh)

	go source.consume()
	return source
}

func (s *NatsSource) consume() {
	for msg := range s.msgCh {
		if s.MessageHandler != nil {
			outMsg := &source.Message{
				Timestamp: time.Now(),
				Value:     msg.Data,
			}
			resolveChan := make(chan error)
			go func(ch chan error) {
				err := <-ch
				if err != nil {
					log.Printf("[NATS-SOURCE] processor error: %s", err)
				}
				close(ch)
			}(resolveChan)
			s.MessageHandler(outMsg, resolveChan)
		}
	}
}
