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
		err := DefineFlow(schema, manifest, flow)
		if err != nil {
			return err
		}
	}

	for _, proxy := range manifest.Proxy {
		err := DefineProxy(schema, manifest, proxy)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineProxy checks and defines the types for the given proxy
func DefineProxy(schema schema.Collection, manifest *specs.Manifest, proxy *specs.Proxy) (err error) {
	for _, call := range proxy.Calls {
		err = DefineCall(schema, manifest, call, proxy)
		if err != nil {
			return err
		}

		if call.Rollback != nil {
			err = DefineCall(schema, manifest, call.Rollback, proxy)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DefineFlow checks and defines the types for the given flow
func DefineFlow(schema schema.Collection, manifest *specs.Manifest, flow *specs.Flow) (err error) {
	if flow.Schema != "" {
		err = DefineFlowSchema(schema, flow)
		if err != nil {
			return err
		}
	}

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
		err = DefineParameterMap(nil, flow.Output, flow)
		if err != nil {
			return err
		}

		if flow.Descriptor != nil {
			err = CheckTypes(flow.Output, flow.Descriptor.GetOutput(), flow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DefineFlowSchema attempts to define the flow input and output based on the given schema method
func DefineFlowSchema(schema schema.Collection, flow *specs.Flow) error {
	service := schema.GetService(GetService(flow.Schema))
	if service == nil {
		return trace.New(trace.WithMessage("undefined service alias '%s' in flow schema '%s'", GetService(flow.Schema), flow.Name))
	}

	method := service.GetMethod(GetMethod(flow.Schema))
	if method == nil {
		return trace.New(trace.WithMessage("undefined method '%s' in flow schema '%s'", GetMethod(flow.Schema), flow.Name))
	}

	flow.Descriptor = method
	flow.Input = specs.ToParameterMap(flow.Input, "", method.GetInput())

	return nil
}

// DefineCall defineds the types for the given parameter map
func DefineCall(schema schema.Collection, manifest *specs.Manifest, call specs.FlowCaller, flow specs.FlowManager) (err error) {
	if call.GetEndpoint() == "" {
		return nil
	}

	service := schema.GetService(GetSchemaService(manifest, GetService(call.GetEndpoint())))
	if service == nil {
		return trace.New(trace.WithMessage("undefined service alias '%s' in flow '%s'", GetService(call.GetEndpoint()), flow.GetName()))
	}

	method := service.GetMethod(GetMethod(call.GetEndpoint()))
	if method == nil {
		return trace.New(trace.WithMessage("undefined method '%s' in flow '%s'", GetMethod(call.GetEndpoint()), flow.GetName()))
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
func DefineParameterMap(call specs.FlowCaller, params specs.Object, flow specs.FlowManager) (err error) {
	for _, header := range params.GetHeader() {
		err = DefineProperty(call, header, flow)
		if err != nil {
			return err
		}
	}

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
func DefineProperty(call specs.FlowCaller, property *specs.Property, flow specs.FlowManager) error {
	if property.Reference == nil {
		return nil
	}

	breakpoint := "output"
	if call != nil {
		breakpoint = call.GetName()
	}

	references := lookup.GetAvailableResources(flow, breakpoint)
	reference := lookup.GetResourceReference(property.Reference, references)
	if reference == nil {
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("undefined resource '%s' in '%s.%s.%s'", property.Reference, flow.GetName(), breakpoint, property.Path))
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
func CheckTypes(object specs.Object, message schema.Object, flow specs.FlowManager) (err error) {
	for _, header := range object.GetHeader() {
		if header.GetType() != types.TypeString {
			return trace.New(trace.WithMessage("cannot use type %s for header.%s in flow %s", header.GetType(), header.GetPath(), flow.GetName()))
		}
	}

	for key, property := range object.GetProperties() {
		field := message.GetField(key)
		if field == nil {
			return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("undefined schema field '%s' in flow '%s'", property.Path, flow.GetName()))
		}

		if property.Type != field.GetType() {
			return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("cannot use (%s) type (%s) in '%s'", property.Type, field.GetType(), field.GetName()))
		}
	}

	for key, nested := range object.GetNestedProperties() {
		field := message.GetField(key)
		if field == nil {
			return trace.New(trace.WithMessage("undefined schema nested message '%s' in flow '%s'", nested.Path, flow.GetName()))
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
			return trace.New(trace.WithMessage("undefined schema repeated message '%s' in flow '%s'", repeated.Path, flow.GetName()))
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
