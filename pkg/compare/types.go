package compare

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
	"go.uber.org/zap"
)

// Types compares the types defined insde the schema definitions against the configured specification
func Types(ctx *broker.Context, services specs.ServiceList, objects specs.Schemas, flows specs.FlowListInterface) (err error) {
	logger.Info(ctx, "Comparing manifest types")

	for _, flow := range flows {
		err := FlowTypes(ctx, services, objects, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// ProxyTypes compares the given proxy against the configured schema types
func ProxyTypes(ctx *broker.Context, services specs.ServiceList, objects specs.Schemas, proxy *specs.Proxy) (err error) {
	logger.Info(ctx, "Compare proxy flow types", zap.String("proxy", proxy.GetName()))

	if proxy.OnError != nil {
		err = CheckParameterMapTypes(ctx, proxy.OnError.Response, objects, proxy)
		if err != nil {
			return err
		}
	}

	for _, node := range proxy.Nodes {
		err = CallTypes(ctx, services, objects, node, node.Call, proxy)
		if err != nil {
			return err
		}

		err = CallTypes(ctx, services, objects, node, node.Rollback, proxy)
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
func FlowTypes(ctx *broker.Context, services specs.ServiceList, objects specs.Schemas, flow specs.FlowInterface) (err error) {
	logger.Info(ctx, "Comparing flow types", zap.String("flow", flow.GetName()))

	if flow.GetInput() != nil {
		err = CheckParameterMapTypes(ctx, flow.GetInput(), objects, flow)
		if err != nil {
			return err
		}
	}

	if flow.GetOnError() != nil {
		err = CheckParameterMapTypes(ctx, flow.GetOnError().Response, objects, flow)
		if err != nil {
			return err
		}
	}

	for _, node := range flow.GetNodes() {
		err = CallTypes(ctx, services, objects, node, node.Call, flow)
		if err != nil {
			return err
		}

		err = CallTypes(ctx, services, objects, node, node.Rollback, flow)
		if err != nil {
			return err
		}
	}

	if flow.GetOutput() != nil {
		message := objects.Get(flow.GetOutput().Schema)
		if message == nil {
			return ErrUndefinedObject{
				Flow:   flow.GetName(),
				Schema: flow.GetOutput().Schema,
			}
		}

		err = CheckParameterMapTypes(ctx, flow.GetOutput(), objects, flow)
		if err != nil {
			return err
		}
	}

	if flow.GetForward() != nil {
		if flow.GetForward().Request != nil && flow.GetForward().Request.Header != nil {
			err = CheckHeader(flow.GetForward().Request.Header, flow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CallTypes compares the given call types against the configured schema types
func CallTypes(ctx *broker.Context, services specs.ServiceList, objects specs.Schemas, node *specs.Node, call *specs.Call, flow specs.FlowInterface) (err error) {
	if call == nil {
		return nil
	}

	if call.Method == "" {
		return nil
	}

	logger.Info(ctx, "Comparing call types", zap.String("call", node.ID), zap.String("method", call.Method), zap.String("service", call.Service))

	service := services.Get(call.Service)
	if service == nil {
		return ErrUndefinedService{
			Flow:    flow.GetName(),
			Service: call.Service,
		}
	}

	method := service.GetMethod(call.Method)
	if method == nil {
		return ErrUndefinedMethod{
			Flow:   flow.GetName(),
			Method: call.Method,
		}
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
func CheckParameterMapTypes(ctx *broker.Context, parameters *specs.ParameterMap, objects specs.Schemas, flow specs.FlowInterface) error {
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
		return ErrUndefinedSchema{
			Path: property.Path,
			Expr: property.Expr,
		}
	}
	if property.Type != schema.Type {
		return ErrTypeMismatch{
			Expr:     property.Expr,
			Type:     property.Type,
			Expected: schema.Type,
			Path:     property.Path,
		}
	}

	if property.Label != schema.Label {
		return ErrLabelMismatch{
			Expr:     property.Expr,
			Label:    property.Label,
			Expected: schema.Label,
			Path:     property.Path,
		}
	}

	if len(property.Nested) > 0 {
		if len(schema.Nested) == 0 {
			return ErrUndeclaredSchema{
				Expr: property.Expr,
				Path: property.Path,
				Name: property.Name,
			}
		}

		for key, nested := range property.Nested {
			object := schema.Nested[key]
			if object == nil {
				return ErrUndeclaredSchemaInProperty{
					Expr: nested.Expr,
					Path: nested.Path,
					Name: flow.GetName(),
				}
			}

			err := CheckPropertyTypes(nested, object, flow)
			if err != nil {
				return err
			}
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
			return ErrHeaderTypeMismatch{
				Type:     header.Type,
				Path:     header.Path,
				Flow:     flow.GetName(),
				Expected: types.String,
			}
		}
	}

	return nil
}
