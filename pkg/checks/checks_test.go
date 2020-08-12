package checks

import (
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/specs"
)

func TestDuplicateManifests(t *testing.T) {
	tests := map[string]specs.FlowListInterface{
		"duplicate flow": {
			&specs.Flow{
				Name: "dup",
			},
			&specs.Flow{
				Name: "dup",
			},
		},
		"duplicate proxy": {
			&specs.Proxy{
				Name: "dup",
			},
			&specs.Proxy{
				Name: "dup",
			},
		},
		"duplicate node": {
			&specs.Flow{
				Name: "first",
				Nodes: specs.NodeList{
					{
						ID: "dup",
					},
					{
						ID: "dup",
					},
				},
			},
		},
		"duplicate proxy node": {
			&specs.Proxy{
				Name: "first",
				Nodes: specs.NodeList{
					{
						ID: "dup",
					},
					{
						ID: "dup",
					},
				},
			},
		},
		"duplicate flow - proxy": {
			&specs.Proxy{
				Name: "first",
			},
			&specs.Flow{
				Name: "first",
			},
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewContext())
			err := FlowDuplicates(ctx, input)
			if err == nil {
				t.Fatal("unexpected pass", input)
			}
		})
	}
}

func TestDuplicateManifestsPass(t *testing.T) {
	tests := map[string]specs.FlowListInterface{
		"multiple flow": {
			&specs.Flow{
				Name: "first",
			},
			&specs.Flow{
				Name: "second",
			},
		},
		"flow and proxy": {
			&specs.Flow{
				Name: "first",
			},
			&specs.Proxy{
				Name: "second",
			},
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewContext())
			err := FlowDuplicates(ctx, input)
			if err != nil {
				t.Fatal("unexpected fail", err)
			}
		})
	}
}

func TestReservedKeywordsManifests(t *testing.T) {
	tests := map[string]specs.FlowListInterface{
		"error": {
			&specs.Flow{
				Name: "first",
				Nodes: specs.NodeList{
					{
						ID: "error",
					},
				},
			},
		},
		"input": {
			&specs.Flow{
				Name: "first",
				Nodes: specs.NodeList{
					{
						ID: "input",
					},
				},
			},
		},
		"stack": {
			&specs.Flow{
				Name: "first",
				Nodes: specs.NodeList{
					{
						ID: "stack",
					},
				},
			},
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewContext())
			err := FlowDuplicates(ctx, input)
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
			Nodes: specs.NodeList{
				{
					ID: "dup",
				},
				{
					ID: "dup",
				},
			},
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewContext())
			err := NodeDuplicates(ctx, input.Name, input.Nodes)
			if err == nil {
				t.Fatal("unexpected pass", input)
			}
		})
	}
}
