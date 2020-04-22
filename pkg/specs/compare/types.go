package compare

import (
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/trace"
	"github.com/jexia/maestro/pkg/specs/types"
	"github.com/sirupsen/logrus"
)

// ManifestTypes compares the types defined insde the schema definitions against the configured specification
func ManifestTypes(ctx instance.Context, services *specs.ServicesManifest, schema *specs.SchemaManifest, flows *specs.FlowsManifest) (err error) {
	ctx.Logger(logger.Core).Info("Comparing manifest types")

	for _, flow := range flows.Flows {
		err := FlowTypes(ctx, services, schema, flows, flow)
		if err != nil {
			return err
		}
	}

	for _, proxy := range flows.Proxy {
		err := ProxyTypes(ctx, services, schema, flows, proxy)
		if err != nil {
			return err
		}
	}

	return nil
}

// ProxyTypes compares the given proxy against the configured schema types
func ProxyTypes(ctx instance.Context, services *specs.ServicesManifest, schema *specs.SchemaManifest, flows *specs.FlowsManifest, proxy *specs.Proxy) (err error) {
	ctx.Logger(logger.Core).WithField("proxy", proxy.GetName()).Info("Compare proxy flow types")

	for _, node := range proxy.Nodes {
		if node.Call != nil {
			err = CallTypes(ctx, services, schema, flows, node, node.Call, proxy)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = CallTypes(ctx, services, schema, flows, node, node.Rollback, proxy)
			if err != nil {
				return err
			}
		}
	}

	if proxy.Forward.Request.Header != nil {
		err = Header(proxy.Forward.Request.Header, proxy)
		if err != nil {
			return err
		}
	}

	return nil
}

// FlowTypes compares the flow types against the configured schema types
func FlowTypes(ctx instance.Context, services *specs.ServicesManifest, schema *specs.SchemaManifest, flows *specs.FlowsManifest, flow *specs.Flow) (err error) {
	ctx.Logger(logger.Core).WithField("flow", flow.GetName()).Info("Comparing flow types")

	if flow.Input != nil {
		message := schema.GetProperty(flow.Input.Schema)
		if message == nil {
			return trace.New(trace.WithMessage("undefined flow input object '%s' in '%s'", flow.Input.Schema, flow.Name))
		}

		err = CheckTypes(flow.Input.Property, message, flow)
		if err != nil {
			return err
		}
	}

	for _, node := range flow.Nodes {
		if node.Call != nil {
			err = CallTypes(ctx, services, schema, flows, node, node.Call, flow)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = CallTypes(ctx, services, schema, flows, node, node.Rollback, flow)
			if err != nil {
				return err
			}
		}
	}

	if flow.Output != nil {
		message := schema.GetProperty(flow.Output.Schema)
		if message == nil {
			return trace.New(trace.WithMessage("undefined flow output object '%s' in '%s'", flow.Output.Schema, flow.Name))
		}

		err = CheckTypes(flow.Output.Property, message, flow)
		if err != nil {
			return err
		}

		if flow.Output.Header != nil {
			err = Header(flow.Output.Header, flow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CallTypes compares the given call types against the configured schema types
func CallTypes(ctx instance.Context, services *specs.ServicesManifest, schema *specs.SchemaManifest, flows *specs.FlowsManifest, node *specs.Node, call *specs.Call, flow specs.FlowResourceManager) (err error) {
	if call.Method != "" {
		ctx.Logger(logger.Core).WithFields(logrus.Fields{
			"call":    node.Name,
			"method":  call.Method,
			"service": call.Service,
		}).Info("Comparing call types")

		service := services.GetService(call.Service)
		if service == nil {
			return trace.New(trace.WithMessage("undefined service '%s' in flow '%s'", call.Service, flow.GetName()))
		}

		method := service.GetMethod(call.Method)
		if method == nil {
			return trace.New(trace.WithMessage("undefined method '%s' in flow '%s'", call.Method, flow.GetName()))
		}

		if call.Request != nil {
			if call.Request.Header != nil {
				err = Header(call.Request.Header, flow)
				if err != nil {
					return err
				}
			}

			err = CheckTypes(call.Request.Property, schema.GetProperty(method.Input), flow)
			if err != nil {
				return err
			}
		}

		if call.Response != nil {
			if call.Response.Header != nil {
				err = Header(call.Request.Header, flow)
				if err != nil {
					return err
				}
			}

			err = CheckTypes(call.Response.Property, schema.GetProperty(method.Output), flow)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CheckTypes checks the given schema against the given schema method types
func CheckTypes(property *specs.Property, schema *specs.Property, flow specs.FlowResourceManager) (err error) {
	if schema == nil {
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("unable to check types for '%s' no schema given", property.Path))
	}

	if property.Type != schema.Type {
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("cannot use (%s) type (%s) in '%s'", property.Type, schema.Type, property.Path))
	}

	if property.Label != schema.Label {
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("cannot use (%s) label (%s) in '%s'", property.Label, schema.Label, property.Path))
	}

	if len(property.Nested) > 0 {
		if len(schema.Nested) == 0 {
			return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("property '%s' has a nested object but schema does not '%s'", property.Path, schema.Name))
		}

		for key, nested := range property.Nested {
			object := schema.Nested[key]
			if object == nil {
				return trace.New(trace.WithExpression(nested.Expr), trace.WithMessage("undefined schema nested message property '%s' in flow '%s'", nested.Path, flow.GetName()))
			}

			err := CheckTypes(nested, object, flow)
			if err != nil {
				return err
			}
		}

		// Set any properties not defined inside the flow but available inside the schema
		for _, prop := range schema.Nested {
			_, has := property.Nested[prop.Name]
			if has {
				continue
			}

			property.Nested[prop.Name] = prop
		}
	}

	// ensure the property position
	property.Position = schema.Position

	return nil
}

// Header compares the given header types
func Header(header specs.Header, flow specs.FlowResourceManager) error {
	for _, header := range header {
		if header.Type != types.String {
			return trace.New(trace.WithMessage("cannot use type (%s) for 'header.%s' in flow '%s', expected (%s)", header.Type, header.Path, flow.GetName(), types.String))
		}
	}

	return nil
}
