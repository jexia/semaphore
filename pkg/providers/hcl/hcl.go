package hcl

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/specs"
	"go.uber.org/zap"
)

// ServicesResolver constructs a schema resolver for the given path.
// The HCL schema resolver relies on other schema registries.
// Those need to be resolved before the HCL schemas are resolved.
func ServicesResolver(path string) providers.ServicesResolver {
	return func(ctx *broker.Context) (specs.ServiceList, error) {
		logger.Debug(ctx, "resolving HCL services", zap.String("path", path))

		definitions, err := ResolvePath(ctx, []string{}, path)
		if err != nil {
			return nil, err
		}

		services := make(specs.ServiceList, 0)

		for _, definition := range definitions {
			result, err := ParseServices(ctx, definition)
			if err != nil {
				return nil, err
			}

			services.Append(result)
		}

		return services, nil
	}
}

// FlowsResolver constructs a resource resolver for the given path
func FlowsResolver(path string) providers.FlowsResolver {
	return func(ctx *broker.Context) (specs.FlowListInterface, error) {
		logger.Debug(ctx, "resolving HCL flows", zap.String("path", path))

		definitions, err := ResolvePath(ctx, []string{}, path)
		if err != nil {
			return nil, err
		}

		var errObject *specs.ParameterMap
		flows := make(specs.FlowListInterface, 0)

		for _, definition := range definitions {
			errResult, result, err := ParseFlows(ctx, definition)
			if err != nil {
				return nil, err
			}

			if errResult != nil {
				errObject = errResult
			}

			flows.Append(result)
		}

		ResolveErrors(flows, errObject)

		return flows, nil
	}
}

// EndpointsResolver constructs a resource resolver for the given path
func EndpointsResolver(path string) providers.EndpointsResolver {
	return func(ctx *broker.Context) (specs.EndpointList, error) {
		logger.Debug(ctx, "resolving HCL endpoints", zap.String("path", path))

		definitions, err := ResolvePath(ctx, []string{}, path)
		if err != nil {
			return nil, err
		}

		endpoints := make(specs.EndpointList, 0, len(definitions))

		for _, definition := range definitions {
			result, err := ParseEndpoints(ctx, definition)
			if err != nil {
				return nil, err
			}

			endpoints.Append(result)
		}

		return endpoints, nil
	}
}

// GetOptions returns the defined options inside the given path
func GetOptions(ctx *broker.Context, path string) (*Options, error) {
	definitions, err := ResolvePath(ctx, []string{}, path)
	if err != nil {
		return nil, err
	}

	options := &Options{}

	for _, definition := range definitions {
		if definition.LogLevel != "" {
			options.LogLevel = definition.LogLevel
		}

		if len(definition.Protobuffers) > 0 {
			options.Protobuffers = append(options.Protobuffers, definition.Protobuffers...)
		}

		if definition.GRPC != nil {
			options.GRPC = definition.GRPC
		}

		if definition.HTTP != nil {
			options.HTTP = definition.HTTP
		}

		if definition.GraphQL != nil {
			options.GraphQL = definition.GraphQL
		}

		if definition.Prometheus != nil {
			options.Prometheus = definition.Prometheus
		}
	}

	return options, nil
}

// ResolvePath resolves the given path and returns the available manifests.
// All defined includes are followed and their manifests are included
func ResolvePath(ctx *broker.Context, ignore []string, path string) ([]Manifest, error) {
	logger.Debug(ctx, "resolving HCL path", zap.String("path", path))

	manifests := make([]Manifest, 0)
	if path == "" {
		return manifests, nil
	}

	files, err := providers.ResolvePath(ctx, ignore, path)
	ignore = append(ignore, path)

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrPathNotFound{
			Path: path,
		}
	}

	logger.Debug(ctx, "files found", zap.String("path", path), zap.Int("files", len(files)))

	for _, file := range files {
		logger.Debug(ctx, "resolving file", zap.String("path", file.Path))

		reader, err := os.Open(file.Path)
		if err != nil {
			return nil, err
		}

		definition, err := UnmarshalHCL(ctx, file.Name(), reader)
		if err != nil {
			return nil, err
		}

		if definition.Protobuffers != nil {
			for index, proto := range definition.Protobuffers {
				if !filepath.IsAbs(proto) {
					proto = filepath.Join(filepath.Dir(file.Path), proto)
				}

				definition.Protobuffers[index] = proto
			}
		}

		manifests = append(manifests, definition)

		for _, include := range definition.Include {
			path := include

			if !filepath.IsAbs(include) {
				path = filepath.Join(filepath.Dir(file.Path), include)
			}

			logger.Info(ctx, "including HCL file", zap.String("path", path))
			results, err := ResolvePath(ctx, ignore, path)
			if err != nil {
				return nil, ErrPathNotFound{
					wrapErr: wrapErr{err},
					Path:    include,
				}
			}

			manifests = append(manifests, results...)
		}
	}

	logger.Debug(ctx, "resolve path result", zap.String("path", path), zap.Int("manifests", len(files)))
	return manifests, nil
}

// UnmarshalHCL unmarshals the given HCL stream into a intermediate resource.
// All environment variables found inside the given file are replaced.
func UnmarshalHCL(parent *broker.Context, filename string, reader io.Reader) (manifest Manifest, _ error) {
	ctx := logger.WithFields(parent, zap.String("file", filename))
	logger.Info(ctx, "reading HCL files")

	bb, err := ioutil.ReadAll(reader)
	if err != nil {
		return manifest, err
	}

	logger.Info(ctx, "parsing HCL syntax")

	// replace all environment variables found inside the given file
	bb = []byte(os.ExpandEnv(string(bb)))

	file, diags := hclsyntax.ParseConfig(bb, filename, hcl.InitialPos)
	if diags.HasErrors() {
		return manifest, errors.New(diags.Error())
	}

	logger.Debug(ctx, "decoding HCL syntax")

	diags = gohcl.DecodeBody(file.Body, nil, &manifest)
	if diags.HasErrors() {
		return manifest, errors.New(diags.Error())
	}

	return manifest, nil
}
