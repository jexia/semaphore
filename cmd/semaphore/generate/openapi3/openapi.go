package openapi3

import (
	"os"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/cmd/semaphore/config"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/providers"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/providers/openapi3"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var params = config.New()

// Command represents the semaphore daemon command
var Command = &cobra.Command{
	Use:   "openapi3",
	Short: "Generates a openapi3 specification",
	Long: `Generates a openapi3 specification.
*NOTE*: This is a experimental feature and not all features are supported yet.`,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	Command.PersistentFlags().StringSliceVar(&params.Protobuffers, "proto", []string{}, "If set are all proto definitions found inside the given path passed as schema definitions, all proto definitions are also passed as imports")
	Command.PersistentFlags().StringSliceVarP(&params.Files, "file", "f", []string{"config.hcl"}, "Parses the given file as a definition file")
	Command.PersistentFlags().StringVar(&params.LogLevel, "level", "warn", "Global logging level, this value will override the defined log level inside the file definitions")
}

func run(cmd *cobra.Command, args []string) error {
	arguments, err := config.ConstructArguments(params)
	if err != nil {
		return err
	}

	ctx := logger.WithLogger(broker.NewContext())
	options, err := semaphore.NewOptions(ctx, arguments...)
	if err != nil {
		return err
	}

	collection, err := providers.Resolve(ctx, functions.Collection{}, options)
	if err != nil {
		return err
	}

	object, err := openapi3.Generate(collection.EndpointList, collection.FlowListInterface)
	if err != nil {
		return err
	}

	bb, err := yaml.Marshal(object)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(bb)
	if err != nil {
		return err
	}

	return nil
}
