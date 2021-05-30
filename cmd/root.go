package cmd

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
	"github.com/spf13/cobra"
)

func getEventCallback(eventSink sink.Sink, resolve chan processor.Response) processor.EventCallback {
	return func(response *processor.Response) {
		if eventSink == nil {
			if response.Data != "" {
				log.Println("[ROOT] WARNING: processor has a return value but no sink is configured")
			}
		}
		if response.Error != "" {
			resolve <- *response
			return
		}
		msg := &sink.Message{
			Key:   response.Key,
			Value: []byte(response.Data),
		}
		err := eventSink.Publish(msg)
		if err != nil {
			log.Printf("[ROOT] publish error: %s", err)
		}
		resolve <- *response
	}
}

var rootCmd = &cobra.Command{
	Use:  "bruco config_file_path.yaml",
	Long: "The streaming pipeline processing tool.",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := conf.LoadConfig(args[0])
		if err != nil {
			panic(err)
		}

		eventSource := factory.GetSourceInstance(cfg)
		eventSink := factory.GetSinkInstance(cfg)
		asyncHandler := cfg.Source.IsFireAndForget()
		if asyncHandler {
			log.Println("[ROOT] running in async mode")
		} else {
			log.Println("[ROOT] running in sync mode")
		}

		workers := factory.GetProcessorWorkerPoolInstance(cfg)

		eventSource.SetMessageHandler(func(msg *source.Message) chan processor.Response {
			resolve := make(chan processor.Response)
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
			return resolve
		})

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		// cleanup
		workers.Destroy()
	},
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
