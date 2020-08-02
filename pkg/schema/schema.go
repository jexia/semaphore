package schema

import (
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/jexia/semaphore/pkg/core/trace"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/sirupsen/logrus"
)

// Define defines the types of all schemas inside the given flow list
func Define(ctx instance.Context, services specs.ServiceList, schemas specs.Objects, flows specs.FlowListInterface) (err error) {
	ctx.Logger(logger.Core).Info("Defining manifest types")

	for _, flow := range flows {
		err := DefineFlow(ctx, services, schemas, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineFlow defines the types for the given flow and the resources within the flow
func DefineFlow(ctx instance.Context, services specs.ServiceList, schemas specs.Objects, flow specs.FlowInterface) (err error) {
	ctx.Logger(logger.Core).WithField("flow", flow.GetName()).Info("Defining flow types")

	if flow.GetInput() != nil {
		input := schemas.Get(flow.GetInput().Schema)
		if input == nil {
			return trace.New(trace.WithMessage("object '%s', is unavailable inside the schema collection", flow.GetInput().Schema))
		}

		flow.GetInput().Property = input
	}

	if flow.GetOnError() != nil {
		err = DefineOnError(ctx, schemas, flow.GetOnError(), flow)
		if err != nil {
			return err
		}
	}

	for _, node := range flow.GetNodes() {
		err = DefineNode(ctx, services, schemas, node, flow)
		if err != nil {
			return err
		}
	}

	if flow.GetOutput() != nil {
		err = DefineParameterMap(ctx, schemas, flow.GetOutput(), flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineNode defines all the references inside the given node
func DefineNode(ctx instance.Context, services specs.ServiceList, schemas specs.Objects, node *specs.Node, flow specs.FlowInterface) (err error) {
	if node.Condition != nil {
		err = DefineParameterMap(ctx, schemas, node.Condition.Params, flow)
		if err != nil {
			return err
		}
	}

	if node.Call != nil {
		err = DefineCall(ctx, services, schemas, node, node.Call, flow)
		if err != nil {
			return err
		}
	}

	if node.Rollback != nil {
		err = DefineCall(ctx, services, schemas, node, node.Rollback, flow)
		if err != nil {
			return err
		}
	}

	if node.OnError != nil {
		err = DefineOnError(ctx, schemas, node.OnError, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineCall defineds the types for the specs call
func DefineCall(ctx instance.Context, services specs.ServiceList, schemas specs.Objects, node *specs.Node, call *specs.Call, flow specs.FlowInterface) (err error) {
	if call.Request != nil {
		err = DefineParameterMap(ctx, schemas, call.Request, flow)
		if err != nil {
			return err
		}
	}

	if call.Method != "" {
		ctx.Logger(logger.Core).WithFields(logrus.Fields{
			"call":    node.ID,
			"method":  call.Method,
			"service": call.Service,
		}).Info("Defining call types")

		service := services.Get(call.Service)
		if service == nil {
			return trace.New(trace.WithMessage("undefined service '%s' in flow '%s'", call.Service, flow.GetName()))
		}

		method := service.GetMethod(call.Method)
		if method == nil {
			return trace.New(trace.WithMessage("undefined method '%s' in flow '%s'", call.Method, flow.GetName()))
		}

		output := schemas.Get(method.Output)
		if output == nil {
			return trace.New(trace.WithMessage("undefined method output property '%s' in flow '%s'", method.Output, flow.GetName()))
		}

		call.Descriptor = method
		call.Response = &specs.ParameterMap{
			Property: output,
		}

		call.Request.Schema = method.Input
		call.Response.Schema = method.Output
	}

	if call.Response != nil {
		err = DefineParameterMap(ctx, schemas, call.Response, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineProperty ensures that all schema properties are defined
func DefineProperty(property *specs.Property, schema *specs.Property, flow specs.FlowInterface) error {
	if property == nil {
		property = schema
		return nil
	}

	for key, nested := range property.Nested {
		object := schema.Nested[key]
		if object == nil {
			return trace.New(trace.WithMessage("undefined schema nested message property '%s' in flow '%s'", key, flow.GetName()))
		}

		err := DefineProperty(nested, object, flow)
		if err != nil {
			return err
		}
	}

	if property.Nested == nil && schema.Nested != nil {
		property.Nested = schema.Nested
		return nil
	}

	// Set any properties not defined inside the flow but available inside the schema
	for _, prop := range schema.Nested {
		_, has := property.Nested[prop.Name]
		if has {
			continue
		}

		property.Nested[prop.Name] = prop
	}

	return nil
}

// DefineParameterMap defines the types for the given parameter map
func DefineParameterMap(ctx instance.Context, schemas specs.Objects, params *specs.ParameterMap, flow specs.FlowInterface) (err error) {
	if params == nil || params.Schema == "" {
		return nil
	}

	schema := schemas.Get(params.Schema)
	if schema == nil {
		return trace.New(trace.WithMessage("object '%s', is unavailable inside the schema collection", params.Schema))
	}

	err = DefineProperty(params.Property, schema, flow)
	if err != nil {
		return err
	}

	return nil
}

// DefineOnError defines references made inside the given on error specs
func DefineOnError(ctx instance.Context, schemas specs.Objects, params *specs.OnError, flow specs.FlowInterface) (err error) {
	if params.Response != nil {
		err = DefineParameterMap(ctx, schemas, params.Response, flow)
		if err != nil {
			return err
		}
	}

	return nil
}
