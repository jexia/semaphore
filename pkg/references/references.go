package references

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/trace"
	"github.com/jexia/semaphore/pkg/lookup"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"go.uber.org/zap"
)

// Resolve all references inside the given flow list
func Resolve(ctx *broker.Context, flows specs.FlowListInterface) (err error) {
	logger.Info(ctx, "defining manifest types")

	for _, flow := range flows {
		err := ResolveFlow(logger.WithFields(ctx, zap.String("flow", flow.GetName())), flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// ResolveFlow all references made within the given flow
func ResolveFlow(ctx *broker.Context, flow specs.FlowInterface) (err error) {
	logger.Info(ctx, "defining flow types")

	if flow.GetOnError() != nil {
		err = ResolveOnError(ctx, nil, flow.GetOnError(), flow)
		if err != nil {
			return err
		}
	}

	for _, node := range flow.GetNodes() {
		err = ResolveNode(ctx, node, flow)
		if err != nil {
			return err
		}
	}

	if flow.GetOutput() != nil {
		err = ResolveParameterMap(ctx, nil, flow.GetOutput(), flow)
		if err != nil {
			return err
		}
	}

	if flow.GetForward() != nil && flow.GetForward().Request != nil {
		for _, header := range flow.GetForward().Request.Header {
			err = ResolveProperty(ctx, nil, header, flow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ResolveNode resolves all references made within the given node
func ResolveNode(ctx *broker.Context, node *specs.Node, flow specs.FlowInterface) (err error) {
	if node.Condition != nil {
		err = ResolveParameterMap(ctx, node, node.Condition.Params, flow)
		if err != nil {
			return err
		}
	}

	if node.Call != nil {
		err = ResolveCall(ctx, node, node.Call, flow)
		if err != nil {
			return err
		}
	}

	if node.Rollback != nil {
		err = ResolveCall(ctx, node, node.Rollback, flow)
		if err != nil {
			return err
		}
	}

	if node.OnError != nil {
		err = ResolveOnError(ctx, node, node.OnError, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// ResolveCall resolves all references made within the given call
func ResolveCall(ctx *broker.Context, node *specs.Node, call *specs.Call, flow specs.FlowInterface) (err error) {
	if call.Request != nil {
		err = ResolveParameterMap(ctx, node, call.Request, flow)
		if err != nil {
			return err
		}
	}

	if call.Response != nil {
		err = ResolveParameterMap(ctx, node, call.Response, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// ResolveParameterMap resolves all references made within the given parameter map
func ResolveParameterMap(ctx *broker.Context, node *specs.Node, params *specs.ParameterMap, flow specs.FlowInterface) (err error) {
	for _, header := range params.Header {
		err = ResolveProperty(ctx, node, header, flow)
		if err != nil {
			return err
		}
	}

	if params.Params != nil {
		err = ResolveParams(ctx, node, params.Params, flow)
		if err != nil {
			return err
		}
	}

	if params.Property != nil {
		err = ResolveProperty(ctx, node, params.Property, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// ResolveParams resolves all references made within the given parameters
func ResolveParams(ctx *broker.Context, node *specs.Node, params map[string]*specs.Property, flow specs.FlowInterface) error {
	for _, param := range params {
		if param.Reference == nil {
			continue
		}

		err := ResolveProperty(ctx, node, param, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// ResolveProperty resolves all references made within the given property
func ResolveProperty(ctx *broker.Context, node *specs.Node, property *specs.Property, flow specs.FlowInterface) error {
	if property == nil {
		return nil
	}

	if len(property.Nested) > 0 {
		for _, nested := range property.Nested {
			err := ResolveProperty(ctx, node, nested, flow)
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
		breakpoint = node.ID

		if node.Rollback != nil && property != nil {
			rollback := node.Rollback.Request.Property
			if InsideProperty(rollback, property) {
				breakpoint = lookup.GetNextResource(flow, breakpoint)
			}
		}
	}

	reference, err := LookupReference(ctx, breakpoint, property.Reference, flow)
	if err != nil {
		logger.Error(ctx, "unable to lookup reference", zap.Error(err))
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("undefined reference '%s' in '%s.%s.%s'", property.Reference, flow.GetName(), breakpoint, property.Path))
	}

	if reference.Reference != nil && reference.Reference.Property == nil {
		err := ResolveProperty(ctx, node, reference, flow)
		if err != nil {
			return err
		}
	}

	logger.Debug(ctx, "references lookup result",
		zap.String("reference", property.Reference.String()),
		zap.String("name", property.Name),
		zap.String("path", property.Path),
	)

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
func LookupReference(ctx *broker.Context, breakpoint string, reference *specs.PropertyReference, flow specs.FlowInterface) (*specs.Property, error) {
	reference.Resource = lookup.ResolveSelfReference(reference.Resource, breakpoint)

	logger.Debug(ctx, "lookup references until breakpoint",
		zap.String("reference", reference.String()),
		zap.String("breakpoint", breakpoint),
	)

	references := lookup.GetAvailableResources(flow, breakpoint)
	result := lookup.GetResourceReference(reference, references, breakpoint)
	if result == nil {
		return nil, trace.New(trace.WithMessage("undefined resource '%s' in '%s'.'%s'", reference, flow.GetName(), breakpoint))
	}

	logger.Debug(ctx, "lookup references result",
		zap.String("reference", reference.String()),
		zap.String("path", result.Path),
		zap.String("name", result.Name),
	)

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

// ResolveOnError resolves references made inside the given on error specs
func ResolveOnError(ctx *broker.Context, node *specs.Node, params *specs.OnError, flow specs.FlowInterface) (err error) {
	if params.Response != nil {
		err = ResolveParameterMap(ctx, node, params.Response, flow)
		if err != nil {
			return err
		}
	}

	err = ResolveProperty(ctx, node, params.Message, flow)
	if err != nil {
		return err
	}

	err = ResolveProperty(ctx, node, params.Status, flow)
	if err != nil {
		return err
	}

	err = ResolveParams(ctx, node, params.Params, flow)
	if err != nil {
		return err
	}

	return nil
}
