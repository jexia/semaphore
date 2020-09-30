package functions

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"go.uber.org/zap"
)

// Custom represents a collection of custom defined functions that could be called inside a template
type Custom map[string]Intermediate

// Intermediate prepares the custom defined function.
// The given arguments represent the exprected types that are passed when called.
// Properties returned should be absolute.
type Intermediate func(args ...*specs.Property) (*specs.Property, Exec, error)

// Exec is a executable function.
// A store should be returned which could be used to encode the function property
type Exec func(store references.Store) error

// Collection represents a collection of stacks grouped by nodes
type Collection map[*specs.ParameterMap]Stack

// Reserve reserves a new function stack for the given node.
// If a stack already exists for the given node is it returned.
func (collection Collection) Reserve(params *specs.ParameterMap) Stack {
	stack, has := collection[params]
	if has {
		return stack
	}

	collection[params] = Stack{}
	return collection[params]
}

// Load attempts to load the function stack for the given parameter map
func (collection Collection) Load(params *specs.ParameterMap) Stack {
	return collection[params]
}

// Stack represents a collection of functions
type Stack map[string]*Function

// Function represents a custom defined function
type Function struct {
	Arguments []*specs.Property
	Fn        Exec
	Returns   *specs.Property
}

var (
	// FunctionPattern is the matching pattern for custom defined functions
	FunctionPattern = regexp.MustCompile(`(\w+)\((.*)\)$`)
)

const (
	// ArgumentDelimiter represents the character delimiting function arguments
	ArgumentDelimiter = ","
)

