package daemon

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jexia/semaphore/cmd/semaphore/daemon/config"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/prettyerr"
	"github.com/spf13/cobra"
)

var flags = &config.Daemon{}

// Command represents the semaphore daemon command
var Command = &cobra.Command{
	Use:          "daemon",
	Short:        "Starts the Semaphore daemon, it will execute with the passed definitions and expose the configured endpoints",
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	Command.PersistentFlags().StringVar(&flags.HTTP.Address, "http", "", "If set starts the HTTP listener on the given TCP address")
	Command.PersistentFlags().StringVar(&flags.GraphQL.Address, "graphql", "", "If set starts the GraphQL listener on the given TCP address")
	Command.PersistentFlags().StringVar(&flags.GRPC.Address, "grpc", "", "If set starts the gRPC listener on the given TCP address")
	Command.PersistentFlags().StringSliceVar(&flags.Protobuffers, "proto", []string{}, "If set are all proto definitions found inside the given path passed as schema definitions, all proto definitions are also passed as imports")
	Command.PersistentFlags().StringSliceVarP(&flags.Files, "file", "f", []string{"config.hcl"}, "Parses the given file as a definition file")
	Command.PersistentFlags().StringVar(&flags.LogLevel, "level", "", "Global logging level, this value will override the defined log level inside the file definitions")
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

	client, err := NewClient(ctx, core, provider)
	if err != nil {
		return err
	}

	go sigterm(client)

	err = client.Serve()
	if err != nil {
		return err
	}

	return nil
}

func sigterm(client *Client) {
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
	<-term
	client.Close()
}
