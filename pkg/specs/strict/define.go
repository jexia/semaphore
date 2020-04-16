package strict

import (
	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/lookup"
	"github.com/jexia/maestro/pkg/specs/template"
	"github.com/jexia/maestro/pkg/specs/trace"
	"github.com/jexia/maestro/pkg/transport"
	"github.com/sirupsen/logrus"
)

// DefineManifest checks and defines the types for the given manifest
func DefineManifest(ctx instance.Context, services *specs.ServicesManifest, schema *specs.SchemaManifest, flows *specs.FlowsManifest) (err error) {
	ctx.Logger(logger.Core).Info("Defining manifest types")

	for _, flow := range flows.Flows {
		err := DefineFlow(ctx, services, schema, flows, flow)
		if err != nil {
			return err
		}
	}

	for _, proxy := range flows.Proxy {
		err := DefineProxy(ctx, services, schema, flows, proxy)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineProxy checks and defines the types for the given proxy
func DefineProxy(ctx instance.Context, services *specs.ServicesManifest, schema *specs.SchemaManifest, flows *specs.FlowsManifest, proxy *specs.Proxy) (err error) {
	ctx.Logger(logger.Core).WithField("proxy", proxy.GetName()).Info("Defining proxy flow types")

	for _, node := range proxy.Nodes {
		if node.Call != nil {
			err = DefineCall(ctx, services, schema, flows, node, node.Call, proxy)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = DefineCall(ctx, services, schema, flows, node, node.Rollback, proxy)
			if err != nil {
				return err
			}
		}
	}

	if proxy.Forward != nil {
		for _, header := range proxy.Forward.Request.Header {
			err = DefineProperty(ctx, nil, header, proxy)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DefineFlow defines the types for the given flow
func DefineFlow(ctx instance.Context, services *specs.ServicesManifest, schema *specs.SchemaManifest, flows *specs.FlowsManifest, flow *specs.Flow) (err error) {
	ctx.Logger(logger.Core).WithField("flow", flow.GetName()).Info("Defining flow types")

	if flow.Input != nil {
		input := schema.GetProperty(flow.Input.Schema)
		if input == nil {
			return trace.New(trace.WithMessage("undefined object '%s' in schema collection", flow.Input.Schema))
		}

		flow.Input = ToParameterMap(flow.Input, "", input)
	}

	for _, node := range flow.Nodes {
		if node.Call != nil {
			err = DefineCall(ctx, services, schema, flows, node, node.Call, flow)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = DefineCall(ctx, services, schema, flows, node, node.Rollback, flow)
			if err != nil {
				return err
			}
		}
	}

	if flow.Output != nil {
		err = DefineParameterMap(ctx, nil, flow.Output, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineCall defineds the types for the specs call
func DefineCall(ctx instance.Context, services *specs.ServicesManifest, schema *specs.SchemaManifest, flows *specs.FlowsManifest, node *specs.Node, call *specs.Call, flow specs.FlowResourceManager) (err error) {
	if call.Request != nil {
		err = DefineParameterMap(ctx, node, call.Request, flow)
		if err != nil {
			return err
		}
	}

	if call.Method != "" {
		ctx.Logger(logger.Core).WithFields(logrus.Fields{
			"call":    node.Name,
			"method":  call.Method,
			"service": call.Service,
		}).Info("Defining call types")

		service := services.GetService(call.Service)
		if service == nil {
			return trace.New(trace.WithMessage("undefined service '%s' in flow '%s'", call.Service, flow.GetName()))
		}

		method := service.GetMethod(call.Method)
		if method == nil {
			return trace.New(trace.WithMessage("undefined method '%s' in flow '%s'", call.Method, flow.GetName()))
		}

		output := schema.GetProperty(method.Output)
		if output == nil {
			return trace.New(trace.WithMessage("undefined method output property '%s' in flow '%s'", method.Output, flow.GetName()))
		}

		call.Descriptor = method
		call.Response = ToParameterMap(nil, "", output)
	}

	if call.Response != nil {
		err = DefineParameterMap(ctx, node, call.Response, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineFunctions defined all properties within the given functions
func DefineFunctions(ctx instance.Context, functions functions.Stack, node *specs.Node, flow specs.FlowResourceManager) error {
	if functions == nil {
		return nil
	}

	for _, function := range functions {
		if function.Arguments != nil {
			for _, arg := range function.Arguments {
				DefineProperty(ctx, node, arg, flow)
			}
		}

		DefineProperty(ctx, node, function.Returns, flow)
	}

	return nil
}

// DefineCaller defineds the types for the given transport caller
func DefineCaller(ctx instance.Context, node *specs.Node, manifest *specs.FlowsManifest, call transport.Call, flow specs.FlowResourceManager) (err error) {
	ctx.Logger(logger.Core).Info("Defining caller references")

	method := call.GetMethod(node.Call.Method)
	for _, prop := range method.References() {
		err = DefineProperty(ctx, node, prop, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineParameterMap defines the types for the given parameter map
func DefineParameterMap(ctx instance.Context, node *specs.Node, params *specs.ParameterMap, flow specs.FlowResourceManager) (err error) {
	if params.Property == nil {
		return nil
	}

	for _, header := range params.Header {
		err = DefineProperty(ctx, node, header, flow)
		if err != nil {
			return err
		}
	}

	err = DefineProperty(ctx, node, params.Property, flow)
	if err != nil {
		return err
	}

	return nil
}

// DefineProperty defines the given property type.
// If any object is references it has to be fixed afterwards and moved into the correct dataset
func DefineProperty(ctx instance.Context, node *specs.Node, property *specs.Property, flow specs.FlowResourceManager) error {
	if len(property.Nested) > 0 {
		for _, nested := range property.Nested {
			err := DefineProperty(ctx, node, nested, flow)
			if err != nil {
				return err
			}
		}
	}

	if property.Reference == nil {
		return nil
	}

	breakpoint := template.OutputResource
	if node != nil {
		breakpoint = node.Name

		if node.Rollback != nil {
			rollback := node.Rollback.Request.Property
			if InsideProperty(rollback, property) {
				breakpoint = lookup.GetNextResource(flow, breakpoint)
			}
		}
	}

	// Self ref issue
	// lookup.ResolveSelfReference(node, property.Reference)

	ctx.Logger(logger.Core).WithFields(logrus.Fields{
		"breakpoint": breakpoint,
		"reference":  property.Reference,
	}).Debug("Lookup references until breakpoint")

	references := lookup.GetAvailableResources(flow, breakpoint)
	reference := lookup.GetResourceReference(property.Reference, references, breakpoint)
	if reference == nil {
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("undefined resource '%s' in '%s.%s.%s'", property.Reference, flow.GetName(), breakpoint, property.Path))
	}

	ctx.Logger(logger.Core).WithFields(logrus.Fields{
		"name": reference.Name,
		"path": reference.Path,
	}).Debug("References lookup result")

	property.Type = reference.Type
	property.Label = reference.Label
	property.Default = reference.Default
	property.Reference.Property = reference

	// TODO: support enum type

	return nil
}

// InsideProperty checks whether the given property is insde the source property
func InsideProperty(source *specs.Property, target *specs.Property) bool {
	if source == target {
		return true
	}

	if len(source.Nested) > 0 {
		for _, nested := range source.Nested {
			is := InsideProperty(nested, target)
			if is {
				return is
			}
		}
	}

	return false
}
