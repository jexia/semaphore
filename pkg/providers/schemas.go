package providers

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
	"go.uber.org/zap"
)

// ResolveSchemaDefinitions ensures that all schema definitions are valid and are
// available inside the schema collection. Call method input and output schemas
// are defined and set inside the given calls.
func ResolveSchemaDefinitions(ctx *broker.Context, services specs.ServiceList, schemas specs.Schemas, flows specs.FlowListInterface) (err error) {
	logger.Info(ctx, "defining manifest types")

	for _, flow := range flows {
		err := resolveFlowSchemaDefinitions(ctx, services, schemas, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

func resolveFlowSchemaDefinitions(parent *broker.Context, services specs.ServiceList, schemas specs.Schemas, flow specs.FlowInterface) (err error) {
	ctx := logger.WithFields(parent, zap.String("flow", flow.GetName()))
	logger.Info(ctx, "defining flow types")

	if flow.GetInput() != nil {
		input := schemas.Get(flow.GetInput().Schema)
		if input == nil {
			return ErrUndefinedObject{
				Schema: flow.GetInput().Schema,
			}
		}

		flow.GetInput().Property = input.Clone()
	}

	if flow.GetOnError() != nil {
		err = resolveOnErrorSchemaDefinitions(ctx, schemas, flow.GetOnError(), flow)
		if err != nil {
			return err
		}
	}

	for _, node := range flow.GetNodes() {
		err = resolveNodeSchemaDefinitions(ctx, services, schemas, node, flow)
		if err != nil {
			return err
		}
	}

	if flow.GetOutput() != nil {
		err = resolveParameterMapSchemaDefinitions(ctx, schemas, flow.GetOutput(), flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// resolveNodeSchemaDefinitions ensures that all schema properties are defined inside the given node
func resolveNodeSchemaDefinitions(ctx *broker.Context, services specs.ServiceList, schemas specs.Schemas, node *specs.Node, flow specs.FlowInterface) (err error) {
	if node.Condition != nil {
		err = resolveParameterMapSchemaDefinitions(ctx, schemas, node.Condition.Params, flow)
		if err != nil {
			return err
		}
	}

	if node.Call != nil {
		err = defineCallSchemaDefinitions(ctx, services, schemas, node, node.Call, flow)
		if err != nil {
			return err
		}
	}

	if node.Rollback != nil {
		err = defineCallSchemaDefinitions(ctx, services, schemas, node, node.Rollback, flow)
		if err != nil {
			return err
		}
	}

	if node.Intermediate != nil {
		err = resolveParameterMapSchemaDefinitions(ctx, schemas, node.Intermediate, flow)
		if err != nil {
			return err
		}
	}

	if node.OnError != nil {
		err = resolveOnErrorSchemaDefinitions(ctx, schemas, node.OnError, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// defineCallSchemaDefinitions defineds the types for the specs call
func defineCallSchemaDefinitions(ctx *broker.Context, services specs.ServiceList, schemas specs.Schemas, node *specs.Node, call *specs.Call, flow specs.FlowInterface) (err error) {
	if call.Request != nil {
		err = resolveParameterMapSchemaDefinitions(ctx, schemas, call.Request, flow)
		if err != nil {
			return err
		}
	}

	if call.Method != "" {
		logger.Info(ctx, "defining call types",
			zap.String("call", node.ID),
			zap.String("method", call.Method),
			zap.String("service", call.Service),
		)

		service := services.Get(call.Service)
		if service == nil {
			return ErrUndefinedService{
				Service: call.Service,
				Flow:    flow.GetName(),
			}
		}

		method := service.GetMethod(call.Method)
		if method == nil {
			return ErrUndefinedMethod{
				Flow:   flow.GetName(),
				Method: call.Method,
			}
		}

		output := schemas.Get(method.Output)
		if output == nil {
			return ErrUndefinedOutput{
				Output: method.Output,
				Flow:   flow.GetName(),
			}
		}

		call.Descriptor = method
		call.Response = &specs.ParameterMap{
			Property: output.Clone(),
		}

		call.Request.Schema = method.Input
		call.Response.Schema = method.Output
	}

	if call.Response != nil {
		err = resolveParameterMapSchemaDefinitions(ctx, schemas, call.Response, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// resolveParameterMapSchemaDefinitions ensures that the given parameter map schema is available
func resolveParameterMapSchemaDefinitions(ctx *broker.Context, schemas specs.Schemas, params *specs.ParameterMap, flow specs.FlowInterface) error {
	if params == nil || params.Schema == "" {
		return nil
	}

	schema := schemas.Get(params.Schema)
	if schema == nil {
		return ErrUndefinedObject{
			Schema: params.Schema,
		}
	}

	return nil
}

// resolveOnErrorSchemaDefinitions ensures that all schema properties are defined inside the given on error object
func resolveOnErrorSchemaDefinitions(ctx *broker.Context, schemas specs.Schemas, params *specs.OnError, flow specs.FlowInterface) (err error) {
	if params.Response != nil {
		err = resolveParameterMapSchemaDefinitions(ctx, schemas, params.Response, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineSchemaPropertyPositions defines property labels and positions
func DefineSchemaPropertyPositions(property, schema *specs.Property, flow specs.FlowInterface) error {
	property.Label = schema.Label
	property.Position = schema.Position

	return nil
}
