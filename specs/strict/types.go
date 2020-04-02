package strict

import (
	"github.com/jexia/maestro/instance"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/specs/types"
	"github.com/sirupsen/logrus"
)

// CompareManifestTypes compares the types defined insde the schema definitions agains the configured specification
func CompareManifestTypes(ctx instance.Context, schema schema.Collection, manifest *specs.Manifest) (err error) {
	ctx.Logger(logger.Core).Info("Comparing manifest types")

	for _, flow := range manifest.Flows {
		err := CompareFlowTypes(ctx, schema, manifest, flow)
		if err != nil {
			return err
		}
	}

	for _, proxy := range manifest.Proxy {
		err := CompareProxyTypes(ctx, schema, manifest, proxy)
		if err != nil {
			return err
		}
	}

	return nil
}

// CompareProxyTypes compares the given proxy agains the configured schema types
func CompareProxyTypes(ctx instance.Context, schema schema.Collection, manifest *specs.Manifest, proxy *specs.Proxy) (err error) {
	ctx.Logger(logger.Core).WithField("proxy", proxy.GetName()).Info("Compare proxy flow types")

	for _, node := range proxy.Nodes {
		if node.Call != nil {
			err = CompareCallTypes(ctx, schema, manifest, node, node.Call, proxy)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = CompareCallTypes(ctx, schema, manifest, node, node.Rollback, proxy)
			if err != nil {
				return err
			}
		}
	}

	// TODO: proxy header type checking

	return nil
}

// CompareFlowTypes compares the flow types agains the configured schema types
func CompareFlowTypes(ctx instance.Context, schema schema.Collection, manifest *specs.Manifest, flow *specs.Flow) (err error) {
	ctx.Logger(logger.Core).WithField("flow", flow.GetName()).Info("Comparing flow types")

	if flow.Input != nil {
		message, err := GetObjectSchema(schema, flow.Input)
		if err != nil {
			return err
		}

		err = CheckTypes(flow.Input.Property, message, flow)
		if err != nil {
			return err
		}
	}

	for _, node := range flow.Nodes {
		if node.Call != nil {
			err = CompareCallTypes(ctx, schema, manifest, node, node.Call, flow)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = CompareCallTypes(ctx, schema, manifest, node, node.Rollback, flow)
			if err != nil {
				return err
			}
		}
	}

	if flow.Output != nil {
		message, err := GetObjectSchema(schema, flow.Output)
		if err != nil {
			return err
		}

		err = CheckTypes(flow.Output.Property, message, flow)
		if err != nil {
			return err
		}

		if flow.Output.Header != nil {
			err = CompareHeader(flow.Output.Header, flow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CompareCallTypes compares the given call types agains the configured schema types
func CompareCallTypes(ctx instance.Context, schema schema.Collection, manifest *specs.Manifest, node *specs.Node, call *specs.Call, flow specs.FlowManager) (err error) {
	if call.Method != "" {
		ctx.Logger(logger.Core).WithFields(logrus.Fields{
			"call":    node.Name,
			"method":  call.Method,
			"service": call.Service,
		}).Info("Comparing call types")

		service := schema.GetService(call.Service)
		if service == nil {
			return trace.New(trace.WithMessage("undefined service '%s' in flow '%s'", call.Service, flow.GetName()))
		}

		method := service.GetMethod(call.Method)
		if method == nil {
			return trace.New(trace.WithMessage("undefined method '%s' in flow '%s'", call.Method, flow.GetName()))
		}

		if call.Request != nil {
			if call.Request.Header != nil {
				err = CompareHeader(call.Request.Header, flow)
				if err != nil {
					return err
				}
			}

			err = CheckTypes(call.Request.Property, method.GetInput(), flow)
			if err != nil {
				return err
			}
		}

		if call.Response != nil {
			if call.Response.Header != nil {
				err = CompareHeader(call.Request.Header, flow)
				if err != nil {
					return err
				}
			}

			err = CheckTypes(call.Response.Property, method.GetOutput(), flow)
			if err != nil {
				return err
			}
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

	if len(property.Nested) > 0 {
		if len(schema.GetNested()) == 0 {
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

// CompareHeader compares the given header types
func CompareHeader(header specs.Header, flow specs.FlowManager) error {
	for _, header := range header {
		if header.Type != types.String {
			return trace.New(trace.WithMessage("cannot use type %s for header.%s in flow %s", header.Type, header.Path, flow.GetName()))
		}
	}

	return nil
}
