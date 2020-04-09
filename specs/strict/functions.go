package strict

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/jexia/maestro/internal/instance"
	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/trace"
)

var (
	// FunctionPattern is the matching pattern for custom defined functions
	FunctionPattern = regexp.MustCompile(`(\w+)\((.*)\)$`)
)

const (
	// FunctionArgumentDelimiter represents the character delimiting function arguments
	FunctionArgumentDelimiter = ","
)

// PrepareManifestFunctions prepares all function definitions inside the given manifest
func PrepareManifestFunctions(ctx instance.Context, functions specs.CustomDefinedFunctions, manifest *specs.Manifest) (err error) {
	ctx.Logger(logger.Core).Info("Comparing manifest types")

	for _, flow := range manifest.Flows {
		err := PrepareFlowFunctions(ctx, functions, manifest, flow)
		if err != nil {
			return err
		}
	}

	for _, proxy := range manifest.Proxy {
		err := PrepareProxyFunctions(ctx, functions, manifest, proxy)
		if err != nil {
			return err
		}
	}

	return nil
}

// PrepareProxyFunctions prepares all function definitions inside the given proxy
func PrepareProxyFunctions(ctx instance.Context, functions specs.CustomDefinedFunctions, manifest *specs.Manifest, proxy *specs.Proxy) (err error) {
	ctx.Logger(logger.Core).WithField("proxy", proxy.GetName()).Info("Prepare proxy functions")

	for _, node := range proxy.Nodes {
		if node.Call != nil {
			err = PrepareCallFunctions(ctx, node, proxy, functions, node.Call)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = PrepareCallFunctions(ctx, node, proxy, functions, node.Rollback)
			if err != nil {
				return err
			}
		}
	}

	// if proxy.Forward.Request.Header != nil {
	// 	err = PrepareHeaderFunctions(proxy.Forward.Request.Header, proxy)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// PrepareFlowFunctions prepares the functions definitions inside the given flow
func PrepareFlowFunctions(ctx instance.Context, functions specs.CustomDefinedFunctions, manifest *specs.Manifest, flow *specs.Flow) (err error) {
	ctx.Logger(logger.Core).WithField("flow", flow.GetName()).Info("Comparing flow functions")

	for _, node := range flow.Nodes {
		if node.Call != nil {
			err = PrepareCallFunctions(ctx, node, flow, functions, node.Call)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = PrepareCallFunctions(ctx, node, flow, functions, node.Rollback)
			if err != nil {
				return err
			}
		}
	}

	if flow.Output != nil {
		err := PrepareParameterMapFunctions(ctx, nil, flow, flow.Output, functions)
		if err != nil {
			return err
		}

		if flow.Output.Header != nil {

		}
	}

	return nil
}

// PrepareCallFunctions prepares the function definitions inside the given flow
func PrepareCallFunctions(ctx instance.Context, node *specs.Node, flow specs.FlowManager, functions specs.CustomDefinedFunctions, call *specs.Call) error {
	if call.Request != nil {
		err := PrepareParameterMapFunctions(ctx, node, flow, call.Request, functions)
		if err != nil {
			return err
		}
	}

	if call.Response != nil {
		err := PrepareParameterMapFunctions(ctx, node, flow, call.Response, functions)
		if err != nil {
			return err
		}
	}

	return nil
}

// PrepareParameterMapFunctions prepares the function definitions inside the given parameter map
func PrepareParameterMapFunctions(ctx instance.Context, node *specs.Node, flow specs.FlowManager, params *specs.ParameterMap, functions specs.CustomDefinedFunctions) error {
	if params.Functions == nil {
		params.Functions = make(specs.Functions)
	}

	err := PreparePropertyFunctions(ctx, node, flow, params.Functions, params.Property, functions)
	if err != nil {
		return err
	}

	return nil
}

// PreparePropertyFunctions prepares the function definitions inside the given property
func PreparePropertyFunctions(ctx instance.Context, node *specs.Node, flow specs.FlowManager, collection specs.Functions, prop *specs.Property, functions specs.CustomDefinedFunctions) error {
	if prop == nil {
		return nil
	}

	if prop.Nested != nil {
		for _, nested := range prop.Nested {
			err := PreparePropertyFunctions(ctx, node, flow, collection, nested, functions)
			if err != nil {
				return err
			}
		}
	}

	prop, err := PrepareFunction(ctx, node, flow, prop, collection, functions)
	if err != nil {
		return err
	}

	if prop != nil {
		// override prop
	}

	return nil
}

// PrepareFunction attempts to parses the given function
func PrepareFunction(ctx instance.Context, node *specs.Node, flow specs.FlowManager, property *specs.Property, collection specs.Functions, methods specs.CustomDefinedFunctions) (*specs.Property, error) {
	if !FunctionPattern.MatchString(property.Raw) {
		return nil, nil
	}

	pattern := FunctionPattern.FindStringSubmatch(property.Raw)
	fn := pattern[1]
	args := strings.Split(pattern[2], FunctionArgumentDelimiter)

	if methods[fn] == nil {
		return nil, trace.New(trace.WithMessage("undefined custom function '%s' in '%s'", fn, property.Raw))
	}

	arguments := make([]*specs.Property, len(args))

	for index, arg := range args {
		result, err := specs.ParseTemplateContent(property.Path, property.Name, strings.TrimSpace(arg))
		if err != nil {
			return nil, err
		}

		err = DefineProperty(ctx, node, result, flow)
		if err != nil {
			return nil, err
		}

		arguments[index] = result
	}

	returns, handle, err := methods[fn](arguments...)
	if err != nil {
		return nil, err
	}

	stack := GenerateStackReference()
	function := &specs.Function{
		Arguments: arguments,
		Fn:        handle,
		Returns:   property,
	}

	collection[stack] = function

	result := &specs.Property{
		Name:    property.Name,
		Path:    property.Path,
		Type:    returns.Type,
		Label:   returns.Label,
		Default: returns.Default,
		Reference: &specs.PropertyReference{
			Resource: specs.JoinPath(specs.StackResource, stack),
			Path:     ".",
			Property: returns,
		},
		Raw: property.Raw,
	}

	return result, nil
}

// GenerateStackReference generates a unique path prefix which could be used to isolate functions
func GenerateStackReference() string {
	bb := make([]byte, 5)
	rand.Read(bb)
	return hex.EncodeToString(bb)
}
