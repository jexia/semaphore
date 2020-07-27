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
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/jexia/semaphore/pkg/core/trace"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/sirupsen/logrus"
)

// ServicesResolver constructs a schema resolver for the given path.
// The HCL schema resolver relies on other schema registries.
// Those need to be resolved before the HCL schemas are resolved.
func ServicesResolver(path string) providers.ServicesResolver {
	return func(ctx instance.Context) (specs.ServiceList, error) {
		ctx.Logger(logger.Core).WithField("path", path).Debug("Resolving HCL services")

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
	return func(ctx instance.Context) (specs.FlowListInterface, error) {
		ctx.Logger(logger.Core).WithField("path", path).Debug("Resolving HCL flows")

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

			flows.Append(result...)
		}

		ResolveErrors(flows, errObject)

		return flows, nil
	}
}

// EndpointsResolver constructs a resource resolver for the given path
func EndpointsResolver(path string) providers.EndpointsResolver {
	return func(ctx instance.Context) (specs.EndpointList, error) {
		ctx.Logger(logger.Core).WithField("path", path).Debug("Resolving HCL endpoints")

		definitions, err := ResolvePath(ctx, []string{}, path)
		if err != nil {
			return nil, err
		}

		endpoints := make(specs.EndpointList, len(definitions))

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
func GetOptions(ctx instance.Context, path string) (*Options, error) {
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

// Resolve represents a resolve object
type Resolve struct {
	File     *providers.FileInfo
	Manifest Manifest
	Err      error
}

// ResolvePath resolves the given path and returns the available manifests.
// All defined includes are followed and their manifests are included
func ResolvePath(ctx instance.Context, ignore []string, path string) ([]Manifest, error) {
	ctx.Logger(logger.Core).WithField("path", path).Info("Resolving HCL path")

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
		return nil, trace.New(trace.WithMessage("unable to resolve path, no files found '%s'", path))
	}

	ctx.Logger(logger.Core).WithFields(logrus.Fields{
		"path":  path,
		"files": len(files),
	}).Debug("Files found after resolving path")

	for _, file := range files {
		ctx.Logger(logger.Core).WithField("path", file.Path).Debug("Resolving file")

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

			ctx.Logger(logger.Core).WithField("path", path).Info("Including HCL file")
			results, err := ResolvePath(ctx, ignore, path)
			if err != nil {
				return nil, trace.New(trace.WithMessage("unable to read include %s: %s", include, err))
			}

			manifests = append(manifests, results...)
		}
	}

	ctx.Logger(logger.Core).WithFields(logrus.Fields{
		"path":      path,
		"manifests": len(files),
	}).Debug("Resolve path result")

	return manifests, nil
}

// UnmarshalHCL unmarshals the given HCL stream into a intermediate resource.
// All environment variables found inside the given file are replaced.
func UnmarshalHCL(ctx instance.Context, filename string, reader io.Reader) (manifest Manifest, _ error) {
	ctx.Logger(logger.Core).WithField("file", filename).Info("Reading HCL files")

	bb, err := ioutil.ReadAll(reader)
	if err != nil {
		return manifest, err
	}

	ctx.Logger(logger.Core).WithField("file", filename).Debug("Parsing HCL syntax")

	// replace all environment variables found inside the given file
	bb = []byte(os.ExpandEnv(string(bb)))

	file, diags := hclsyntax.ParseConfig(bb, filename, hcl.InitialPos)
	if diags.HasErrors() {
		return manifest, errors.New(diags.Error())
	}

	ctx.Logger(logger.Core).WithField("file", filename).Debug("Decoding HCL syntax")

	diags = gohcl.DecodeBody(file.Body, nil, &manifest)
	if diags.HasErrors() {
		return manifest, errors.New(diags.Error())
	}

	return manifest, nil
}
