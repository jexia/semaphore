package strict

import (
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/lookup"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/specs/types"
)

// Define checks and defines the types for the given manifest
func Define(schema schema.Collection, manifest *specs.Manifest) (err error) {
	for _, flow := range manifest.Flows {
		for _, call := range flow.Calls {
			err = DefineCall(schema, manifest, call, flow)
			if err != nil {
				return err
			}

			if call.Rollback != nil {
				err = DefineCall(schema, manifest, call.Rollback, flow)
				if err != nil {
					return err
				}
			}
		}

		if flow.Output != nil {
			final := flow.Calls[len(flow.Calls)-1]
			err = DefineParameterMap(final, flow.Output, flow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DefineCall defineds the types for the given parameter map
func DefineCall(schema schema.Collection, manifest *specs.Manifest, call specs.FlowCaller, flow *specs.Flow) (err error) {
	if call.GetEndpoint() == "" {
		return nil
	}

	service := schema.GetService(GetSchemaService(manifest, GetService(call.GetEndpoint())))
	if service == nil {
		return trace.New(trace.WithMessage("undefined service alias '%s' in flow '%s'", GetService(call.GetEndpoint()), flow.Name))
	}

	method := service.GetMethod(GetMethod(call.GetEndpoint()))
	if method == nil {
		return trace.New(trace.WithMessage("undefined method '%s' in flow '%s'", GetMethod(call.GetEndpoint()), flow.Name))
	}

	call.SetDescriptor(method)

	if call.GetRequest() != nil {
		err = DefineParameterMap(call, call.GetRequest(), flow)
		if err != nil {
			return err
		}

		err = CheckTypes(call.GetRequest(), method.GetInput(), flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineParameterMap defines the types for the given parameter map
func DefineParameterMap(call specs.FlowCaller, params specs.Object, flow *specs.Flow) (err error) {
	for _, property := range params.GetProperties() {
		err = DefineProperty(call, property, flow)
		if err != nil {
			return err
		}
	}

	for _, nested := range params.GetNestedProperties() {
		err = DefineParameterMap(call, nested, flow)
		if err != nil {
			return err
		}
	}

	for _, repeated := range params.GetRepeatedProperties() {
		err = DefineParameterMap(call, repeated, flow)
		if err != nil {
			return err
		}
	}

	ResolvePropertyObjectReferences(params)

	return nil
}

// DefineProperty defines the given property type.
// If any object is references it has to be fixed afterwards and moved into the correct dataset
func DefineProperty(call specs.FlowCaller, property *specs.Property, flow *specs.Flow) error {
	if property.Reference == nil {
		return nil
	}

	references := lookup.GetAvailableResources(flow, call.GetName())
	reference := lookup.GetResourceReference(property.Reference, references)
	if reference == nil {
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("undefined resource '%s' in '%s.%s.%s'", property.Reference, flow.Name, call.GetName(), property.Path))
	}

	property.Type = reference.GetType()
	property.Default = reference.GetDefault()
	property.Reference.Object = reference.GetObject()

	if reference.GetObject() != nil {
		property.Reference.Label = reference.GetObject().GetLabel()
	}

	// TODO: support enum type

	return nil
}

// CheckTypes checks the given call against the given schema method types
func CheckTypes(object specs.Object, message schema.Object, flow *specs.Flow) (err error) {
	for key, property := range object.GetProperties() {
		field := message.GetField(key)
		if field == nil {
			return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("undefined schema field '%s' in flow '%s'", property.Path, flow.Name))
		}

		if property.Type != field.GetType() {
			return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("cannot use (%s) type (%s) in '%s'", property.Type, field.GetType(), field.GetName()))
		}
	}

	for key, nested := range object.GetNestedProperties() {
		field := message.GetField(key)
		if field == nil {
			return trace.New(trace.WithMessage("undefined schema nested message '%s' in flow '%s'", nested.Path, flow.Name))
		}

		if field.GetType() != types.TypeMessage {
			return trace.New(trace.WithMessage("cannot use (%s) type in schema (%s)", types.TypeMessage, field.GetType()))
		}

		err = CheckTypes(nested.GetObject(), field.GetObject(), flow)
		if err != nil {
			return err
		}
	}

	for key, repeated := range object.GetRepeatedProperties() {
		field := message.GetField(key)
		if field == nil {
			return trace.New(trace.WithMessage("undefined schema repeated message '%s' in flow '%s'", repeated.Path, flow.Name))
		}

		if field.GetType() != types.TypeMessage {
			return trace.New(trace.WithMessage("cannot use (%s) type in schema (%s)", types.TypeMessage, field.GetType()))
		}

		err = CheckTypes(repeated.GetObject(), field.GetObject(), flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// ResolvePropertyObjectReferences moves any property object reference into the correct data structure
func ResolvePropertyObjectReferences(params specs.Object) {
	for key, property := range params.GetProperties() {
		if property.Reference == nil {
			continue
		}

		if property.Reference.Object == nil {
			continue
		}

		delete(params.GetProperties(), key)

		if property.Reference.Object.GetLabel() == types.LabelRepeated {
			repeated := property.Reference.Object.(*specs.RepeatedParameterMap)
			params.GetRepeatedProperties()[key] = repeated.Clone(key, property.Path)
			continue
		}

		nested := property.Reference.Object.(*specs.NestedParameterMap)
		params.GetNestedProperties()[key] = nested.Clone(key, property.Path)
	}
}

// GetSchemaService attempts to find a service matching the alias name and return the schema name
func GetSchemaService(manifest *specs.Manifest, alias string) string {
	for _, service := range manifest.Services {
		if service.Alias == alias {
			return service.Schema
		}
	}

	return ""
}
