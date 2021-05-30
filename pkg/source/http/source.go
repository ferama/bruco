package httpsource

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ferama/bruco/pkg/source"
)

type HttpSource struct {
	source.SourceBase

	port int
}

func NewHttpSource(cfg *HttpSourceConf) *HttpSource {
	port := 8080
	if cfg.Port != 0 {
		port = cfg.Port
	}
	source := &HttpSource{
		port: port,
	}

	http.HandleFunc("/", source.httpHandler)
	go func() {
		addr := fmt.Sprintf(":%d", source.port)
		log.Printf("[HTTP-SOURCE] listening on: %s", addr)
		log.Fatal(http.ListenAndServe(addr, nil))
	}()
	return source
}

func (s *HttpSource) httpHandler(w http.ResponseWriter, r *http.Request) {
	if s.MessageHandler == nil {
		log.Panicln("[HTTP-SOURCE] you need to set a message handler for http source")
	}
	r.ParseForm()
	body := r.Form.Get("")
	outMsg := &source.Message{
		Timestamp: time.Now(),
		Value:     []byte(body),
	}
	resolveChan := s.MessageHandler(outMsg)
	response := <-resolveChan
	if response.Error != "" {
		log.Printf("[HTTP-SOURCE] processor error: %s", response.Error)
	}
	w.Header().Add("Content-Type", response.ContentType)
	fmt.Fprintf(w, response.Data)
}
