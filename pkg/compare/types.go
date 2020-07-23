package compare

import (
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/jexia/semaphore/pkg/core/trace"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/sirupsen/logrus"
)

// ManifestTypes compares the types defined insde the schema definitions against the configured specification
func ManifestTypes(ctx instance.Context, services *specs.ServicesManifest, objects specs.Objects, flows *specs.FlowsManifest) (err error) {
	ctx.Logger(logger.Core).Info("Comparing manifest types")

	for _, flow := range flows.Flows {
		err := FlowTypes(ctx, services, objects, flows, flow)
		if err != nil {
			return err
		}
	}

	for _, proxy := range flows.Proxy {
		err := ProxyTypes(ctx, services, objects, flows, proxy)
		if err != nil {
			return err
		}
	}

	return nil
}

// ProxyTypes compares the given proxy against the configured schema types
func ProxyTypes(ctx instance.Context, services *specs.ServicesManifest, objects specs.Objects, flows *specs.FlowsManifest, proxy *specs.Proxy) (err error) {
	ctx.Logger(logger.Core).WithField("proxy", proxy.GetName()).Info("Compare proxy flow types")

	if proxy.OnError != nil {
		err = CheckParameterMapTypes(ctx, proxy.OnError.Response, objects, proxy)
		if err != nil {
			return err
		}
	}

	for _, node := range proxy.Nodes {
		err = CallTypes(ctx, services, objects, flows, node, node.Call, proxy)
		if err != nil {
			return err
		}

		err = CallTypes(ctx, services, objects, flows, node, node.Rollback, proxy)
		if err != nil {
			return err
		}
	}

	if proxy.Forward.Request.Header != nil {
		err = CheckHeader(proxy.Forward.Request.Header, proxy)
		if err != nil {
			return err
		}
	}

	return nil
}

// FlowTypes compares the flow types against the configured schema types
func FlowTypes(ctx instance.Context, services *specs.ServicesManifest, objects specs.Objects, flows *specs.FlowsManifest, flow *specs.Flow) (err error) {
	ctx.Logger(logger.Core).WithField("flow", flow.GetName()).Info("Comparing flow types")

	err = CheckParameterMapTypes(ctx, flow.Input, objects, flow)
	if err != nil {
		return err
	}

	if flow.OnError != nil {
		err = CheckParameterMapTypes(ctx, flow.OnError.Response, objects, flow)
		if err != nil {
			return err
		}
	}

	for _, node := range flow.Nodes {
		err = CallTypes(ctx, services, objects, flows, node, node.Call, flow)
		if err != nil {
			return err
		}

		err = CallTypes(ctx, services, objects, flows, node, node.Rollback, flow)
		if err != nil {
			return err
		}
	}

	if flow.Output != nil {
		message := objects.Get(flow.Output.Schema)
		if message == nil {
			return trace.New(trace.WithMessage("undefined flow output object '%s' in '%s'", flow.Output.Schema, flow.Name))
		}

		err = CheckParameterMapTypes(ctx, flow.Output, objects, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// CallTypes compares the given call types against the configured schema types
func CallTypes(ctx instance.Context, services *specs.ServicesManifest, objects specs.Objects, flows *specs.FlowsManifest, node *specs.Node, call *specs.Call, flow specs.FlowInterface) (err error) {
	if call == nil {
		return nil
	}

	if call.Method == "" {
		return nil
	}

	ctx.Logger(logger.Core).WithFields(logrus.Fields{
		"call":    node.ID,
		"method":  call.Method,
		"service": call.Service,
	}).Info("Comparing call types")

	service := services.GetService(call.Service)
	if service == nil {
		return trace.New(trace.WithMessage("undefined service '%s' in flow '%s'", call.Service, flow.GetName()))
	}

	method := service.GetMethod(call.Method)
	if method == nil {
		return trace.New(trace.WithMessage("undefined method '%s' in flow '%s'", call.Method, flow.GetName()))
	}

	err = CheckParameterMapTypes(ctx, call.Request, objects, flow)
	if err != nil {
		return err
	}

	err = CheckParameterMapTypes(ctx, call.Response, objects, flow)
	if err != nil {
		return err
	}

	if node.OnError != nil {
		err = CheckParameterMapTypes(ctx, node.OnError.Response, objects, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckParameterMapTypes checks the given parameter map against the configured schema property
func CheckParameterMapTypes(ctx instance.Context, parameters *specs.ParameterMap, objects specs.Objects, flow specs.FlowInterface) error {
	if parameters == nil {
		return nil
	}

	if parameters.Header != nil {
		err := CheckHeader(parameters.Header, flow)
		if err != nil {
			return err
		}
	}

	err := CheckPropertyTypes(parameters.Property, objects.Get(parameters.Schema), flow)
	if err != nil {
		return err
	}

	return nil
}

// CheckPropertyTypes checks the given schema against the given schema method types
func CheckPropertyTypes(property *specs.Property, schema *specs.Property, flow specs.FlowInterface) (err error) {
	if schema == nil {
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("unable to check types for '%s' no schema given", property.Path))
	}

	if property.Type != schema.Type {
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("cannot use type (%s) for '%s', expected (%s)", property.Type, property.Path, schema.Type))
	}

	if property.Label != schema.Label {
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("cannot use label (%s) for '%s', expected (%s)", property.Label, property.Path, schema.Label))
	}

	if len(property.Nested) > 0 {
		if len(schema.Nested) == 0 {
			return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("property '%s' has a nested object but schema does not '%s'", property.Path, schema.Name))
		}

		for key, nested := range property.Nested {
			object := schema.Nested[key]
			if object == nil {
				return trace.New(trace.WithExpression(nested.Expr), trace.WithMessage("undefined schema nested message property '%s' in flow '%s'", nested.Path, flow.GetName()))
			}

			err := CheckPropertyTypes(nested, object, flow)
			if err != nil {
				return err
			}
		}

		// Set any properties not defined inside the flow but available inside the schema
		for _, prop := range schema.Nested {
			_, has := property.Nested[prop.Name]
			if has {
				continue
			}

			property.Nested[prop.Name] = prop
		}
	}

	// ensure the property position
	property.Position = schema.Position

	return nil
}

// CheckHeader compares the given header types
func CheckHeader(header specs.Header, flow specs.FlowInterface) error {
	for _, header := range header {
		if header.Type != types.String {
			return trace.New(trace.WithMessage("cannot use type (%s) for 'header.%s' in flow '%s', expected (%s)", header.Type, header.Path, flow.GetName(), types.String))
		}
	}

	return nil
}
