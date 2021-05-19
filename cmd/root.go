package cmd

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	rootcmd "github.com/ferama/bruco/cmd/root"
	"github.com/ferama/bruco/pkg/conf"
	"github.com/ferama/bruco/pkg/processor"
	"github.com/ferama/bruco/pkg/sink"
	"github.com/ferama/bruco/pkg/source"
	"github.com/spf13/cobra"
)

func getEventCallback(eventSink sink.Sink, resolve chan error) processor.EventCallback {
	return func(response *processor.Response) {
		if eventSink == nil {
			if response.Data != "" {
				log.Println("[ROOT] WARNING: processor has a return value but no sink is configured")
			}
			resolve <- nil
			return
		}
		if response.Error != "" {
			resolve <- errors.New(response.Error)
			return
		}
		if len(response.Data) == 0 {
			resolve <- nil
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
		resolve <- err
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

		eventSource, sourceConf := rootcmd.GetEventSource(cfg)
		eventSink, _ := rootcmd.GetEventSink(cfg)
		asyncHandler := sourceConf.IsFireAndForget()

		workers := processor.NewPool(cfg.GerProcessorConf())

		eventSource.SetMessageHandler(func(msg *source.Message, resolve chan error) {
			if asyncHandler {
				// NOTE: the async handler version will not guarantee
				// messages handling order between same partition
				workers.HandleEventAsync(msg.Value, getEventCallback(eventSink, resolve))
			} else {
				response, err := workers.HandleEvent(msg.Value)
				if err != nil {
					resolve <- err
				} else {
					getEventCallback(eventSink, resolve)(response)
				}
			}
		})

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
	},
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
