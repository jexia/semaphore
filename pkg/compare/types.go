package compare

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/trace"
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
			return trace.New(trace.WithMessage("undefined flow output object '%s' in '%s'", flow.GetOutput().Schema, flow.GetName()))
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

	err := parameters.Property.Compare(objects.Get(parameters.Schema))
	if err != nil {
		return fmt.Errorf("flow '%s' mismatch: %w", flow.GetName(), err)
	}

	return nil
}

// CheckHeader compares the given header types
func CheckHeader(header specs.Header, flow specs.FlowInterface) error {
	for _, header := range header {
		if header.Type() != types.String {
			return trace.New(trace.WithMessage("cannot use type (%s) for 'header.%s' in flow '%s', expected (%s)", header.Type(), header.Path, flow.GetName(), types.String))
		}
	}

	return nil
}
