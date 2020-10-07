package openapi3

import (
	"os"

	"github.com/jexia/semaphore/cmd/semaphore/daemon/config"
	"github.com/jexia/semaphore/cmd/semaphore/daemon/providers"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/prettyerr"
	"github.com/jexia/semaphore/pkg/providers/openapi3"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var flags = &config.Daemon{}

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
	Command.PersistentFlags().StringSliceVar(&flags.Protobuffers, "proto", []string{}, "If set are all proto definitions found inside the given path passed as schema definitions, all proto definitions are also passed as imports")
	Command.PersistentFlags().StringSliceVarP(&flags.Files, "file", "f", []string{"config.hcl"}, "Parses the given file as a definition file")
	Command.PersistentFlags().StringVar(&flags.LogLevel, "level", "warn", "Global logging level, this value will override the defined log level inside the file definitions")
}

func run(cmd *cobra.Command, args []string) (err error) {

	defer func() {
		if err != nil {
			err = prettyerr.StandardErr(err)
		}
	}()

	ctx := logger.WithLogger(broker.NewContext())
	err = config.SetOptions(ctx, flags)
	if err != nil {
		return err
	}

	core, err := config.NewCore(ctx, flags)
	if err != nil {
		return err
	}

	provider, err := config.NewProviders(ctx, core, flags)
	if err != nil {
		return err
	}

	collection, err := providers.Resolve(ctx, functions.Collection{}, provider)
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
