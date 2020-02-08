package types

import (
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jhump/protoreflect/desc"
)

// Define checks and defines the types for the given manifest
func Define(desc []*desc.FileDescriptor, manifest *specs.Manifest) (err error) {
	for _, flow := range manifest.Flows {
		for _, call := range flow.Calls {
			err = DefineCall(desc, manifest, call, flow)
			if err != nil {
				return err
			}

			if call.Rollback != nil {
				err = DefineCall(desc, manifest, call.Rollback, flow)
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
func DefineCall(desc []*desc.FileDescriptor, manifest *specs.Manifest, call specs.FlowCaller, flow *specs.Flow) (err error) {
	if call.GetEndpoint() == "" {
		return nil
	}

	service := FindService(desc, GetServiceProto(manifest, GetService(call.GetEndpoint())))
	if service == nil {
		return trace.New(trace.WithMessage("undefined service alias '%s' in flow '%s'", GetServiceProto(manifest, GetService(call.GetEndpoint())), flow.Name))
	}

	method := FindMethod(service, GetMethod(call.GetEndpoint()))
	if method == nil {
		return trace.New(trace.WithMessage("undefined method '%s' in flow '%s'", GetMethod(call.GetEndpoint()), flow.Name))
	}

	call.SetDescriptor(method)

	if call.GetRequest() != nil {
		err = DefineParameterMap(call, call.GetRequest(), flow)
		if err != nil {
			return err
		}

		err = CheckProtoTypes(call.GetRequest(), method.GetInputType(), flow)
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

	references := specs.GetAvailableResources(flow, call.GetName())
	reference := specs.GetResourceReference(property.Reference, references)
	if reference == nil {
		return trace.New(trace.WithExpression(property.Expr), trace.WithExpression(property.Expr), trace.WithExpression(property.Expr), trace.WithMessage("undefined resource '%s' in flow '%s'", property.Reference, flow.Name))
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

// CheckProtoTypes checks the given call against the given proto method types
func CheckProtoTypes(object specs.Object, message *desc.MessageDescriptor, flow *specs.Flow) (err error) {
	for key, property := range object.GetProperties() {
		field := GetField(message, key)
		if field == nil {
			return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("undefined proto field '%s' in flow '%s'", property.Path, flow.Name))
		}

		if property.Type != specs.GetType(field.GetType()) {
			return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("cannot use (%s) type (%s) in '%s'", property.Type, specs.GetType(field.GetType()), field.GetFullyQualifiedName()))
		}
	}

	for key, nested := range object.GetNestedProperties() {
		field := GetField(message, key)
		if field == nil {
			return trace.New(trace.WithMessage("undefined proto nested message '%s' in flow '%s'", nested.Path, flow.Name))
		}

		if specs.GetType(field.GetType()) != specs.TypeMessage {
			return trace.New(trace.WithMessage("cannot use (%s) type in proto (%s)", specs.TypeMessage, specs.GetType(field.GetType())))
		}

		err = CheckProtoTypes(nested.GetObject(), field.GetMessageType(), flow)
		if err != nil {
			return err
		}
	}

	for key, repeated := range object.GetRepeatedProperties() {
		field := GetField(message, key)
		if field == nil {
			return trace.New(trace.WithMessage("undefined proto repeated message '%s' in flow '%s'", repeated.Path, flow.Name))
		}

		if specs.GetType(field.GetType()) != specs.TypeMessage {
			return trace.New(trace.WithMessage("cannot use (%s) type in proto (%s)", specs.TypeMessage, specs.GetType(field.GetType())))
		}

		err = CheckProtoTypes(repeated.GetObject(), field.GetMessageType(), flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetField attempts to get the given field inside the given method
func GetField(message *desc.MessageDescriptor, name string) *desc.FieldDescriptor {
	for _, field := range message.GetFields() {
		if field.GetName() == name {
			return field
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

		if property.Reference.Object.GetLabel() == specs.LabelRepeated {
			repeated := property.Reference.Object.(*specs.RepeatedParameterMap)
			params.GetRepeatedProperties()[key] = repeated.Clone(key, property.Path)
			continue
		}

		nested := property.Reference.Object.(*specs.NestedParameterMap)
		params.GetNestedProperties()[key] = nested.Clone(key, property.Path)
	}
}
