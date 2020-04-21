package config

import (
	"github.com/jexia/maestro"
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
)

// ConstructArguments constructs the option arguments from the given parameters
func ConstructArguments(params *Maestro) ([]constructor.Option, error) {
	arguments := []constructor.Option{
		maestro.WithCodec(json.NewConstructor()),
		maestro.WithCodec(proto.NewConstructor()),
		maestro.WithCaller(micro.New("micro-grpc", microGRPC.NewService())),
		maestro.WithCaller(grpc.NewCaller()),
		maestro.WithCaller(http.NewCaller()),
	}

	for _, path := range params.Files {
		arguments = append(arguments, maestro.WithFlows(hcl.FlowsResolver(path)))
		arguments = append(arguments, maestro.WithServices(hcl.ServicesResolver(path)))
		arguments = append(arguments, maestro.WithEndpoints(hcl.EndpointsResolver(path)))

		options, err := hcl.GetOptions(path)
		if err != nil {
			return nil, err
		}

		Parse(options, params)
	}

	for _, path := range params.Protobuffers {
		arguments = append(arguments, maestro.WithSchema(protoc.SchemaResolver(params.Protobuffers, path)))
		arguments = append(arguments, maestro.WithServices(protoc.ServiceResolver(params.Protobuffers, path)))
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

	arguments = append(arguments, maestro.WithLogLevel(logger.Global, params.LogLevel))

	return arguments, nil
}
