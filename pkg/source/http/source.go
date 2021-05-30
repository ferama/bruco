package httpsource

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ferama/bruco/pkg/processor"
	"github.com/ferama/bruco/pkg/source"
)

type HttpSource struct {
	source.SourceBase

	port string
}

func NewHttpSource(cfg *HttpSourceConf) *HttpSource {
	source := &HttpSource{
		port: "8080",
	}

	http.HandleFunc("/", source.httpHandler)
	go func() {
		addr := fmt.Sprintf(":%s", source.port)
		log.Printf("[HTTP-SOURCE] listening on: %s", addr)
		log.Fatal(http.ListenAndServe(addr, nil))
	}()
	return source
}

func (s *HttpSource) httpHandler(w http.ResponseWriter, r *http.Request) {
	if s.MessageHandler == nil {
		log.Panicln("[HTTP-SOURCE] you need to set a message handler for http source")
	}
	outMsg := &source.Message{
		Timestamp: time.Now(),
		Value:     []byte("prova"),
	}
	resolveChan := s.MessageHandler(outMsg)
	go func(ch chan processor.Response) {
		response := <-ch
		if response.Error != "" {
			log.Printf("[HTTP-SOURCE] processor error: %s", response.Error)
		}
		close(ch)
	}(resolveChan)

}
