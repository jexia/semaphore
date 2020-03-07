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

	// TODO: proxy header type checking

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

		err = CheckHeader(flow.Output.Header, flow)
		if err != nil {
			return err
		}

		err = CheckTypes(flow.Output.Property, message, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetObjectSchema attempts to fetch the defined schema object for the given parameter map
func GetObjectSchema(schema schema.Collection, params *specs.ParameterMap) (schema.Property, error) {
	prop := schema.GetProperty(params.Schema)
	if prop == nil {
		return nil, trace.New(trace.WithMessage("undefined object '%s' in schema collection", params.Schema))
	}

	return prop, nil
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

		err = CheckHeader(call.GetRequest().Header, flow)
		if err != nil {
			return err
		}

		err = CheckTypes(call.GetRequest().Property, method.GetInput(), flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineParameterMap defines the types for the given parameter map
func DefineParameterMap(node *specs.Node, call *specs.Call, params *specs.ParameterMap, flow specs.FlowManager) (err error) {
	for _, header := range params.Header {
		err = DefineProperty(node, header, flow)
		if err != nil {
			return err
		}
	}

	err = DefineProperty(node, params.Property, flow)
	if err != nil {
		return err
	}

	ResolvePropertyReferences(params.Property)

	return nil
}

// DefineProperty defines the given property type.
// If any object is references it has to be fixed afterwards and moved into the correct dataset
func DefineProperty(node *specs.Node, property *specs.Property, flow specs.FlowManager) error {
	if property.Nested != nil {
		for _, nested := range property.Nested {
			err := DefineProperty(node, nested, flow)
			if err != nil {
				return err
			}
		}
	}

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

	property.Type = reference.Type
	property.Label = reference.Label
	property.Default = reference.Default
	property.Reference.Property = reference

	// TODO: support enum type

	return nil
}

// CheckHeader checks the given header types
func CheckHeader(header specs.Header, flow specs.FlowManager) error {
	for _, header := range header {
		if header.Type != types.TypeString {
			return trace.New(trace.WithMessage("cannot use type %s for header.%s in flow %s", header.Type, header.Path, flow.GetName()))
		}
	}

	return nil
}

// CheckTypes checks the given schema against the given schema method types
func CheckTypes(property *specs.Property, schema schema.Property, flow specs.FlowManager) (err error) {
	property.Desciptor = schema

	if property.Type != schema.GetType() {
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("cannot use (%s) type (%s) in '%s'", property.Type, schema.GetType(), property.Path))
	}

	if property.Label != schema.GetLabel() {
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("cannot use (%s) label (%s) in '%s'", property.Label, schema.GetLabel(), property.Path))
	}

	if property.Nested != nil {
		if schema.GetNested() == nil {
			return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("property '%s' has a nested object but schema does not '%s'", property.Path, schema.GetName()))
		}

		for key, nested := range property.Nested {
			object := schema.GetNested()[key]
			if object == nil {
				return trace.New(trace.WithExpression(nested.Expr), trace.WithMessage("undefined schema nested message property '%s' in flow '%s'", nested.Path, flow.GetName()))
			}

			err := CheckTypes(nested, object, flow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ResolvePropertyReferences moves any property reference into the correct data structure
func ResolvePropertyReferences(property *specs.Property) {
	if property.Nested != nil {
		for _, nested := range property.Nested {
			ResolvePropertyReferences(nested)
		}
	}

	if property.Reference == nil {
		return
	}

	if property.Reference.Property == nil {
		return
	}

	// reference := property.Reference.Property

	// if reference.Label == types.LabelRepeated {
	// 	clone := property.Clone(key, property.Path)
	// 	clone.Reference = property.Reference

	// 	SetReferences(property.Reference, clone)
	// 	params.GetRepeatedProperties()[key] = clone
	// 	return
	// }

	// ref := ref.Clone()
	// ref.Path = specs.JoinPath(ref.Path, key)
	// prop.Reference = ref
}

// // SetReferences sets all the references within a given object to the given resource and path
// func SetReferences(ref *specs.PropertyReference, object *specs.Property) {
// 	for key, prop := range object.GetProperties() {
// 		ref := ref.Clone()
// 		ref.Path = specs.JoinPath(ref.Path, key)
// 		prop.Reference = ref
// 	}

// 	for key, nested := range object.GetNestedProperties() {
// 		ref := ref.Clone()
// 		ref.Path = specs.JoinPath(ref.Path, key)
// 		SetReferences(ref, nested)
// 	}

// 	for key, repeated := range object.GetRepeatedProperties() {
// 		ref := ref.Clone()
// 		ref.Path = specs.JoinPath(ref.Path, key)
// 		SetReferences(ref, repeated)
// 	}
// }

// GetSchemaService attempts to find a service matching the alias name and return the schema name
func GetSchemaService(manifest *specs.Manifest, name string) string {
	for _, service := range manifest.Services {
		if service.Name == name {
			return service.Schema
		}
	}

	return ""
}
