package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jexia/maestro"
	"github.com/jexia/maestro/cmd/maestro/config"
	"github.com/jexia/maestro/internal/constructor"
	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/pkg/codec/json"
	"github.com/jexia/maestro/pkg/codec/proto"
	"github.com/jexia/maestro/pkg/definitions/hcl"
	"github.com/jexia/maestro/pkg/definitions/protoc"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport/graphql"
	"github.com/jexia/maestro/pkg/transport/grpc"
	"github.com/jexia/maestro/pkg/transport/http"
	"github.com/jexia/maestro/pkg/transport/micro"
	microGRPC "github.com/micro/go-micro/service/grpc"
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
	Cmd.PersistentFlags().StringVar(&global.GRPC.Address, "grpc", "", "If set starts the gRPC listener on the given TCP address")
	Cmd.PersistentFlags().StringSliceVar(&global.Protobuffers, "proto", []string{}, "If set are all proto definitions found inside the given path passed as schema definitions, all proto definitions are also passed as imports")
	Cmd.PersistentFlags().StringSliceVar(&global.HCL, "hcl", []string{}, "If set are all definitions inside the given path parsed")
	Cmd.PersistentFlags().StringVar(&global.LogLevel, "level", "info", "Logging level")
}

func run(cmd *cobra.Command, args []string) error {
	err := config.Read(cmd, global)
	if err != nil {
		return err
	}

	options := []constructor.Option{
		maestro.WithLogLevel(logger.Global, global.LogLevel),
		maestro.WithCodec(json.NewConstructor()),
		maestro.WithCodec(proto.NewConstructor()),
		maestro.WithCaller(micro.New("micro-grpc", microGRPC.NewService())),
		maestro.WithCaller(grpc.NewCaller()),
		maestro.WithCaller(http.NewCaller()),
	}

	for _, path := range global.HCL {
		options = append(options, maestro.WithFlows(hcl.FlowsResolver(path)))
		options = append(options, maestro.WithServices(hcl.ServicesResolver(path)))
		options = append(options, maestro.WithEndpoints(hcl.EndpointsResolver(path)))
	}

	for _, path := range global.Protobuffers {
		options = append(options, maestro.WithSchema(protoc.SchemaResolver(global.Protobuffers, path)))
		options = append(options, maestro.WithServices(protoc.ServiceResolver(global.Protobuffers, path)))
	}

	if global.HTTP.Address != "" {
		options = append(options, maestro.WithListener(http.NewListener(global.HTTP.Address, specs.Options{})))
	}

	if global.GraphQL.Address != "" {
		options = append(options, maestro.WithListener(graphql.NewListener(global.GraphQL.Address, specs.Options{})))
	}

	if global.GRPC.Address != "" {
		options = append(options, maestro.WithListener(grpc.NewListener(global.GRPC.Address, specs.Options{})))
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
