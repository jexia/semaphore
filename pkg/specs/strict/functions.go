package strict

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/template"
	"github.com/jexia/maestro/pkg/specs/trace"
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
func PrepareManifestFunctions(ctx instance.Context, mem functions.Collection, functions functions.Custom, manifest *specs.FlowsManifest) (err error) {
	ctx.Logger(logger.Core).Info("Comparing manifest types")

	for _, flow := range manifest.Flows {
		err := PrepareFlowFunctions(ctx, mem, functions, manifest, flow)
		if err != nil {
			return err
		}
	}

	for _, proxy := range manifest.Proxy {
		err := PrepareProxyFunctions(ctx, mem, functions, manifest, proxy)
		if err != nil {
			return err
		}
	}

	return nil
}

// PrepareProxyFunctions prepares all function definitions inside the given proxy
func PrepareProxyFunctions(ctx instance.Context, mem functions.Collection, functions functions.Custom, manifest *specs.FlowsManifest, proxy *specs.Proxy) (err error) {
	ctx.Logger(logger.Core).WithField("proxy", proxy.GetName()).Info("Prepare proxy functions")

	for _, node := range proxy.Nodes {
		if node.Call != nil {
			err = PrepareCallFunctions(ctx, node, proxy, mem, functions, node.Call)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = PrepareCallFunctions(ctx, node, proxy, mem, functions, node.Rollback)
			if err != nil {
				return err
			}
		}
	}

	if proxy.Forward.Request.Header != nil {
		stack := mem.Reserve(nil)
		err = PrepareHeaderFunctions(ctx, proxy, stack, proxy.Forward.Request.Header, functions)
		if err != nil {
			return err
		}
	}

	return nil
}

// PrepareFlowFunctions prepares the functions definitions inside the given flow
func PrepareFlowFunctions(ctx instance.Context, mem functions.Collection, functions functions.Custom, manifest *specs.FlowsManifest, flow *specs.Flow) (err error) {
	ctx.Logger(logger.Core).WithField("flow", flow.GetName()).Info("Comparing flow functions")

	for _, node := range flow.Nodes {
		if node.Call != nil {
			err = PrepareCallFunctions(ctx, node, flow, mem, functions, node.Call)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = PrepareCallFunctions(ctx, node, flow, mem, functions, node.Rollback)
			if err != nil {
				return err
			}
		}
	}

	if flow.Output != nil {
		stack := mem.Reserve(nil)
		err := PrepareParameterMapFunctions(ctx, nil, flow, stack, flow.Output, functions)
		if err != nil {
			return err
		}

		if flow.Output.Header != nil {

		}
	}

	return nil
}

// PrepareCallFunctions prepares the function definitions inside the given flow
func PrepareCallFunctions(ctx instance.Context, node *specs.Node, flow specs.FlowResourceManager, mem functions.Collection, functions functions.Custom, call *specs.Call) error {
	if call.Request != nil {
		stack := mem.Reserve(node)
		err := PrepareParameterMapFunctions(ctx, node, flow, stack, call.Request, functions)
		if err != nil {
			return err
		}
	}

	if call.Response != nil {
		stack := mem.Reserve(node)
		err := PrepareParameterMapFunctions(ctx, node, flow, stack, call.Response, functions)
		if err != nil {
			return err
		}
	}

	return nil
}

// PrepareParameterMapFunctions prepares the function definitions inside the given parameter map
func PrepareParameterMapFunctions(ctx instance.Context, node *specs.Node, flow specs.FlowResourceManager, stack functions.Stack, params *specs.ParameterMap, functions functions.Custom) error {
	if params.Header != nil {
		err := PrepareHeaderFunctions(ctx, flow, stack, params.Header, functions)
		if err != nil {
			return err
		}
	}

	if params.Property != nil {
		err := PreparePropertyFunctions(ctx, node, flow, stack, params.Property, functions)
		if err != nil {
			return err
		}
	}

	return nil
}

// PrepareHeaderFunctions prepares the function definitions inside the given header
func PrepareHeaderFunctions(ctx instance.Context, flow specs.FlowResourceManager, stack functions.Stack, header specs.Header, functions functions.Custom) error {
	for _, prop := range header {
		err := PrepareFunction(ctx, nil, flow, prop, stack, functions)
		if err != nil {
			return err
		}
	}

	return nil
}

// PreparePropertyFunctions prepares the function definitions inside the given property
func PreparePropertyFunctions(ctx instance.Context, node *specs.Node, flow specs.FlowResourceManager, stack functions.Stack, prop *specs.Property, functions functions.Custom) error {
	if prop == nil {
		return nil
	}

	if prop.Nested != nil {
		for _, nested := range prop.Nested {
			err := PreparePropertyFunctions(ctx, node, flow, stack, nested, functions)
			if err != nil {
				return err
			}
		}
	}

	err := PrepareFunction(ctx, node, flow, prop, stack, functions)
	if err != nil {
		return err
	}

	return nil
}

// PrepareFunction attempts to parses the given function
func PrepareFunction(ctx instance.Context, node *specs.Node, flow specs.FlowResourceManager, property *specs.Property, stack functions.Stack, methods functions.Custom) error {
	if !FunctionPattern.MatchString(property.Raw) {
		return nil
	}

	pattern := FunctionPattern.FindStringSubmatch(property.Raw)
	fn := pattern[1]
	args := strings.Split(pattern[2], FunctionArgumentDelimiter)

	if methods[fn] == nil {
		return trace.New(trace.WithMessage("undefined custom function '%s' in '%s'", fn, property.Raw))
	}

	arguments := make([]*specs.Property, len(args))

	for index, arg := range args {
		result, err := template.ParseContent(property.Path, property.Name, strings.TrimSpace(arg))
		if err != nil {
			return err
		}

		err = DefineProperty(ctx, node, result, flow)
		if err != nil {
			return err
		}

		arguments[index] = result
	}

	returns, handle, err := methods[fn](arguments...)
	if err != nil {
		return err
	}

	ref := GenerateStackReference()
	function := &functions.Function{
		Arguments: arguments,
		Fn:        handle,
		Returns:   property,
	}

	stack[ref] = function

	property.Type = returns.Type
	property.Label = returns.Label
	property.Default = returns.Default
	property.Reference = &specs.PropertyReference{
		Resource: template.JoinPath(template.StackResource, ref),
		Path:     ".",
		Property: returns,
	}

	return nil
}

// GenerateStackReference generates a unique path prefix which could be used to isolate functions
func GenerateStackReference() string {
	bb := make([]byte, 5)
	rand.Read(bb)
	return hex.EncodeToString(bb)
}
