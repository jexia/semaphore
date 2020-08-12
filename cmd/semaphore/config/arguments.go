package config

import (
	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/cmd/semaphore/middleware"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/config"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/codec/proto"
	"github.com/jexia/semaphore/pkg/metrics/prometheus"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/providers/protobuffers"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport/graphql"
	"github.com/jexia/semaphore/pkg/transport/grpc"
	"github.com/jexia/semaphore/pkg/transport/http"
	"github.com/jexia/semaphore/pkg/transport/micro"
	microGRPC "github.com/micro/go-micro/v2/service/grpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ConstructArguments constructs the option arguments from the given parameters
func ConstructArguments(params *Semaphore) ([]config.Option, error) {
	arguments := []config.Option{
		semaphore.WithCodec(json.NewConstructor()),
		semaphore.WithCodec(proto.NewConstructor()),
		semaphore.WithCaller(micro.NewCaller("micro-grpc", microGRPC.NewService())),
		semaphore.WithCaller(grpc.NewCaller()),
		semaphore.WithCaller(http.NewCaller()),
	}

	ctx := logger.WithLogger(broker.NewContext())

	level := zapcore.InfoLevel
	err := level.UnmarshalText([]byte(params.LogLevel))
	if err != nil {
		logger.Error(ctx, "unable to unmarshal log level", zap.String("level", params.LogLevel))
	}

	logger.SetLevel(ctx, "*", level)

	for _, path := range params.Files {
		arguments = append(arguments, semaphore.WithFlows(hcl.FlowsResolver(path)))
		arguments = append(arguments, semaphore.WithServices(hcl.ServicesResolver(path)))
		arguments = append(arguments, semaphore.WithEndpoints(hcl.EndpointsResolver(path)))
		arguments = append(arguments, semaphore.AfterConstructor(middleware.ServiceSelector(path)))

		options, err := hcl.GetOptions(ctx, path)
		if err != nil {
			return nil, err
		}

		Parse(options, params)
	}

	for _, path := range params.Protobuffers {
		arguments = append(arguments, semaphore.WithSchema(protobuffers.SchemaResolver(params.Protobuffers, path)))
		arguments = append(arguments, semaphore.WithServices(protobuffers.ServiceResolver(params.Protobuffers, path)))
	}

	if params.HTTP.Address != "" {
		arguments = append(
			arguments,
			semaphore.WithListener(
				http.NewListener(
					params.HTTP.Address,
					http.WithOrigins(params.HTTP.Origin),
					http.WithReadTimeout(params.HTTP.ReadTimeout),
					http.WithWriteTimeout(params.HTTP.WriteTimeout),
					http.WithKeyFile(params.HTTP.KeyFile),
					http.WithCertFile(params.HTTP.CertFile),
				),
			),
		)
	}

	if params.GraphQL.Address != "" {
		arguments = append(arguments, semaphore.WithListener(graphql.NewListener(params.GraphQL.Address, specs.Options{})))
	}

	if params.GRPC.Address != "" {
		arguments = append(arguments, semaphore.WithListener(grpc.NewListener(params.GRPC.Address, specs.Options{})))
	}

	if params.Prometheus.Address != "" {
		arguments = append(arguments, semaphore.WithMiddleware(prometheus.New(params.Prometheus.Address)))
	}

	arguments = append([]config.Option{semaphore.WithLogLevel("*", params.LogLevel)}, arguments...)

	return arguments, nil
}
