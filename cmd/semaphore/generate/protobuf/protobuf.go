package protobuf

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/cmd/semaphore/config"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/endpoints"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/providers"
	"github.com/jexia/semaphore/pkg/codec/proto"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/transport"
	"github.com/jexia/semaphore/pkg/transport/grpc"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/spf13/cobra"
)

var (
	params = config.New()

	// Command represents the semaphore daemon command
	Command = &cobra.Command{
		Use:   "protobuf",
		Short: "Generates a protobuf definitions",
		Long: `Generates a protobuf definitions.
	*NOTE*: This is a experimental feature and not all features are supported yet.`,
		RunE:         run,
		SilenceUsage: true,
	}

	output string
)

func init() {
	Command.PersistentFlags().StringSliceVar(&params.Protobuffers, "proto", []string{}, "If set are all proto definitions found inside the given path passed as schema definitions, all proto definitions are also passed as imports")
	Command.PersistentFlags().StringSliceVarP(&params.Files, "file", "f", []string{"config.hcl"}, "Parses the given file as a definition file")
	Command.PersistentFlags().StringVar(&params.LogLevel, "level", "warn", "Global logging level, this value will override the defined log level inside the file definitions")
	Command.PersistentFlags().StringVarP(&output, "output", "o", "", "Output directory (all missing directories will be created)")
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

	services, err := generate(ctx, transporters)
	if err != nil {
		return err
	}

	// TODO: configure printer with additional options
	printer := &protoprint.Printer{}

	for key, service := range services {

		dst, err := getOutput(output, key)
		if err != nil {
			return fmt.Errorf("failed to set the output for generator: %w", err)
		}

		descriptor, err := service.FileDescriptor()
		if err != nil {
			return fmt.Errorf("cannot generate file descriptor for %q: %w", key, err)
		}

		if err := printer.PrintProtoFile(descriptor, dst); err != nil {
			return err
		}

		dst.Close()
	}

	return nil
}

func getOutput(output, pkg string) (io.WriteCloser, error) {
	if output == "" {
		return os.Stdout, nil
	}

	filePath := path.Join(append([]string{output}, strings.Split(pkg, ".")...)...) + ".proto"

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(filePath)
}

func generate(ctx *broker.Context, endpoints transport.EndpointList) (map[string]*proto.Service, error) {
	var services = make(map[string]*proto.Service)

	for _, endpoint := range endpoints {
		if endpoint.Listener != "grpc" {
			continue
		}

		options, err := grpc.ParseEndpointOptions(endpoint)
		if err != nil {
			return nil, err
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
			return nil, err
		}

		services[service.String()].Methods[method.String()] = method
	}

	return services, nil
}
