package hcl

import (
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
)

// ParseServices parses the given intermediate manifest to a schema
func ParseServices(ctx instance.Context, manifest Manifest) (specs.ServiceList, error) {
	ctx.Logger(logger.Core).Info("Parsing intermediate manifest to schema")

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
func ParseIntermediateService(ctx instance.Context, manifest Service) (*specs.Service, error) {
	ctx.Logger(logger.Core).WithField("service", manifest.Name).Debug("Parsing intermediate service to schema")

	methods, err := ParseIntermediateMethods(ctx, manifest.Methods)
	if err != nil {
		return nil, err
	}

	result := &specs.Service{
		Package:            manifest.Package,
		FullyQualifiedName: template.JoinPath(manifest.Package, manifest.Name),
		Name:               manifest.Name,
		Transport:          manifest.Transport,
		Host:               manifest.Host,
		Codec:              manifest.Codec,
		Methods:            methods,
		Options:            ParseIntermediateDefinitionOptions(manifest.Options),
	}

	return result, nil
}

// ParseIntermediateMethods parses the given methods for the given service
func ParseIntermediateMethods(ctx instance.Context, methods []Method) ([]*specs.Method, error) {
	result := make([]*specs.Method, len(methods))

	for index, manifest := range methods {
		ctx.Logger(logger.Core).WithFields(logrus.Fields{
			"method": manifest.Name,
		}).Debug("Parsing intermediate method to schema")

		result[index] = &specs.Method{
			Name:    manifest.Name,
			Input:   manifest.Input,
			Output:  manifest.Output,
			Options: ParseIntermediateDefinitionOptions(manifest.Options),
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
