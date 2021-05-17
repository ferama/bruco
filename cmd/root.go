package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ferama/bruco/pkg/conf"
	"github.com/ferama/bruco/pkg/processor"
	"github.com/ferama/bruco/pkg/source"
	"github.com/ferama/bruco/pkg/source/kafka"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func handleEventResponse(response *processor.Response) {
	// log.Println(response.Data)
}

var rootCmd = &cobra.Command{
	Use:  "bruco config_file_path.yaml",
	Long: "The pipeline processing tool.",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := conf.LoadConfig(args[0])
		if err != nil {
			panic(err)
		}

		var eventSource source.Source
		sourceKind := cfg.Source["kind"]
		asyncHandler := true

		switch sourceKind {
		case "kafka":
			m, _ := yaml.Marshal(cfg.Source)
			conf := &kafka.KafkaSourceConf{}
			yaml.Unmarshal(m, conf)
			eventSource = kafka.NewKafkaSource(conf)
			asyncHandler = conf.AsyncHandler
		default:
			log.Fatalf("Invalid source kind: %s", sourceKind)
		}

		workers := processor.NewPool(cfg.Workers, cfg.LambdaPath)

		eventSource.SetMessageHandler(func(msg *source.Message) {
			if asyncHandler {
				// NOTE: the async handler version will not guarantee
				// messages handling order between same partition
				workers.HandleEventAsync(msg.Value, handleEventResponse)
			} else {
				response, err := workers.HandleEvent(msg.Value)
				if err != nil {
					log.Println(err)
					return
				}
				handleEventResponse(response)
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
