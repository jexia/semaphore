package strict

import (
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/lookup"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/specs/types"
	log "github.com/sirupsen/logrus"
)

// Define checks and defines the types for the given manifest
func Define(schema schema.Collection, manifest *specs.Manifest) (err error) {
	log.Info("Defining manifest types")

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
	log.WithField("proxy", proxy.GetName()).Info("Defining proxy flow types")

	for _, node := range proxy.Nodes {
		if node.Call != nil {
			err = DefineCall(schema, manifest, node, node.Call, proxy)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = DefineCall(schema, manifest, node, node.Rollback, proxy)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DefineFlow checks and defines the types for the given flow
func DefineFlow(schema schema.Collection, manifest *specs.Manifest, flow *specs.Flow) (err error) {
	log.WithField("flow", flow.GetName()).Info("Defining flow types")

	if flow.Input != nil {
		message, err := GetObjectSchema(schema, flow.Input)
		if err != nil {
			return err
		}

		flow.Input = specs.ToParameterMap(flow.Input, "", message)
		flow.Input.SetDescriptor(message)
	}

	for _, node := range flow.Nodes {
		if node.Call != nil {
			err = DefineCall(schema, manifest, node, node.Call, flow)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = DefineCall(schema, manifest, node, node.Rollback, flow)
			if err != nil {
				return err
			}
		}
	}

	if flow.Output != nil {
		err = DefineParameterMap(nil, nil, flow.Output, flow)
		if err != nil {
			return err
		}

		message, err := GetObjectSchema(schema, flow.Output)
		if err != nil {
			return err
		}

		err = CheckTypes(flow.Output, message, flow)
		if err != nil {
			return err
		}

		flow.Output.SetDescriptor(message)
	}

	return nil
}

// GetObjectSchema attempts to fetch the defined schema object for the given parameter map
func GetObjectSchema(schema schema.Collection, params *specs.ParameterMap) (schema.Object, error) {
	object := schema.GetObject(params.Schema)
	if object == nil {
		return nil, trace.New(trace.WithMessage("undefined object '%s' in schema collection", params.Schema))
	}

	return object, nil
}

// DefineCall defineds the types for the given parameter map
func DefineCall(schema schema.Collection, manifest *specs.Manifest, node *specs.Node, call *specs.Call, flow specs.FlowManager) (err error) {
	if call.GetMethod() == "" {
		return nil
	}

	log.WithFields(log.Fields{
		"call":   node.GetName(),
		"method": call.GetMethod(),
	}).Info("Defining call types")

	service := schema.GetService(GetSchemaService(manifest, call.GetService()))
	if service == nil {
		return trace.New(trace.WithMessage("undefined service alias '%s' in flow '%s'", call.GetService(), flow.GetName()))
	}

	method := service.GetMethod(call.GetMethod())
	if method == nil {
		return trace.New(trace.WithMessage("undefined method '%s' in flow '%s'", call.GetMethod(), flow.GetName()))
	}

	call.SetDescriptor(method)

	if call.GetRequest() != nil {
		err = DefineParameterMap(node, call, call.GetRequest(), flow)
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
func DefineParameterMap(node *specs.Node, call *specs.Call, params specs.Object, flow specs.FlowManager) (err error) {
	for _, header := range params.GetHeader() {
		err = DefineProperty(node, call, header, flow)
		if err != nil {
			return err
		}
	}

	for _, property := range params.GetProperties() {
		err = DefineProperty(node, call, property, flow)
		if err != nil {
			return err
		}
	}

	for _, nested := range params.GetNestedProperties() {
		err = DefineParameterMap(node, call, nested, flow)
		if err != nil {
			return err
		}
	}

	for _, repeated := range params.GetRepeatedProperties() {
		err = DefineParameterMap(node, call, repeated, flow)
		if err != nil {
			return err
		}
	}

	ResolvePropertyObjectReferences(params)

	return nil
}

// DefineProperty defines the given property type.
// If any object is references it has to be fixed afterwards and moved into the correct dataset
func DefineProperty(node *specs.Node, call *specs.Call, property *specs.Property, flow specs.FlowManager) error {
	if property.Reference == nil {
		return nil
	}

	breakpoint := "output"
	if node != nil {
		breakpoint = node.GetName()
	}

	log.WithFields(log.Fields{
		"breakpoint": breakpoint,
		"reference":  property.Reference,
	}).Debug("Lookup references until breakpoint")

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
	object.SetDescriptor(message)

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
			clone := repeated.Clone(key, property.Path)
			clone.Template = property.Reference

			SetObjectReferences(property.Reference, clone)
			params.GetRepeatedProperties()[key] = clone
			continue
		}

		nested := property.Reference.Object.(*specs.NestedParameterMap)
		params.GetNestedProperties()[key] = nested.Clone(key, property.Path)
	}
}

// SetObjectReferences sets all the references within a given object to the given resource and path
func SetObjectReferences(ref *specs.PropertyReference, object specs.Object) {
	for key, prop := range object.GetProperties() {
		ref := ref.Clone()
		ref.Path = specs.JoinPath(ref.Path, key)
		prop.Reference = ref
	}

	for key, nested := range object.GetNestedProperties() {
		ref := ref.Clone()
		ref.Path = specs.JoinPath(ref.Path, key)
		SetObjectReferences(ref, nested)
	}

	for key, repeated := range object.GetRepeatedProperties() {
		ref := ref.Clone()
		ref.Path = specs.JoinPath(ref.Path, key)
		SetObjectReferences(ref, repeated)
	}
}

// GetSchemaService attempts to find a service matching the alias name and return the schema name
func GetSchemaService(manifest *specs.Manifest, name string) string {
	for _, service := range manifest.Services {
		if service.Name == name {
			return service.Schema
		}
	}

	return ""
}
