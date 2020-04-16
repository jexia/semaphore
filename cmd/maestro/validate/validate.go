package validate

import (
	"github.com/jexia/maestro"
	"github.com/jexia/maestro/cmd/maestro/config"
	"github.com/jexia/maestro/internal/constructor"
	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/pkg/codec/json"
	"github.com/jexia/maestro/pkg/codec/proto"
	"github.com/jexia/maestro/pkg/definitions/hcl"
	"github.com/jexia/maestro/pkg/definitions/protoc"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/transport/grpc"
	"github.com/jexia/maestro/pkg/transport/http"
	"github.com/jexia/maestro/pkg/transport/micro"
	microGRPC "github.com/micro/go-micro/service/grpc"
	"github.com/spf13/cobra"
)

var global = config.New()

// Cmd represents the maestro validate command
var Cmd = &cobra.Command{
	Use:          "validate",
	Short:        "Validate the flow definitions with the configured schema format(s)",
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	Cmd.PersistentFlags().StringP("config", "c", "", "Config file path")
	Cmd.PersistentFlags().StringSliceVar(&global.Protobuffers, "proto", []string{}, "If set are all proto definitions found inside the given path passed as schema definitions, all proto definitions are also passed as imports")
	Cmd.PersistentFlags().StringSliceVar(&global.Flows, "flow", []string{}, "If set are all flow definitions inside the given path passed as flow definitions")
	Cmd.PersistentFlags().StringVar(&global.LogLevel, "level", "error", "Logging level")
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
		maestro.WithCaller(http.NewCaller()),
		maestro.WithCaller(grpc.NewCaller()),
	}

	for _, flow := range global.Flows {
		options = append(options, maestro.WithFlows(hcl.FlowsResolver(flow)))
	}

	for _, path := range global.Protobuffers {
		options = append(options, maestro.WithSchema(protoc.SchemaResolver(global.Protobuffers, path)))
	}

	ctx := instance.NewContext()
	_, _, _, _, err = constructor.Specs(ctx, maestro.NewOptions(ctx, options...))
	if err != nil {
		return err
	}

	return nil
}
