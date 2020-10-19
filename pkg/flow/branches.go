package flow

import (
	"github.com/jexia/semaphore/pkg/lookup"
)

// ConstructBranches constructs the node branches based on the made references
func ConstructBranches(nodes []*Node) {
	for _, node := range nodes {
		for _, dependency := range node.DependsOn {
			if dependency == nil {
				continue
			}

			ConstructDependency(node, dependency.ID, nodes)
		}

		for _, reference := range node.References {
			target, _ := lookup.ParseResource(reference.Resource)
			ConstructDependency(node, target, nodes)
		}
	}
}

// ConstructDependency constructs a dependency for the given node
func ConstructDependency(node *Node, target string, nodes []*Node) {
	for _, parent := range nodes {
		if parent.Name == target {
			if !node.Previous.Has(parent.Name) {
				node.Previous = append(node.Previous, parent)
			}

			if !parent.Next.Has(node.Name) {
				parent.Next = append(parent.Next, node)
			}
		}
	}
}

// FetchStarting constructs the starting seeds for the given nodes
func FetchStarting(nodes []*Node) (result []*Node) {
	for _, node := range nodes {
		if len(node.Previous) == 0 {
			result = append(result, node)
		}
	}

	return result
}
