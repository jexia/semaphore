package hcl

import (
	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"github.com/jexia/semaphore/v2/pkg/specs/template"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
)

// ParseServices parses the given intermediate manifest to a schema
func ParseServices(ctx *broker.Context, manifest Manifest) (specs.ServiceList, error) {
	logger.Info(ctx, "parsing intermediate manifest to schema")

	result := make(specs.ServiceList, len(manifest.Services))

	for index, intermediate := range manifest.Services {
		service, err := ParseIntermediateService(ctx, intermediate)
		if err != nil {
			return nil, err
		}
		result[index] = service
	}

	return result, nil
}

// ParseIntermediateService parses the given intermediate service to a specs service
func ParseIntermediateService(parent *broker.Context, manifest Service) (*specs.Service, error) {
	ctx := logger.WithFields(parent, zap.String("service", manifest.Name))
	logger.Debug(ctx, "parsing intermediate service to schema")

	methods, err := ParseIntermediateMethods(ctx, manifest.Methods)
	if err != nil {
		return nil, err
	}

	resolver := "dns"
	if manifest.Resolver != "" {
		resolver = manifest.Resolver
	}

	result := &specs.Service{
		Package:            manifest.Package,
		FullyQualifiedName: template.JoinPath(manifest.Package, manifest.Name),
		Name:               manifest.Name,
		Transport:          manifest.Transport,
		Host:               manifest.Host,
		RequestCodec:       manifest.Codec,
		ResponseCodec:      manifest.Codec,
		Methods:            methods,
		Options:            ParseIntermediateDefinitionOptions(manifest.Options),
		Resolver:           resolver,
	}

	return result, nil
}

// ParseIntermediateMethods parses the given methods for the given service
func ParseIntermediateMethods(ctx *broker.Context, methods []Method) ([]*specs.Method, error) {
	result := make([]*specs.Method, len(methods))

	for index, method := range methods {
		logger.Debug(ctx, "parsing intermediate method to schema", zap.String("method", method.Name))

		result[index] = &specs.Method{
			Name:    method.Name,
			Input:   method.Input,
			Output:  method.Output,
			Options: ParseIntermediateDefinitionOptions(method.Options),
		}
	}

	return result, nil
}

// ParseIntermediateDefinitionOptions parses the given intermediate options to a definitions options
func ParseIntermediateDefinitionOptions(options *BlockOptions) specs.Options {
	if options == nil {
		return specs.Options{}
	}

	result := specs.Options{}
	attrs, _ := options.Body.JustAttributes()

	for key, val := range attrs {
		val, _ := val.Expr.Value(nil)
		if val.Type() != cty.String {
			continue
		}

		result[key] = val.AsString()
	}

	return result
}
