package httpsource

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ferama/bruco/pkg/source"
)

type HttpSource struct {
	source.SourceBase

	ignoreProcessorResponse bool
	port                    int
}

func NewHttpSource(cfg *HttpSourceConf) *HttpSource {
	port := 8080
	if cfg.Port != 0 {
		port = cfg.Port
	}
	source := &HttpSource{
		port:                    port,
		ignoreProcessorResponse: cfg.IgnoreProcessorResponse,
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

	body, _ := ioutil.ReadAll(r.Body)
	outMsg := &source.Message{
		Timestamp: time.Now(),
		Value:     []byte(body),
	}
	resolveChan := s.MessageHandler(outMsg)

	if !s.ignoreProcessorResponse {
		response := <-resolveChan
		if response.Error != "" {
			log.Printf("[HTTP-SOURCE] processor error: %s", response.Error)
			http.Error(w, fmt.Sprintf("processor error: %s", response.Error), 400)
			return
		}
		if response.ContentType != "" {
			w.Header().Add("Content-Type", response.ContentType)
		}
		fmt.Fprintf(w, response.Data)
	} else {
		fmt.Fprintf(w, "ok")
	}

}