// PrepareFunctions prepares all function definitions inside the given flows
func PrepareFunctions(ctx *broker.Context, mem Collection, functions Custom, flows specs.FlowListInterface) (err error) {
	logger.Info(ctx, "preparing manifest functions")

	for _, flow := range flows {
		err := PrepareFlowFunctions(logger.WithFields(ctx, zap.String("flow", flow.GetName())), mem, functions, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// PrepareFlowFunctions prepares the functions definitions inside the given flow
func PrepareFlowFunctions(ctx *broker.Context, mem Collection, functions Custom, flow specs.FlowInterface) (err error) {
	logger.Info(ctx, "comparing flow functions")

	for _, node := range flow.GetNodes() {
		err = PrepareNodeFunctions(ctx, mem, functions, flow, node)
		if err != nil {
			return err
		}
	}

	if flow.GetOutput() != nil {
		stack := mem.Reserve(flow.GetOutput())
		err := PrepareParameterMapFunctions(ctx, nil, flow, stack, flow.GetOutput(), functions)
		if err != nil {
			return err
		}
	}

	if flow.GetForward() != nil {
		if flow.GetForward().Request != nil && flow.GetForward().Request.Header != nil {
			stack := mem.Reserve(flow.GetForward().Request)
			err = PrepareHeaderFunctions(ctx, flow, stack, flow.GetForward().Request.Header, functions)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// PrepareNodeFunctions prepares the available functions within the given node
func PrepareNodeFunctions(ctx *broker.Context, mem Collection, functions Custom, flow specs.FlowInterface, node *specs.Node) (err error) {
	if node.Intermediate != nil {
		stack := mem.Reserve(node.Intermediate)
		err = PrepareParameterMapFunctions(ctx, node, flow, stack, node.Intermediate, functions)
		if err != nil {
			return err
		}
	}

	if node.Condition != nil {
		stack := mem.Reserve(node.Condition.Params)
		err = PrepareParameterMapFunctions(ctx, node, flow, stack, node.Condition.Params, functions)
		if err != nil {
			return err
		}
	}

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

	return nil
}

// PrepareCallFunctions prepares the function definitions inside the given flow
func PrepareCallFunctions(ctx *broker.Context, node *specs.Node, flow specs.FlowInterface, mem Collection, functions Custom, call *specs.Call) error {
	if call.Request != nil {
		stack := mem.Reserve(call.Request)
		err := PrepareParameterMapFunctions(ctx, node, flow, stack, call.Request, functions)
		if err != nil {
			return err
		}
	}

	if call.Response != nil {
		stack := mem.Reserve(call.Response)
		err := PrepareParameterMapFunctions(ctx, node, flow, stack, call.Response, functions)
		if err != nil {
			return err
		}
	}

	return nil
}

// PrepareParameterMapFunctions prepares the function definitions inside the given parameter map
func PrepareParameterMapFunctions(ctx *broker.Context, node *specs.Node, flow specs.FlowInterface, stack Stack, params *specs.ParameterMap, functions Custom) error {
	if params.Header != nil {
		err := PrepareHeaderFunctions(ctx, flow, stack, params.Header, functions)
		if err != nil {
			return err
		}
	}

	if params.Params != nil {
		err := PrepareParamsFunctions(ctx, node, flow, stack, params.Params, functions)
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

	params.Stack = make(map[string]*specs.Property, len(stack))
	for ref, fn := range stack {
		params.Stack[ref] = fn.Returns
	}

	return nil
}

// PrepareHeaderFunctions prepares the function definitions inside the given header
func PrepareHeaderFunctions(ctx *broker.Context, flow specs.FlowInterface, stack Stack, header specs.Header, functions Custom) error {
	for _, prop := range header {
		err := PrepareFunction(ctx, nil, flow, prop, stack, functions)
		if err != nil {
			return err
		}
	}

	return nil
}

// PrepareParamsFunctions prepares the function definitions inside the given property
func PrepareParamsFunctions(ctx *broker.Context, node *specs.Node, flow specs.FlowInterface, stack Stack, params map[string]*specs.Property, functions Custom) error {
	if params == nil {
		return nil
	}

	for _, param := range params {
		err := PrepareFunction(ctx, node, flow, param, stack, functions)
		if err != nil {
			return err
		}
	}

	return nil
}

// PreparePropertyFunctions prepares the function definitions inside the given property
func PreparePropertyFunctions(ctx *broker.Context, node *specs.Node, flow specs.FlowInterface, stack Stack, prop *specs.Property, functions Custom) error {
	if prop == nil {
		return nil
	}

	if prop.Message != nil {
		for _, nested := range prop.Message {
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
func PrepareFunction(ctx *broker.Context, node *specs.Node, flow specs.FlowInterface, property *specs.Property, stack Stack, methods Custom) error {
	if property == nil {
		return nil
	}

	if !FunctionPattern.MatchString(property.Raw) {
		return nil
	}

	pattern := FunctionPattern.FindStringSubmatch(property.Raw)
	fn := pattern[1]
	args := strings.Split(pattern[2], ArgumentDelimiter)

	if methods[fn] == nil {
		return ErrUndefinedFunction{
			Function: fn,
			Property: property.Raw,
		}
	}

	arguments := make([]*specs.Property, len(args))

	for index, arg := range args {
		result, err := template.ParseContent(property.Path, property.Name, strings.TrimSpace(arg))
		if err != nil {
			return err
		}

		err = references.ResolveProperty(ctx, node, result, flow)
		if err != nil {
			return err
		}

		err = PrepareFunction(ctx, node, flow, result, stack, methods)
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
	function := &Function{
		Arguments: arguments,
		Fn:        handle,
		Returns:   property,
	}

	stack[ref] = function

	property.Template = returns.Template
	property.Label = returns.Label
	property.Reference = &specs.PropertyReference{
		Resource: template.JoinPath(template.StackResource, ref),
		Path:     ".",
		Property: returns,
	}

	references.ScopeNestedReferences(&returns.Template, &property.Template)
	return nil
}

// GenerateStackReference generates a unique path prefix which could be used to isolate functions
func GenerateStackReference() string {
	bb := make([]byte, 5)
	rand.Read(bb)
	return hex.EncodeToString(bb)
}
