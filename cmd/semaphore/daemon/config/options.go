package config

import (
	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/cmd/semaphore/daemon/providers"
	"github.com/jexia/semaphore/cmd/semaphore/functions"
	"github.com/jexia/semaphore/cmd/semaphore/middleware"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/codec/proto"
	formencoded "github.com/jexia/semaphore/pkg/codec/www-form-urlencoded"
	"github.com/jexia/semaphore/pkg/codec/xml"
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

// Daemon configurations
type Daemon struct {
	LogLevel     string
	HTTP         HTTP
	GraphQL      GraphQL
	GRPC         GRPC
	Prometheus   Prometheus
	Protobuffers []string
	Files        []string
}

// Prometheus configurations
type Prometheus struct {
	Address string
}

// HTTP configurations
type HTTP struct {
	Address      string
	Origin       []string
	CertFile     string
	KeyFile      string
	ReadTimeout  string
	WriteTimeout string
}

// GRPC configurations
type GRPC struct {
	Address string
}

// GraphQL configurations
type GraphQL struct {
	Address string
}

// SetOptions reads the available configuration files and attempts to set the
// read options.
func SetOptions(ctx *broker.Context, flags *Daemon) error {
	for _, path := range flags.Files {
		options, err := hcl.GetOptions(ctx, path)
		if err != nil {
			return err
		}

		parseHCL(options, flags)
	}

	level := zapcore.InfoLevel
	err := level.UnmarshalText([]byte(flags.LogLevel))
	if err != nil {
		logger.Error(ctx, "unable to unmarshal log level", zap.String("level", flags.LogLevel))
	}

	err = logger.SetLevel(ctx, "*", level)
	if err != nil {
		logger.Error(ctx, "unable to set log level", zap.Error(err))
	}

	return nil
}

func parseHCL(options *hcl.Options, target *Daemon) {
	if target.LogLevel == "" && options.LogLevel != "" {
		target.LogLevel = options.LogLevel
	}

	if len(options.Protobuffers) > 0 {
		target.Protobuffers = append(target.Protobuffers, options.Protobuffers...)
	}

	if options.GraphQL != nil && target.GraphQL.Address == "" {
		target.GraphQL = GraphQL(*options.GraphQL)
	}

	if options.HTTP != nil && target.HTTP.Address == "" {
		target.HTTP = HTTP(*options.HTTP)
	}

	if options.GRPC != nil && target.GRPC.Address == "" {
		target.GRPC = GRPC(*options.GRPC)
	}

	if options.Prometheus != nil && target.Prometheus.Address == "" {
		target.Prometheus = Prometheus{
			Address: options.Prometheus.Address,
		}
	}
}

// NewCore constructs new core options from the given parameters
func NewCore(ctx *broker.Context, flags *Daemon) (semaphore.Options, error) {
	options := []semaphore.Option{
		semaphore.WithCodec(json.NewConstructor()),
		semaphore.WithCodec(proto.NewConstructor()),
		semaphore.WithCodec(formencoded.NewConstructor()),
		semaphore.WithCodec(xml.NewConstructor()),
		semaphore.WithCaller(micro.NewCaller("micro-grpc", microGRPC.NewService())),
		semaphore.WithCaller(grpc.NewCaller()),
		semaphore.WithCaller(http.NewCaller()),
		semaphore.WithFunctions(functions.Default),
	}

	for _, path := range flags.Files {
		options = append(options, semaphore.WithFlows(hcl.FlowsResolver(path)))
	}

	if flags.Prometheus.Address != "" {
		options = append(options, semaphore.WithMiddleware(prometheus.New(flags.Prometheus.Address)))
	}

	options = append([]semaphore.Option{semaphore.WithLogLevel("*", flags.LogLevel)}, options...)

	return semaphore.NewOptions(ctx, options...)
}

// NewProviders constructs new providers options from the given parameters
func NewProviders(ctx *broker.Context, core semaphore.Options, params *Daemon) (providers.Options, error) {
	options := []providers.Option{}

	for _, path := range params.Files {
		options = append(options, providers.WithServices(hcl.ServicesResolver(path)))
		options = append(options, providers.WithEndpoints(hcl.EndpointsResolver(path)))
		options = append(options, providers.WithAfterConstructor(middleware.ServiceSelector(path)))
	}

	for _, path := range params.Protobuffers {
		options = append(options, providers.WithSchema(protobuffers.SchemaResolver(params.Protobuffers, path)))
		options = append(options, providers.WithServices(protobuffers.ServiceResolver(params.Protobuffers, path)))
	}

	if params.HTTP.Address != "" {
		options = append(options,
			providers.WithListener(
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
		options = append(options, providers.WithListener(graphql.NewListener(params.GraphQL.Address, specs.Options{})))
	}

	if params.GRPC.Address != "" {
		options = append(options, providers.WithListener(grpc.NewListener(params.GRPC.Address, specs.Options{})))
	}

	return providers.NewOptions(ctx, core, options...)
}
