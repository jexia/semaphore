package daemon

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jexia/maestro"
	"github.com/jexia/maestro/cmd/maestro/config"
	"github.com/spf13/cobra"
)

var params = config.New()

// Command represents the maestro daemon command
var Command = &cobra.Command{
	Use:   "daemon",
	Short: "Starts the Maestro daemon, it will execute with the passed definitions and expose the configured endpoints",
	RunE:  run,
}

func init() {
	Command.PersistentFlags().StringVar(&params.HTTP.Address, "http", "", "If set starts the HTTP listener on the given TCP address")
	Command.PersistentFlags().StringVar(&params.GraphQL.Address, "graphql", "", "If set starts the GraphQL listener on the given TCP address")
	Command.PersistentFlags().StringVar(&params.GRPC.Address, "grpc", "", "If set starts the gRPC listener on the given TCP address")
	Command.PersistentFlags().StringSliceVar(&params.Protobuffers, "proto", []string{}, "If set are all proto definitions found inside the given path passed as schema definitions, all proto definitions are also passed as imports")
	Command.PersistentFlags().StringSliceVarP(&params.Files, "file", "f", []string{"config.hcl"}, "Parses the given file as a definition file")
	Command.PersistentFlags().StringVar(&params.LogLevel, "level", "", "Global logging level, this value will override the defined log level inside the file definitions")
}

func run(cmd *cobra.Command, args []string) error {
	arguments, err := config.ConstructArguments(params)
	if err != nil {
		return err
	}

	client, err := maestro.New(arguments...)
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

func sigterm(client *maestro.Client) {
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
	<-term
	client.Close()
}
