package strict

import (
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/lookup"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/specs/types"
	log "github.com/sirupsen/logrus"
)

// DefineManifest checks and defines the types for the given manifest
func DefineManifest(schema schema.Collection, manifest *specs.Manifest) (err error) {
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
		err = DefineParameterMap(nil, flow.Output, flow)
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
	prop := schema.GetMessage(params.Schema)
	if prop == nil {
		return nil, trace.New(trace.WithMessage("undefined object '%s' in schema collection", params.Schema))
	}

	return prop, nil
}

// DefineCall defineds the types for the specs call
func DefineCall(schema schema.Collection, manifest *specs.Manifest, node *specs.Node, call *specs.Call, flow specs.FlowManager) (err error) {
	if call.GetMethod() == "" {
		return nil
	}

	log.WithFields(log.Fields{
		"call":    node.GetName(),
		"method":  call.GetMethod(),
		"service": call.GetService(),
	}).Info("Defining call types")

	service := schema.GetService(call.GetService())
	if service == nil {
		return trace.New(trace.WithMessage("undefined service '%s' in flow '%s'", call.GetService(), flow.GetName()))
	}

	method := service.GetMethod(call.GetMethod())
	if method == nil {
		return trace.New(trace.WithMessage("undefined method '%s' in flow '%s'", call.GetMethod(), flow.GetName()))
	}

	call.SetDescriptor(method)

	if call.GetRequest() != nil {
		err = DefineParameterMap(node, call.GetRequest(), flow)
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

		err = CheckTypes(call.GetResponse().Property, method.GetOutput(), flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineCaller defineds the types for the given protocol caller
func DefineCaller(node *specs.Node, manifest *specs.Manifest, call protocol.Call, flow specs.FlowManager) (err error) {
	log.Info("Defining caller references")

	method := call.GetMethod(node.Call.GetMethod())
	for _, prop := range method.References() {
		err = DefineProperty(node, prop, flow)
		if err != nil {
			return err
		}

		ResolvePropertyReferences(prop)
	}

	return nil
}

// DefineParameterMap defines the types for the given parameter map
func DefineParameterMap(node *specs.Node, params *specs.ParameterMap, flow specs.FlowManager) (err error) {
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

	breakpoint := specs.OutputResource
	if node != nil {
		breakpoint = node.GetName()

		if node.Rollback != nil {
			rollback := node.Rollback.GetRequest().Property
			if InsideProperty(rollback, property) {
				breakpoint = lookup.GetNextResource(flow, breakpoint)
			}
		}
	}

	log.WithFields(log.Fields{
		"breakpoint": breakpoint,
		"reference":  property.Reference,
	}).Debug("Lookup references until breakpoint")

	references := lookup.GetAvailableResources(flow, breakpoint)
	reference := lookup.GetResourceReference(property.Reference, references, breakpoint)
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

// InsideProperty checks whether the given property is insde the source property
func InsideProperty(source *specs.Property, target *specs.Property) bool {
	if source == target {
		return true
	}

	if source.Nested != nil {
		for _, nested := range source.Nested {
			is := InsideProperty(nested, target)
			if is {
				return is
			}
		}
	}

	return false
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
	if schema == nil {
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("unable to check types for '%s' no schema given", property.Path))
	}

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

		for _, prop := range schema.GetNested() {
			_, has := property.Nested[prop.GetName()]
			if has {
				continue
			}

			property.Nested[prop.GetName()] = SchemaToProperty(property.Path, prop)
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

		return
	}

	if property.Reference == nil {
		return
	}

	if property.Reference.Property == nil {
		return
	}

	clone := property.Reference.Property.Clone(property.Reference, property.Name, property.Path)
	property.Reference = clone.Reference
	property.Nested = clone.Nested
}

// SchemaToProperty parses the given schema property to a specs property
func SchemaToProperty(path string, prop schema.Property) *specs.Property {
	result := &specs.Property{
		Name:      prop.GetName(),
		Path:      specs.JoinPath(path, prop.GetName()),
		Type:      prop.GetType(),
		Label:     prop.GetLabel(),
		Desciptor: prop,
	}

	if len(prop.GetNested()) > 0 {
		result.Nested = make(map[string]*specs.Property, len(prop.GetNested()))

		for key, prop := range prop.GetNested() {
			result.Nested[key] = SchemaToProperty(specs.JoinPath(path, key), prop)
		}
	}

	return result
}
