package providers

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/trace"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"go.uber.org/zap"
)

// ResolveSchemas ensures that all schema properties are defined inside the given flows
func ResolveSchemas(ctx *broker.Context, services specs.ServiceList, schemas specs.Schemas, flows specs.FlowListInterface) (err error) {
	logger.Info(ctx, "defining manifest types")

	for _, flow := range flows {
		err := ResolveFlow(ctx, services, schemas, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// ResolveFlow ensures that all schema properties are defined inside the given flow
func ResolveFlow(parent *broker.Context, services specs.ServiceList, schemas specs.Schemas, flow specs.FlowInterface) (err error) {
	ctx := logger.WithFields(parent, zap.String("flow", flow.GetName()))
	logger.Info(ctx, "defining flow types")

	if flow.GetInput() != nil {
		input := schemas.Get(flow.GetInput().Schema)
		if input == nil {
			return ErrUndefinedObject{
				Schema: flow.GetInput().Schema,
			}
		}

		flow.GetInput().Property = input.Clone()
	}

	if flow.GetOnError() != nil {
		err = ResolveOnError(ctx, schemas, flow.GetOnError(), flow)
		if err != nil {
			return err
		}
	}

	for _, node := range flow.GetNodes() {
		err = ResolveNode(ctx, services, schemas, node, flow)
		if err != nil {
			return err
		}
	}

	if flow.GetOutput() != nil {
		err = ResolveParameterMap(ctx, schemas, flow.GetOutput(), flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// ResolveNode ensures that all schema properties are defined inside the given node
func ResolveNode(ctx *broker.Context, services specs.ServiceList, schemas specs.Schemas, node *specs.Node, flow specs.FlowInterface) (err error) {
	if node.Condition != nil {
		err = ResolveParameterMap(ctx, schemas, node.Condition.Params, flow)
		if err != nil {
			return err
		}
	}

	if node.Call != nil {
		err = DefineCall(ctx, services, schemas, node, node.Call, flow)
		if err != nil {
			return err
		}
	}

	if node.Rollback != nil {
		err = DefineCall(ctx, services, schemas, node, node.Rollback, flow)
		if err != nil {
			return err
		}
	}

	if node.Intermediate != nil {
		err = ResolveParameterMap(ctx, schemas, node.Intermediate, flow)
		if err != nil {
			return err
		}
	}

	if node.OnError != nil {
		err = ResolveOnError(ctx, schemas, node.OnError, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineCall defineds the types for the specs call
func DefineCall(ctx *broker.Context, services specs.ServiceList, schemas specs.Schemas, node *specs.Node, call *specs.Call, flow specs.FlowInterface) (err error) {
	if call.Request != nil {
		err = ResolveParameterMap(ctx, schemas, call.Request, flow)
		if err != nil {
			return err
		}
	}

	if call.Method != "" {
		logger.Info(ctx, "defining call types",
			zap.String("call", node.ID),
			zap.String("method", call.Method),
			zap.String("service", call.Service),
		)

		service := services.Get(call.Service)
		if service == nil {
			return ErrUndefinedService{
				Service: call.Service,
				Flow:    flow.GetName(),
			}
		}

		method := service.GetMethod(call.Method)
		if method == nil {
			return ErrUndefinedMethod{
				Flow:   flow.GetName(),
				Method: call.Method,
			}
		}

		output := schemas.Get(method.Output)
		if output == nil {
			return ErrUndefinedOutput{
				Output: method.Output,
				Flow:   flow.GetName(),
			}
		}

		call.Descriptor = method
		call.Response = &specs.ParameterMap{
			Property: output.Clone(),
		}

		call.Request.Schema = method.Input
		call.Response.Schema = method.Output
	}

	if call.Response != nil {
		err = ResolveParameterMap(ctx, schemas, call.Response, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// ResolveParameterMap ensures that all schema properties are defined inisde the given parameter map
func ResolveParameterMap(ctx *broker.Context, schemas specs.Schemas, params *specs.ParameterMap, flow specs.FlowInterface) (err error) {
	if params == nil || params.Schema == "" {
		return nil
	}

	schema := schemas.Get(params.Schema)
	if schema == nil {
		return ErrUndefinedObject{
			Schema: params.Schema,
		}
	}

	err = ResolveProperty(params.Property, schema.Clone(), flow)
	if err != nil {
		return err
	}

	return nil
}

// ResolveOnError ensures that all schema properties are defined inside the given on error object
func ResolveOnError(ctx *broker.Context, schemas specs.Schemas, params *specs.OnError, flow specs.FlowInterface) (err error) {
	if params.Response != nil {
		err = ResolveParameterMap(ctx, schemas, params.Response, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

func resolveMessage(message, schema specs.Message, flow specs.FlowInterface) error {
	for _, nested := range message {
		if nested == nil {
			continue
		}

		object := schema[nested.Name]
		if object == nil {
			return trace.New(trace.WithMessage("undefined schema nested message property '%s' in flow '%s'", nested.Name, flow.GetName()))
		}

		if err := ResolveProperty(nested, object.Clone(), flow); err != nil {
			return err
		}
	}

	return nil
}

func resolveRepeated(repeated, schema specs.Repeated, flow specs.FlowInterface) error {
	if len(repeated) != len(schema) {
		return trace.New(trace.WithMessage("the length of repeated does not match the schema"))
	}

	// FIXME: flow and schema repeated could have different type orders.
	for pos, template := range repeated {
		object := schema[pos]

		if err := ResolveProperty(
			&specs.Property{
				Template: template,
			},
			(&specs.Property{
				Template: object,
			}).Clone(),
			flow,
		); err != nil {
			return err
		}
	}

	return nil
}

func setMessage(message, schema specs.Message) {
	// Set any properties not defined inside the flow but available inside the schema
	for _, prop := range schema {
		if _, ok := message[prop.Name]; ok {
			continue
		}

		message[prop.Name] = prop.Clone()
	}
}

func setRepeated(repeated, schema specs.Repeated) {
	// FIXME: how do we match schemas and repeated?
	// FIXME: out of range panic
	for pos, prop := range schema {
		repeated[pos] = prop.Clone()
	}
}

// ResolveProperty ensures that all schema properties are defined inside the given property
func ResolveProperty(property, schema *specs.Property, flow specs.FlowInterface) error {
	if property == nil {
		property = schema.Clone()
		return nil
	}

	switch {
	case property.Message != nil:
		if err := resolveMessage(property.Message, schema.Message, flow); err != nil {
			return ErrUndefinedProperty{
				Property: property.Name,
				Flow:     flow.GetName(),
			}
		}

		property.Label = schema.Label

		break
	case property.Repeated != nil:
		if err := resolveRepeated(property.Repeated, schema.Repeated, flow); err != nil {
			return err
		}

		property.Label = schema.Label

		break
	}

	switch {
	case schema.Message != nil:
		if property.Message == nil {
			property.Message = schema.Message.Clone()
		}

		setMessage(property.Message, schema.Message)

		break
	case schema.Repeated != nil:
		if property.Repeated == nil {
			property.Repeated = schema.Repeated.Clone()
		}

		setRepeated(property.Repeated, schema.Repeated)

		break
	}

	return nil
}

func ConstructReferencedPathsProperty(referenced Paths, property *specs.Property) *specs.Property {
	paths := allPosssiblePaths(referenced)
	result := property.Clone()
	constructReferencedPathsTemplate(paths, "", result.Template)

	return result
}

func constructReferencedPathsTemplate(referenced Paths, path string, template specs.Template) {
	switch {
	case template.Message != nil:
		for key, nested := range template.Message {
			if nested == nil {
				delete(template.Message, key)
				continue
			}

			if _, has := referenced[nested.Path]; !has {
				delete(template.Message, key)
				continue
			}

			constructReferencedPathsTemplate(referenced, nested.Path, nested.Template)
		}
	case template.Repeated != nil:
		if _, has := referenced[path]; !has {
			template.Repeated = nil
			break
		}

		for _, item := range template.Repeated {
			constructReferencedPathsTemplate(referenced, path, item)
		}
	}
}

func allPosssiblePaths(referenced Paths) Paths {
	result := make(Paths, len(referenced))

	for path := range referenced {
		absolute := ""
		parts := template.SplitPath(path)
		for _, part := range parts {
			current := template.JoinPath(absolute, part)
			result[current] = struct{}{}
			absolute = current
		}
	}

	return result
}

type Paths map[string]struct{}

func FlowReferencedResourcePaths(target Paths, flow specs.FlowInterface, resource string) {
	for _, node := range flow.GetNodes() {
		if node.Call != nil {
			ParameterMapReferencedResourcePaths(target, node.Call.Request, resource)
			ParameterMapReferencedResourcePaths(target, node.Call.Response, resource)
		}

		ParameterMapReferencedResourcePaths(target, node.Intermediate, resource)

		if node.Condition != nil {
			ParameterMapReferencedResourcePaths(target, node.Condition.Params, resource)
		}
	}

	ParameterMapReferencedResourcePaths(target, flow.GetOutput(), resource)
}

func ParameterMapReferencedResourcePaths(target Paths, parameters *specs.ParameterMap, resource string) {
	if parameters == nil {
		return
	}

	for _, header := range parameters.Header {
		if header == nil {
			continue
		}

		PropertyReferencedResourcePaths(target, header.Template, resource)
	}

	for _, params := range parameters.Params {
		if params == nil {
			continue
		}

		PropertyReferencedResourcePaths(target, params.Template, resource)
	}

	if parameters.Property != nil {
		PropertyReferencedResourcePaths(target, parameters.Property.Template, resource)
	}

	for _, stack := range parameters.Stack {
		if stack == nil {
			continue
		}

		PropertyReferencedResourcePaths(target, stack.Template, resource)
	}
}

func PropertyReferencedResourcePaths(target Paths, template specs.Template, resource string) {
	if template.Reference != nil && template.Reference.Resource == resource {
		target[template.Reference.Path] = struct{}{}
	}

	switch {
	case template.Message != nil:
		for _, nested := range template.Message {
			PropertyReferencedResourcePaths(target, nested.Template, resource)
		}
	case template.Repeated != nil:
		for _, item := range template.Repeated {
			PropertyReferencedResourcePaths(target, item, resource)
		}
	}
}
