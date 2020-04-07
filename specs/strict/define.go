package strict

import (
	"github.com/jexia/maestro/internal/instance"
	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/lookup"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/transport"
	"github.com/sirupsen/logrus"
)

// DefineManifest checks and defines the types for the given manifest
func DefineManifest(ctx instance.Context, schema schema.Collection, manifest *specs.Manifest) (err error) {
	ctx.Logger(logger.Core).Info("Defining manifest types")

	for _, flow := range manifest.Flows {
		err := DefineFlow(ctx, schema, manifest, flow)
		if err != nil {
			return err
		}
	}

	for _, proxy := range manifest.Proxy {
		err := DefineProxy(ctx, schema, manifest, proxy)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineProxy checks and defines the types for the given proxy
func DefineProxy(ctx instance.Context, schema schema.Collection, manifest *specs.Manifest, proxy *specs.Proxy) (err error) {
	ctx.Logger(logger.Core).WithField("proxy", proxy.GetName()).Info("Defining proxy flow types")

	for _, node := range proxy.Nodes {
		if node.Call != nil {
			err = DefineCall(ctx, schema, manifest, node, node.Call, proxy)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = DefineCall(ctx, schema, manifest, node, node.Rollback, proxy)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DefineFlow defines the types for the given flow
func DefineFlow(ctx instance.Context, schema schema.Collection, manifest *specs.Manifest, flow *specs.Flow) (err error) {
	ctx.Logger(logger.Core).WithField("flow", flow.GetName()).Info("Defining flow types")

	if flow.Input != nil {
		message, err := GetObjectSchema(schema, flow.Input)
		if err != nil {
			return err
		}

		flow.Input = specs.ToParameterMap(flow.Input, "", message)
	}

	for _, node := range flow.Nodes {
		if node.Call != nil {
			err = DefineCall(ctx, schema, manifest, node, node.Call, flow)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = DefineCall(ctx, schema, manifest, node, node.Rollback, flow)
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

// GetObjectSchema attempts to fetch the defined schema object for the given parameter map
func GetObjectSchema(schema schema.Collection, params *specs.ParameterMap) (schema.Property, error) {
	prop := schema.GetMessage(params.Schema)
	if prop == nil {
		return nil, trace.New(trace.WithMessage("undefined object '%s' in schema collection", params.Schema))
	}

	return prop, nil
}

// DefineCall defineds the types for the specs call
func DefineCall(ctx instance.Context, schema schema.Collection, manifest *specs.Manifest, node *specs.Node, call *specs.Call, flow specs.FlowManager) (err error) {
	if call.Request != nil {
		err = DefineParameterMap(ctx, node, call.Request, flow)
		if err != nil {
			return err
		}

		err = DefineFunctions(ctx, call.Request.Functions, node, flow)
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

		service := schema.GetService(call.Service)
		if service == nil {
			return trace.New(trace.WithMessage("undefined service '%s' in flow '%s'", call.Service, flow.GetName()))
		}

		method := service.GetMethod(call.Method)
		if method == nil {
			return trace.New(trace.WithMessage("undefined method '%s' in flow '%s'", call.Method, flow.GetName()))
		}

		method.GetInput()
		call.SetDescriptor(method)
		call.SetResponse(specs.ToParameterMap(nil, "", method.GetOutput()))
	}

	if call.Response != nil {
		err = DefineParameterMap(ctx, node, call.Response, flow)
		if err != nil {
			return err
		}

		err = DefineFunctions(ctx, call.Response.Functions, node, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineFunctions defined all properties within the given functions
func DefineFunctions(ctx instance.Context, functions specs.Functions, node *specs.Node, flow specs.FlowManager) error {
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
func DefineCaller(ctx instance.Context, node *specs.Node, manifest *specs.Manifest, call transport.Call, flow specs.FlowManager) (err error) {
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
func DefineParameterMap(ctx instance.Context, node *specs.Node, params *specs.ParameterMap, flow specs.FlowManager) (err error) {
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
func DefineProperty(ctx instance.Context, node *specs.Node, property *specs.Property, flow specs.FlowManager) error {
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

	breakpoint := specs.OutputResource
	if node != nil {
		breakpoint = node.Name

		if node.Rollback != nil {
			rollback := node.Rollback.Request.Property
			if InsideProperty(rollback, property) {
				breakpoint = lookup.GetNextResource(flow, breakpoint)
			}
		}
	}

	ctx.Logger(logger.Core).WithFields(logrus.Fields{
		"breakpoint": breakpoint,
		"reference":  property.Reference,
	}).Debug("Lookup references until breakpoint")

	references := lookup.GetAvailableResources(flow, breakpoint)
	reference := lookup.GetResourceReference(property.Reference, references, breakpoint)
	if reference == nil {
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("undefined resource '%s' in '%s.%s.%s'", property.Reference, flow.GetName(), breakpoint, property.Path))
	}

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

// SchemaToProperty parses the given schema property to a specs property
func SchemaToProperty(path string, prop schema.Property) *specs.Property {
	result := &specs.Property{
		Name:      prop.GetName(),
		Path:      specs.JoinPath(path, prop.GetName()),
		Type:      prop.GetType(),
		Label:     prop.GetLabel(),
		Desciptor: prop,
	}

	if len(prop.GetNested()) > 0 {
		result.Nested = make(map[string]*specs.Property, len(prop.GetNested()))

		for key, prop := range prop.GetNested() {
			result.Nested[key] = SchemaToProperty(specs.JoinPath(path, key), prop)
		}
	}

	return result
}
