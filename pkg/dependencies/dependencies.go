package dependencies

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
)

// Unresolved represents a collection of unresolved references
type Unresolved map[string]struct{}

// ResolveFlows resolves all dependencies inside the given manifest
func ResolveFlows(ctx *broker.Context, flows specs.FlowListInterface) error {
	logger.Info(ctx, "resolving flow dependencies")

	for _, flow := range flows {
		for _, node := range flow.GetNodes() {
			call := node.Call
			if node.Call != nil && call.Request != nil {
				err := Resolve(flow, call.Request.DependsOn, node.ID, make(Unresolved))
				if err != nil {
					return err
				}

				node.DependsOn = node.DependsOn.Append(call.Request.DependsOn)
			}

			intermediate := node.Intermediate
			if intermediate != nil {
				err := Resolve(flow, intermediate.DependsOn, node.ID, make(Unresolved))
				if err != nil {
					return err
				}

				node.DependsOn = node.DependsOn.Append(intermediate.DependsOn)
			}

			err := Resolve(flow, node.DependsOn, node.ID, make(Unresolved))
			if err != nil {
				return err
			}
		}

		output := flow.GetOutput()
		if output != nil {
			err := Resolve(flow, output.DependsOn, template.OutputResource, make(Unresolved))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Resolve resolves the given call dependencies and attempts to detect any circular dependencies
func Resolve(manager specs.FlowInterface, dependencies specs.Dependencies, id string, unresolved Unresolved) error {
	if len(dependencies) == 0 {
		return nil
	}

	unresolved[id] = struct{}{}

	for edge := range dependencies {
		// Remove any self references
		if edge == id {
			delete(unresolved, edge)
			delete(dependencies, edge)
			continue
		}

		if _, unresolv := unresolved[edge]; unresolv {
			return ErrCircularDependency{
				Flow: manager.GetName(),
				From: id,
				To:   edge,
			}
		}

		result := manager.GetNodes().Get(edge)
		if result == nil {
			continue
		}

		err := Resolve(manager, result.DependsOn, result.ID, unresolved)
		if err != nil {
			return err
		}

		dependencies[edge] = result
	}

	delete(unresolved, id)

	return nil
}
