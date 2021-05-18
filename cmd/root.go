package cmd

import (
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
	return func(response *processor.Response) {
		eventSink.Publish(response.Key, []byte(response.Data))
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
		asyncHandler := sourceConf.IsAsyncHandler()

		workers := processor.NewPool(cfg.GerProcessorConf())

		eventSource.SetMessageHandler(func(msg *source.Message) {
			if asyncHandler {
				// NOTE: the async handler version will not guarantee
				// messages handling order between same partition
				workers.HandleEventAsync(msg.Value, getEventCallback(eventSink))
			} else {
				response, err := workers.HandleEvent(msg.Value)
				if err != nil {
					log.Println(err)
					return
				}
				getEventCallback(eventSink)(response)
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
