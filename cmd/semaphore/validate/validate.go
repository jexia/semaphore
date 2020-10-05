package validate

import (
	"github.com/jexia/semaphore/cmd/semaphore/daemon/config"
	"github.com/jexia/semaphore/cmd/semaphore/daemon/providers"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/prettyerr"
	"github.com/spf13/cobra"
)

var flags = &config.Daemon{}

// Command represents the semaphore validate command
var Command = &cobra.Command{
	Use:          "validate",
	Short:        "Validate the flow definitions with the configured schema format(s)",
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

	_, err = providers.Resolve(ctx, functions.Collection{}, provider)
	if err != nil {
		return err
	}

	return nil
}
