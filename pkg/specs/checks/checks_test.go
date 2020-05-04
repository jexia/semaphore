package checks

import (
	"testing"

	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/specs"
)

func TestDuplicateManifests(t *testing.T) {
	tests := map[string]*specs.FlowsManifest{
		"duplicate flow": {

			Flows: []*specs.Flow{
				{
					Name: "dup",
				},
				{
					Name: "dup",
				},
			},
		},
		"duplicate proxy": {

			Proxy: []*specs.Proxy{
				{
					Name: "dup",
				},
				{
					Name: "dup",
				},
			},
		},
		"duplicate node": {

			Flows: []*specs.Flow{
				{
					Name: "first",
					Nodes: []*specs.Node{
						{
							Name: "dup",
						},
						{
							Name: "dup",
						},
					},
				},
			},
		},
		"duplicate proxy node": {

			Proxy: []*specs.Proxy{
				{
					Name: "first",
					Nodes: []*specs.Node{
						{
							Name: "dup",
						},
						{
							Name: "dup",
						},
					},
				},
			},
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
			err := ManifestDuplicates(ctx, input)
			if err == nil {
				t.Fatal("unexpected pass", input)
			}
		})
	}
}

func TestDuplicateNodes(t *testing.T) {
	tests := map[string]*specs.Flow{
		"simple": {
			Name: "first",
			Nodes: []*specs.Node{
				{
					Name: "dup",
				},
				{
					Name: "dup",
				},
			},
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
			err := NodeDuplicates(ctx, input.Name, input.Nodes)
			if err == nil {
				t.Fatal("unexpected pass", input)
			}
		})
	}
}
