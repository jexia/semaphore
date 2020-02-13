package specs

import (
	"fmt"
)

// ResolveManifestDependencies resolves all dependencies inside the given manifest
func ResolveManifestDependencies(manifest *Manifest) error {
	for _, flow := range manifest.Flows {
		err := ResolveFlowDependencies(manifest, flow, make(map[string]*Flow))
		if err != nil {
			return err
		}

		for _, call := range flow.Calls {
			err := ResolveCallDependencies(flow, call, make(map[string]*Call))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ResolveFlowDependencies resolves the given flow dependencies and attempts to detect any circular dependencies
func ResolveFlowDependencies(manifest *Manifest, node *Flow, unresolved map[string]*Flow) error {
	unresolved[node.Name] = node

lookup:
	for edge := range node.DependsOn {
		_, unresolv := unresolved[edge]
		if unresolv {
			return fmt.Errorf("Circular flow dependency detected: %s <-> %s", node.Name, edge)
		}

		for _, flow := range manifest.Flows {
			if flow.Name == edge {
				err := ResolveFlowDependencies(manifest, flow, unresolved)
				if err != nil {
					return err
				}

				node.DependsOn[edge] = flow
				continue lookup
			}
		}
	}

	delete(unresolved, node.Name)
	return nil
}

// ResolveCallDependencies resolves the given call dependencies and attempts to detect any circular dependencies
func ResolveCallDependencies(flow *Flow, node *Call, unresolved map[string]*Call) error {
	unresolved[node.Name] = node

lookup:
	for edge := range node.DependsOn {
		_, unresolv := unresolved[edge]
		if unresolv {
			return fmt.Errorf("Circular call dependency detected: %s.%s <-> %s.%s", flow.Name, node.Name, flow.Name, edge)
		}

		for _, call := range flow.Calls {
			if call.Name == edge {
				err := ResolveCallDependencies(flow, call, unresolved)
				if err != nil {
					return err
				}

				node.DependsOn[edge] = call
				continue lookup
			}
		}
	}

	delete(unresolved, node.Name)
	return nil
}
