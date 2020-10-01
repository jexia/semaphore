package references

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
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
			return ErrUnresolvedFlow{
				wrapErr: wrapErr{err},
				Name:    flow.GetName(),
			}
		}
	}

	return nil
}

// ResolveFlow all references made within the given flow
func ResolveFlow(ctx *broker.Context, flow specs.FlowInterface) (err error) {
	if flow.GetOnError() != nil {
		err = ResolveOnError(ctx, nil, flow.GetOnError(), flow)
		if err != nil {
			return ErrUnresolvedOnError{
				wrapErr: wrapErr{err},
				OnError: flow.GetOnError(),
			}
		}
	}

	for _, node := range flow.GetNodes() {
		err = ResolveNode(ctx, node, flow)
		if err != nil {
			return ErrUnresolvedNode{
				wrapErr: wrapErr{err},
				Node:    node,
			}
		}
	}

	if flow.GetOutput() != nil {
		err = ResolveParameterMap(ctx, nil, flow.GetOutput(), flow)
		if err != nil {
			return ErrUnresolvedParameterMap{
				wrapErr:   wrapErr{err},
				Parameter: flow.GetOutput(),
			}
		}
	}

	if flow.GetForward() != nil && flow.GetForward().Request != nil {
		for _, header := range flow.GetForward().Request.Header {
			err = ResolveProperty(ctx, nil, header, flow)
			if err != nil {
				return ErrUnresolvedProperty{
					wrapErr:  wrapErr{err},
					Property: header,
				}
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
			return ErrUnresolvedParameterMap{
				wrapErr:   wrapErr{err},
				Parameter: node.Condition.Params,
			}
		}
	}

	if node.Intermediate != nil {
		err = ResolveParameterMap(ctx, node, node.Intermediate, flow)
		if err != nil {
			return ErrUnresolvedParameterMap{
				wrapErr:   wrapErr{err},
				Parameter: node.Intermediate,
			}
		}
	}

	if node.Call != nil {
		err = ResolveCall(ctx, node, node.Call, flow)
		if err != nil {
			return ErrUnresolvedCall{
				wrapErr: wrapErr{err},
				Call:    node.Call,
			}
		}
	}

	if node.Rollback != nil {
		err = ResolveCall(ctx, node, node.Rollback, flow)
		if err != nil {
			return ErrUnresolvedCall{
				wrapErr: wrapErr{err},
				Call:    node.Rollback,
			}
		}
	}

	if node.OnError != nil {
		err = ResolveOnError(ctx, node, node.OnError, flow)
		if err != nil {
			return ErrUnresolvedOnError{
				wrapErr: wrapErr{err},
				OnError: node.OnError,
			}
		}
	}

	return nil
}

// ResolveCall resolves all references made within the given call
func ResolveCall(ctx *broker.Context, node *specs.Node, call *specs.Call, flow specs.FlowInterface) (err error) {
	if call.Request != nil {
		err = ResolveParameterMap(ctx, node, call.Request, flow)
		if err != nil {
			return ErrUnresolvedParameterMap{
				wrapErr:   wrapErr{err},
				Parameter: call.Request,
			}
		}
	}

	if call.Response != nil {
		err = ResolveParameterMap(ctx, node, call.Response, flow)
		if err != nil {
			return ErrUnresolvedParameterMap{
				wrapErr:   wrapErr{err},
				Parameter: call.Response,
			}
		}
	}

	return nil
}

// ResolveParameterMap resolves all references made within the given parameter map
func ResolveParameterMap(ctx *broker.Context, node *specs.Node, params *specs.ParameterMap, flow specs.FlowInterface) (err error) {
	for _, header := range params.Header {
		err = ResolveProperty(ctx, node, header, flow)
		if err != nil {
			return ErrUnresolvedProperty{
				wrapErr:  wrapErr{err},
				Property: header,
			}
		}
	}

	if params.Params != nil {
		err = ResolveParams(ctx, node, params.Params, flow)
		if err != nil {
			return ErrUnresolvedParams{
				wrapErr: wrapErr{err},
				Params:  params.Params,
			}
		}
	}

	if params.Property != nil {
		err = ResolveProperty(ctx, node, params.Property, flow)
		if err != nil {
			return ErrUnresolvedProperty{
				wrapErr:  wrapErr{err},
				Property: params.Property,
			}
		}
	}

	return nil
}

// ResolveOnError resolves references made inside the given on error specs
func ResolveOnError(ctx *broker.Context, node *specs.Node, params *specs.OnError, flow specs.FlowInterface) (err error) {
	if params.Response != nil {
		err = ResolveParameterMap(ctx, node, params.Response, flow)
		if err != nil {
			return ErrUnresolvedParameterMap{
				wrapErr:   wrapErr{err},
				Parameter: params.Response,
			}
		}
	}

	err = ResolveProperty(ctx, node, params.Message, flow)
	if err != nil {
		return ErrUnresolvedProperty{
			wrapErr:  wrapErr{err},
			Property: params.Message,
		}
	}

	err = ResolveProperty(ctx, node, params.Status, flow)
	if err != nil {
		return ErrUnresolvedProperty{
			wrapErr:  wrapErr{err},
			Property: params.Status,
		}
	}

	err = ResolveParams(ctx, node, params.Params, flow)
	if err != nil {
		return ErrUnresolvedParams{
			wrapErr: wrapErr{err},
			Params:  params.Params,
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
			return ErrUnresolvedProperty{
				wrapErr:  wrapErr{err},
				Property: param,
			}
		}
	}

	return nil
}

func resolveProperty(ctx *broker.Context, node *specs.Node, property *specs.Property, flow specs.FlowInterface) error {
	switch {
	case property.Message != nil:
		for _, nested := range property.Message {
			err := ResolveProperty(ctx, node, nested, flow)
			if err != nil {
				return NewErrUnresolvedProperty(err, nested)
			}
		}

		return nil
	case property.Repeated != nil:
		for _, repeated := range property.Repeated {
			property := &specs.Property{
				Template: repeated,
			}

			err := ResolveProperty(ctx, node, property, flow)
			if err != nil {
				return NewErrUnresolvedProperty(err, property)
			}
		}

		return nil
	default:
		return nil
	}
}

// ResolveProperty resolves all references made within the given property
func ResolveProperty(ctx *broker.Context, node *specs.Node, property *specs.Property, flow specs.FlowInterface) error {
	if property == nil {
		return nil
	}

	if err := resolveProperty(ctx, node, property, flow); err != nil {
		return NewErrUnresolvedProperty(err, property)
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
		return NewErrUndefinedReference(err, property, breakpoint)
	}

	if reference.Reference != nil && reference.Reference.Property == nil {
		err := ResolveProperty(ctx, node, reference, flow)
		if err != nil {
			return NewErrUnresolvedProperty(err, reference)
		}
	}

	logger.Debug(ctx, "references lookup result",
		zap.String("reference", property.Reference.String()),
		zap.String("name", property.Name),
		zap.String("path", property.Path),
	)

	property.Label = reference.Label
	property.Reference.Property = reference

	ScopeNestedReferences(&reference.Template, &property.Template)

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
		return nil, ErrUndefinedResource{
			Reference:           reference,
			Breakpoint:          breakpoint,
			AvailableReferences: references,
		}
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

	switch {
	case source.Message != nil:
		for _, nested := range source.Message {
			if InsideProperty(nested, target) {
				return true
			}
		}
	case source.Repeated != nil:
		for _, repeated := range source.Repeated {
			property := &specs.Property{
				Template: repeated,
			}

			if InsideProperty(property, target) {
				return true
			}
		}
	}

	return false
}

// ScopeNestedReferences scopes all nested references available inside the reference property
func ScopeNestedReferences(source, target *specs.Template) {
	if source == nil || target == nil {
		return
	}

	switch {
	case source.Scalar != nil:
		if target.Scalar == nil {
			target.Scalar = &specs.Scalar{}
		}

		target.Scalar.Default = source.Scalar.Default
		target.Scalar.Type = source.Scalar.Type
		break
	case source.Enum != nil:
		target.Enum = source.Enum
		break
	case source.Message != nil:
		if target.Message == nil {
			target.Message = make(specs.Message, len(source.Message))
		}

		for _, item := range source.Message {
			nested, ok := target.Message[item.Name]
			if !ok {
				nested = item.Clone()
				target.Message[item.Name] = nested
			}

			if nested.Reference == nil {
				nested.Reference = &specs.PropertyReference{
					Resource: target.Reference.Resource,
					Path:     template.JoinPath(target.Reference.Path, item.Name),
					Property: item,
				}
			}

			if len(item.Message) > 0 {
				ScopeNestedReferences(&item.Template, &nested.Template)
			}
		}

		break
	case source.Repeated != nil:
		if target.Repeated != nil {
			return
		}

		target.Repeated = make(specs.Repeated, len(source.Repeated))

		for index, item := range source.Repeated {
			cloned := item.Clone()

			if cloned.Reference == nil {
				cloned.Reference = &specs.PropertyReference{
					Resource: target.Reference.Resource,
					Path:     target.Reference.Path,
				}
			}

			if len(item.Message) > 0 {
				ScopeNestedReferences(&item, &cloned)
			}

			target.Repeated[index] = cloned
		}

		break
	}
}
