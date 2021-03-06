package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ferama/bruco/pkg/conf"
	"github.com/ferama/bruco/pkg/factory"
	"github.com/ferama/bruco/pkg/processor"
	"github.com/ferama/bruco/pkg/sink"
	"github.com/ferama/bruco/pkg/source"
)

func getEventCallback(eventSink sink.Sink, resolve chan processor.Response) processor.EventCallback {
	return func(response *processor.Response) {
		if eventSink == nil {
			resolve <- *response
			return
		}
		if response.Error != "" {
			resolve <- *response
			return
		}
		// If we are here, build sink output message
		msg := &sink.Message{
			Key:   response.Key,
			Value: []byte(response.Data),
		}
		err := eventSink.Publish(msg)
		if err != nil {
			log.Printf("[BRUCO] publish error: %s", err)
		}
		resolve <- *response
	}
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalf("[BRUCO] Usage: bruco function_path")
	}
	cfg, err := conf.LoadConfig(args[0])
	if err != nil {
		cfg.Cleanup()
		log.Fatalf("[BRUCO] %s", err)
	}

	eventSource := factory.GetSourceInstance(cfg)
	eventSink := factory.GetSinkInstance(cfg)
	asyncHandler := cfg.Source.IsFireAndForget()
	if asyncHandler {
		log.Println("[BRUCO] running in async mode")
	} else {
		log.Println("[BRUCO] running in sync mode")
	}

	workers := factory.GetProcessorWorkerPoolInstance(cfg)

	eventSource.SetMessageHandler(func(msg *source.Message) chan processor.Response {
		// IMPORTANT: resolve chan need to be buffered. The buffer
		// size should be exatcly 1 (one response for each request)
		// It I make an unbuffered channel the channel writer will block
		// until a reader is ready to consume the message.
		resolve := make(chan processor.Response, 1)
		if workers == nil {
			// no processor defined. Copy source to sink
			response := processor.Response{
				Data:  string(msg.Value),
				Error: "",
			}
			getEventCallback(eventSink, resolve)(&response)
		} else {
			if asyncHandler {
				// NOTE: the async handler version will not guarantee
				// messages handling order between same partition
				workers.HandleEventAsync(msg.Value, getEventCallback(eventSink, resolve))
			} else {
				response, err := workers.HandleEvent(msg.Value)
				if err != nil {
					resolve <- processor.Response{
						Data:  "",
						Error: err.Error(),
					}
				} else {
					getEventCallback(eventSink, resolve)(response)
				}
			}
		}
		return resolve
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	// cleanup
	if workers != nil {
		workers.Destroy()
	}
	cfg.Cleanup()
}
