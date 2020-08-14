package protobuf

import (
	"fmt"
	"os"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/cmd/semaphore/config"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/endpoints"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/providers"
	"github.com/jexia/semaphore/pkg/codec/proto"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/transport/grpc"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/spf13/cobra"
)

var params = config.New()

// Command represents the semaphore daemon command
var Command = &cobra.Command{
	Use:   "protobuf",
	Short: "Generates a protobuf definitions",
	Long: `Generates a protobuf definitions.
*NOTE*: This is a experimental feature and not all features are supported yet.`,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	Command.PersistentFlags().StringSliceVar(&params.Protobuffers, "proto", []string{}, "If set are all proto definitions found inside the given path passed as schema definitions, all proto definitions are also passed as imports")
	Command.PersistentFlags().StringSliceVarP(&params.Files, "file", "f", []string{"config.hcl"}, "Parses the given file as a definition file")
	Command.PersistentFlags().StringVar(&params.LogLevel, "level", "warn", "Global logging level, this value will override the defined log level inside the file definitions")
}

func run(cmd *cobra.Command, args []string) error {
	arguments, err := config.ConstructArguments(params)
	if err != nil {
		return err
	}

	ctx := logger.WithLogger(broker.NewContext())
	options, err := semaphore.NewOptions(ctx, arguments...)
	if err != nil {
		return err
	}

	collection, err := providers.Resolve(ctx, functions.Collection{}, options)
	if err != nil {
		return err
	}

	transporters, err := endpoints.Transporters(ctx, collection.EndpointList, collection.FlowListInterface,
		endpoints.WithServices(collection.ServiceList),
		endpoints.WithOptions(options),
		endpoints.WithFunctions(functions.Collection{}),
	)
	if err != nil {
		return err
	}

	var (
		services = make(map[string]*proto.Service)
		// descriptors = make(map[string]*desc.FileDescriptor)
	)

	for _, endpoint := range transporters {
		options, err := grpc.ParseEndpointOptions(endpoint)
		if err != nil {
			return err
		}

		protoService := proto.Service{
			Package: options.Package,
			Name:    options.Service,
			Methods: make(proto.Methods),
		}

		service, ok := services[protoService.String()]
		if !ok {
			service, services[protoService.String()] = &protoService, &protoService
		}

		method := &grpc.Method{
			Service:  service,
			Endpoint: endpoint,
			Name:     options.Method,
			Flow:     endpoint.Flow,
		}

		if err := method.NewCodec(ctx, proto.NewConstructor()); err != nil {
			return err
		}

		services[service.String()].Methods[method.String()] = method
	}

	printer := &protoprint.Printer{}

	for key, service := range services {
		descriptor, err := service.FileDescriptor()
		if err != nil {
			return fmt.Errorf("cannot generate file descriptor for %q: %w", key, err)
		}

		if err := printer.PrintProtoFile(descriptor, os.Stdout); err != nil {
			return err
		}
	}

	return nil
}
