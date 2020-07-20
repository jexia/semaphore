package validate

import (
	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/cmd/semaphore/config"
	"github.com/jexia/semaphore/pkg/core"
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/spf13/cobra"
)

var params = config.New()

// Command represents the semaphore validate command
var Command = &cobra.Command{
	Use:          "validate",
	Short:        "Validate the flow definitions with the configured schema format(s)",
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

	ctx := instance.NewContext()
	options, err := semaphore.NewOptions(ctx, arguments...)
	if err != nil {
		return err
	}

	_, err = core.Specs(ctx, functions.Collection{}, options)
	if err != nil {
		return err
	}

	return nil
}
