package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jexia/maestro"
	"github.com/jexia/maestro/cmd/maestro/config"
	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/codec/proto"
	"github.com/jexia/maestro/constructor"
	"github.com/jexia/maestro/definitions/hcl"
	"github.com/jexia/maestro/protocol/graphql"
	"github.com/jexia/maestro/protocol/http"
	"github.com/jexia/maestro/protocol/micro"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jexia/maestro/specs"
	"github.com/micro/go-micro/service/grpc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var global = config.New()

// Cmd represents the maestro run command
var Cmd = &cobra.Command{
	Use:   "run",
	Short: "Run the flow definitions with the configured schema format(s)",
	RunE:  run,
}

func init() {
	Cmd.PersistentFlags().StringP("config", "c", "", "Config file path")
	Cmd.PersistentFlags().StringVar(&global.HTTP.Address, "http", "", "If set starts the HTTP listener on the given TCP address")
	Cmd.PersistentFlags().StringVar(&global.GraphQL.Address, "graphql", "", "If set starts the GraphQL listener on the given TCP address")
	Cmd.PersistentFlags().StringSliceVar(&global.Protobuffers, "proto", []string{}, "If set are all proto definitions found inside the given path passed as schema definitions, all proto definitions are also passed as imports")
	Cmd.PersistentFlags().StringSliceVar(&global.Flows, "flow", []string{}, "If set are all flow definitions inside the given path passed as flow definitions")
	Cmd.PersistentFlags().StringVar(&global.LogLevel, "level", "info", "Logging level")
}

func run(cmd *cobra.Command, args []string) error {
	err := config.Read(cmd, global)
	if err != nil {
		return err
	}

	level, err := logrus.ParseLevel(global.LogLevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(level)

	options := []constructor.Option{
		maestro.WithCodec(json.NewConstructor()),
		maestro.WithCodec(proto.NewConstructor()),
		maestro.WithCaller(micro.New("grpc", grpc.NewService())),
		maestro.WithCaller(http.NewCaller()),
	}

	for _, flow := range global.Flows {
		options = append(options, maestro.WithDefinitions(hcl.DefinitionResolver(flow)))
	}

	for _, path := range global.Protobuffers {
		resolver, err := protoc.Collect(global.Protobuffers, path)
		if err != nil {
			return err
		}

		options = append(options, maestro.WithSchema(resolver))
	}

	if global.HTTP.Address != "" {
		options = append(options, maestro.WithListener(http.NewListener(global.HTTP.Address, specs.Options{})))
	}

	if global.GraphQL.Address != "" {
		options = append(options, maestro.WithListener(graphql.NewListener(global.GraphQL.Address, specs.Options{})))
	}

	client, err := maestro.New(options...)
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
