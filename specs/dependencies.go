package specs

import (
	"fmt"
)

// ResolveManifestDependencies resolves all dependencies inside the given manifest
func ResolveManifestDependencies(manifest *Manifest) error {
	for _, flow := range manifest.Flows {
		err := ResolveFlowManagerDependencies(manifest, flow, make(map[string]FlowManager))
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

	for _, proxy := range manifest.Proxy {
		err := ResolveFlowManagerDependencies(manifest, proxy, make(map[string]FlowManager))
		if err != nil {
			return err
		}

		for _, call := range proxy.Calls {
			err := ResolveCallDependencies(proxy, call, make(map[string]*Call))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ResolveCallDependencies resolves the given call dependencies and attempts to detect any circular dependencies
func ResolveCallDependencies(manager FlowManager, node *Call, unresolved map[string]*Call) error {
	unresolved[node.Name] = node

lookup:
	for edge := range node.DependsOn {
		_, unresolv := unresolved[edge]
		if unresolv {
			return fmt.Errorf("Circular dependency detected: %s.%s <-> %s.%s", manager.GetName(), node.Name, manager.GetName(), edge)
		}

		for _, call := range manager.GetCalls() {
			if call.Name == edge {
				err := ResolveCallDependencies(manager, call, unresolved)
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

// ResolveFlowManagerDependencies resolves the given flow dependencies and attempts to detect any circular dependencies
func ResolveFlowManagerDependencies(manifest *Manifest, node FlowManager, unresolved map[string]FlowManager) error {
	unresolved[node.GetName()] = node

lookup:
	for edge := range node.GetDependencies() {
		_, unresolv := unresolved[edge]
		if unresolv {
			return fmt.Errorf("Circular dependency detected: %s <-> %s", node.GetName(), edge)
		}

		for _, flow := range manifest.Flows {
			if flow.Name == edge {
				err := ResolveFlowManagerDependencies(manifest, flow, unresolved)
				if err != nil {
					return err
				}

				node.GetDependencies()[edge] = flow
				continue lookup
			}
		}
	}

	delete(unresolved, node.GetName())
	return nil
}
