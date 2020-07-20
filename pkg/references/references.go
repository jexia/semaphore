package references

import (
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/jexia/semaphore/pkg/core/trace"
	"github.com/jexia/semaphore/pkg/lookup"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
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

	if proxy.Input != nil && proxy.Input.Schema != "" {
		input := schema.GetProperty(proxy.Input.Schema)
		if input == nil {
			return trace.New(trace.WithMessage("object '%s', is unavailable inside the schema collection", proxy.Input.Schema))
		}

		proxy.Input = ToParameterMap(proxy.Input, "", input)
	}

	if proxy.OnError != nil {
		err = DefineOnError(ctx, schema, nil, proxy.OnError, proxy)
		if err != nil {
			return err
		}
	}

	for _, node := range proxy.Nodes {
		err = DefineNode(ctx, services, schema, flows, node, proxy)
		if err != nil {
			return err
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
			return trace.New(trace.WithMessage("object '%s', is unavailable inside the schema collection", flow.Input.Schema))
		}

		flow.Input = ToParameterMap(flow.Input, "", input)
	}

	if flow.OnError != nil {
		err = DefineOnError(ctx, schema, nil, flow.OnError, flow)
		if err != nil {
			return err
		}
	}

	for _, node := range flow.Nodes {
		err = DefineNode(ctx, services, schema, flows, node, flow)
		if err != nil {
			return err
		}
	}

	if flow.Output != nil {
		err = DefineParameterMap(ctx, schema, nil, flow.Output, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineNode defines all the references inside the given node
func DefineNode(ctx instance.Context, services *specs.ServicesManifest, schema *specs.SchemaManifest, flows *specs.FlowsManifest, node *specs.Node, flow specs.FlowResourceManager) (err error) {
	if node.Condition != nil {
		err = DefineParameterMap(ctx, schema, node, node.Condition.Params, flow)
		if err != nil {
			return err
		}
	}

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

	if node.OnError != nil {
		err = DefineOnError(ctx, schema, node, node.OnError, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineCall defineds the types for the specs call
func DefineCall(ctx instance.Context, services *specs.ServicesManifest, schema *specs.SchemaManifest, flows *specs.FlowsManifest, node *specs.Node, call *specs.Call, flow specs.FlowResourceManager) (err error) {
	if call.Request != nil {
		err = DefineParameterMap(ctx, schema, node, call.Request, flow)
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

		call.Request.Schema = method.Input
		call.Response.Schema = method.Output
	}

	if call.Response != nil {
		err = DefineParameterMap(ctx, schema, node, call.Response, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineParameterMap defines the types for the given parameter map
func DefineParameterMap(ctx instance.Context, schema *specs.SchemaManifest, node *specs.Node, params *specs.ParameterMap, flow specs.FlowResourceManager) (err error) {
	if params.Schema != "" {
		result := schema.GetProperty(params.Schema)
		if result == nil {
			return trace.New(trace.WithMessage("object '%s', is unavailable inside the schema collection", params.Schema))
		}
	}

	for _, header := range params.Header {
		err = DefineProperty(ctx, node, header, flow)
		if err != nil {
			return err
		}
	}

	if params.Params != nil {
		err = DefineParams(ctx, node, params.Params, flow)
		if err != nil {
			return err
		}
	}

	if params.Property != nil {
		err = DefineProperty(ctx, node, params.Property, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineParams defines all types inside the given params
func DefineParams(ctx instance.Context, node *specs.Node, params map[string]*specs.Property, flow specs.FlowResourceManager) error {
	for _, param := range params {
		if param.Reference == nil {
			continue
		}

		err := DefineProperty(ctx, node, param, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineProperty defines the given property type.
// If any object is references it has to be fixed afterwards and moved into the correct dataset
func DefineProperty(ctx instance.Context, node *specs.Node, property *specs.Property, flow specs.FlowResourceManager) error {
	if property == nil {
		return nil
	}

	if len(property.Nested) > 0 {
		for _, nested := range property.Nested {
			err := DefineProperty(ctx, node, nested, flow)
			if err != nil {
				return err
			}
		}
	}

	// ensure property references to be looked up once
	if property.Reference == nil {
		return nil
	}

	breakpoint := template.OutputResource
	if node != nil {
		breakpoint = node.Name

		if node.Rollback != nil && property != nil {
			rollback := node.Rollback.Request.Property
			if InsideProperty(rollback, property) {
				breakpoint = lookup.GetNextResource(flow, breakpoint)
			}
		}
	}

	reference, err := LookupReference(ctx, node, breakpoint, property.Reference, flow)
	if err != nil {
		ctx.Logger(logger.Core).WithField("err", err).Debug("Unable to lookup reference")
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("undefined reference '%s' in '%s.%s.%s'", property.Reference, flow.GetName(), breakpoint, property.Path))
	}

	ctx.Logger(logger.Core).WithFields(logrus.Fields{
		"reference": property.Reference,
		"name":      reference.Name,
		"path":      reference.Path,
	}).Debug("References lookup result")

	property.Type = reference.Type
	property.Label = reference.Label
	property.Default = reference.Default
	property.Reference.Property = reference

	if reference.Enum != nil {
		property.Enum = reference.Enum
	}

	return nil
}

// LookupReference looks up the given reference
func LookupReference(ctx instance.Context, node *specs.Node, breakpoint string, reference *specs.PropertyReference, flow specs.FlowResourceManager) (*specs.Property, error) {
	reference.Resource = lookup.ResolveSelfReference(reference.Resource, breakpoint)

	ctx.Logger(logger.Core).WithFields(logrus.Fields{
		"breakpoint": breakpoint,
		"reference":  reference,
	}).Debug("Lookup references until breakpoint")

	references := lookup.GetAvailableResources(flow, breakpoint)
	result := lookup.GetResourceReference(reference, references, breakpoint)
	if result == nil {
		return nil, trace.New(trace.WithMessage("undefined resource '%s' in '%s'.'%s'", reference, flow.GetName(), breakpoint))
	}

	ctx.Logger(logger.Core).WithFields(logrus.Fields{
		"breakpoint": breakpoint,
		"reference":  result,
	}).Debug("Lookup references result")

	return result, nil
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

// DefineOnError defines references made inside the given on error specs
func DefineOnError(ctx instance.Context, schema *specs.SchemaManifest, node *specs.Node, params *specs.OnError, flow specs.FlowResourceManager) (err error) {
	if params.Response != nil {
		err = DefineParameterMap(ctx, schema, node, params.Response, flow)
		if err != nil {
			return err
		}
	}

	err = DefineProperty(ctx, node, params.Message, flow)
	if err != nil {
		return err
	}

	err = DefineProperty(ctx, node, params.Status, flow)
	if err != nil {
		return err
	}

	err = DefineParams(ctx, node, params.Params, flow)
	if err != nil {
		return err
	}

	return nil
}
