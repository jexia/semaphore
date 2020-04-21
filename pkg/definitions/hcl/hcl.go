package hcl

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/internal/utils"
	"github.com/jexia/maestro/pkg/definitions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/specs"
)

// ServicesResolver constructs a schema resolver for the given path.
// The HCL schema resolver relies on other schema registries.
// Those need to be resolved before the HCL schemas are resolved.
func ServicesResolver(path string) definitions.ServicesResolver {
	return func(ctx instance.Context) ([]*specs.ServicesManifest, error) {
		definitions, err := ResolvePath(ctx, path)
		if err != nil {
			return nil, err
		}

		services := make([]*specs.ServicesManifest, len(definitions))

		for index, definition := range definitions {
			manifest, err := ParseServices(ctx, definition)
			if err != nil {
				return nil, err
			}

			services[index] = manifest
		}

		return services, nil
	}
}

// FlowsResolver constructs a resource resolver for the given path
func FlowsResolver(path string) definitions.FlowsResolver {
	return func(ctx instance.Context) ([]*specs.FlowsManifest, error) {
		definitions, err := ResolvePath(ctx, path)
		if err != nil {
			return nil, err
		}

		flows := make([]*specs.FlowsManifest, len(definitions))

		for index, definition := range definitions {
			manifest, err := ParseFlows(ctx, definition)
			if err != nil {
				return nil, err
			}

			flows[index] = manifest
		}

		return flows, nil
	}
}

// EndpointsResolver constructs a resource resolver for the given path
func EndpointsResolver(path string) definitions.EndpointsResolver {
	return func(ctx instance.Context) ([]*specs.EndpointsManifest, error) {
		definitions, err := ResolvePath(ctx, path)
		if err != nil {
			return nil, err
		}

		endpoints := make([]*specs.EndpointsManifest, len(definitions))

		for index, definition := range definitions {
			manifest, err := ParseEndpoints(ctx, definition)
			if err != nil {
				return nil, err
			}

			endpoints[index] = manifest
		}

		return endpoints, nil
	}
}

// GetOptions returns the defined options inside the given path
func GetOptions(path string) (*Options, error) {
	definitions, err := ResolvePath(instance.NewContext(), path)
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
	}

	return options, nil
}

// Resolve represents a resolve object
type Resolve struct {
	File     *utils.FileInfo
	Manifest Manifest
	Err      error
}

// ResolvePath resolves the given path and returns the available manifests.
// All defined includes are followed and their manifests are included
func ResolvePath(ctx instance.Context, path string) ([]Manifest, error) {
	files, err := utils.ResolvePath(path)
	if err != nil {
		return nil, err
	}

	definitions := make([]Manifest, 0)

	for _, file := range files {
		reader, err := os.Open(file.Path)
		if err != nil {
			return nil, err
		}

		definition, err := UnmarshalHCL(ctx, file.Name(), reader)
		if err != nil {
			return nil, err
		}

		definitions = append(definitions, definition)

		for _, include := range definition.Include {
			manifests, err := ResolvePath(ctx, include)
			if err != nil {
				return nil, err
			}

			definitions = append(definitions, manifests...)
		}
	}

	return definitions, nil
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
