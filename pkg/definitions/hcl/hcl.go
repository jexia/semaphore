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
	return func(ctx instance.Context) (*specs.ServicesManifest, error) {
		files, err := utils.ResolvePath(path)
		if err != nil {
			return nil, err
		}

		services := &specs.ServicesManifest{}

		for _, file := range files {
			reader, err := os.Open(file.Path)
			if err != nil {
				return nil, err
			}

			definition, err := UnmarshalHCL(ctx, file.Name(), reader)
			if err != nil {
				return nil, err
			}

			collection, err := ParseServices(ctx, definition)
			if err != nil {
				return nil, err
			}

			services.Merge(collection)
		}

		return services, nil
	}
}

// FlowsResolver constructs a resource resolver for the given path
func FlowsResolver(path string) definitions.FlowsResolver {
	return func(ctx instance.Context) (*specs.FlowsManifest, error) {
		files, err := utils.ResolvePath(path)
		if err != nil {
			return nil, err
		}

		flows := &specs.FlowsManifest{}

		for _, file := range files {
			reader, err := os.Open(file.Path)
			if err != nil {
				return nil, err
			}

			definition, err := UnmarshalHCL(ctx, file.Name(), reader)
			if err != nil {
				return nil, err
			}

			manifest, err := ParseFlows(ctx, definition)
			if err != nil {
				return nil, err
			}

			flows.Merge(manifest)
		}

		return flows, nil
	}
}

// EndpointsResolver constructs a resource resolver for the given path
func EndpointsResolver(path string) definitions.EndpointsResolver {
	return func(ctx instance.Context) (*specs.EndpointsManifest, error) {
		files, err := utils.ResolvePath(path)
		if err != nil {
			return nil, err
		}

		endpoints := &specs.EndpointsManifest{}

		for _, file := range files {
			reader, err := os.Open(file.Path)
			if err != nil {
				return nil, err
			}

			definition, err := UnmarshalHCL(ctx, file.Name(), reader)
			if err != nil {
				return nil, err
			}

			manifest, err := ParseEndpoints(ctx, definition)
			if err != nil {
				return nil, err
			}

			endpoints.Merge(manifest)
		}

		return endpoints, nil
	}
}

// UnmarshalHCL unmarshals the given HCL stream into a intermediate resource.
func UnmarshalHCL(ctx instance.Context, filename string, reader io.Reader) (manifest Manifest, _ error) {
	ctx.Logger(logger.Core).WithField("file", filename).Info("Reading HCL files")

	bb, err := ioutil.ReadAll(reader)
	if err != nil {
		return manifest, err
	}

	ctx.Logger(logger.Core).WithField("file", filename).Debug("Parsing HCL syntax")

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
