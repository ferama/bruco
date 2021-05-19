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

func getEventCallback(eventSink sink.Sink) processor.EventCallback {
	return func(response *processor.Response) error {
		if eventSink == nil {
			if response.Data != "" {
				log.Println("[ROOT] WARNING: processor has a return value but no sink is configured")
			}
			return nil
		}
		if response.Error != "" {
			return errors.New(response.Error)
		}
		if len(response.Data) == 0 {
			return nil
		}
		msg := &sink.Message{
			Key:   response.Key,
			Value: []byte(response.Data),
		}
		eventSink.Publish(msg)
		return nil
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

		eventSource.SetMessageHandler(func(msg *source.Message) error {
			if asyncHandler {
				// NOTE: the async handler version will not guarantee
				// messages handling order between same partition
				workers.HandleEventAsync(msg.Value, getEventCallback(eventSink))
				// the async handler never returns an error to the source consumers.
				// The source will always consider the message as successfully processesd
				return nil
			} else {
				response, err := workers.HandleEvent(msg.Value)
				if err != nil {
					log.Println("event processing error: ", err)
					return err
				}
				return getEventCallback(eventSink)(response)
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
