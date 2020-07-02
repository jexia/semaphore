package config

import (
	"github.com/jexia/maestro"
	"github.com/jexia/maestro/cmd/maestro/middleware"
	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/core/api"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/core/logger"
	"github.com/jexia/maestro/pkg/metrics/prometheus"
	"github.com/jexia/maestro/pkg/providers/hcl"
	"github.com/jexia/maestro/pkg/providers/proto"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport/graphql"
	"github.com/jexia/maestro/pkg/transport/grpc"
	"github.com/jexia/maestro/pkg/transport/http"
	"github.com/jexia/maestro/pkg/transport/micro"
	microGRPC "github.com/micro/go-micro/v2/service/grpc"
)

// ConstructArguments constructs the option arguments from the given parameters
func ConstructArguments(params *Maestro) ([]api.Option, error) {
	arguments := []api.Option{
		maestro.WithCodec(codec.JSON()),
		maestro.WithCodec(codec.Proto()),
		maestro.WithCaller(micro.NewCaller("micro-grpc", microGRPC.NewService())),
		maestro.WithCaller(grpc.NewCaller()),
		maestro.WithCaller(http.NewCaller()),
	}

	ctx := instance.NewContext()
	ctx.SetLevel(logger.Global, params.LogLevel)

	for _, path := range params.Files {
		arguments = append(arguments, maestro.WithFlows(hcl.FlowsResolver(path)))
		arguments = append(arguments, maestro.WithServices(hcl.ServicesResolver(path)))
		arguments = append(arguments, maestro.WithEndpoints(hcl.EndpointsResolver(path)))
		arguments = append(arguments, maestro.AfterConstructor(middleware.ServiceSelector(path)))

		options, err := hcl.GetOptions(ctx, path)
		if err != nil {
			return nil, err
		}

		Parse(options, params)
	}

	for _, path := range params.Protobuffers {
		arguments = append(arguments, maestro.WithSchema(proto.SchemaResolver(params.Protobuffers, path)))
		arguments = append(arguments, maestro.WithServices(proto.ServiceResolver(params.Protobuffers, path)))
	}

	if params.HTTP.Address != "" {
		arguments = append(arguments, maestro.WithListener(http.NewListener(params.HTTP.Address, specs.Options{})))
	}

	if params.GraphQL.Address != "" {
		arguments = append(arguments, maestro.WithListener(graphql.NewListener(params.GraphQL.Address, specs.Options{})))
	}

	if params.GRPC.Address != "" {
		arguments = append(arguments, maestro.WithListener(grpc.NewListener(params.GRPC.Address, specs.Options{})))
	}

	if params.Prometheus.Address != "" {
		arguments = append(arguments, maestro.WithMiddleware(prometheus.New(params.Prometheus.Address)))
	}

	arguments = append([]api.Option{maestro.WithLogLevel(logger.Global, params.LogLevel)}, arguments...)

	return arguments, nil
}
