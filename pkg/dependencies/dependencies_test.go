package dependencies

import (
	"testing"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"github.com/jexia/semaphore/v2/pkg/specs/template"
)

func TestResolveDependencies(t *testing.T) {
	flows := specs.FlowListInterface{
		&specs.Flow{
			Name: "first",
			Nodes: specs.NodeList{
				{
					ID: "first",
				},
				{
					ID: "second",
					DependsOn: map[string]*specs.Node{
						"first": nil,
					},
				},
				{
					ID: "third",
					DependsOn: map[string]*specs.Node{
						"first":  nil,
						"second": nil,
					},
				},
				{
					ID: "fourth",
					Intermediate: &specs.ParameterMap{
						DependsOn: map[string]*specs.Node{
							"first":  nil,
							"second": nil,
						},
					},
				},
				{
					ID: "fifth",
					Call: &specs.Call{
						Request: &specs.ParameterMap{
							DependsOn: map[string]*specs.Node{
								"first":  nil,
								"second": nil,
							},
						},
					},
				},
			},
		},
		&specs.Flow{
			Name: "second",
			Nodes: specs.NodeList{
				{
					ID: "first",
				},
				{
					ID: "second",
					DependsOn: map[string]*specs.Node{
						"first": nil,
					},
				},
				{
					ID: "third",
					DependsOn: map[string]*specs.Node{
						"first":  nil,
						"second": nil,
					},
				},
			},
		},
		&specs.Flow{
			Name: "third",
			Nodes: specs.NodeList{
				{
					ID: "first",
				},
				{
					ID: "second",
					DependsOn: map[string]*specs.Node{
						"first": nil,
					},
				},
				{
					ID: "third",
					DependsOn: map[string]*specs.Node{
						"first":  nil,
						"second": nil,
					},
				},
			},
			Output: &specs.ParameterMap{
				DependsOn: map[string]*specs.Node{
					"third": nil,
				},
			},
		},
		&specs.Proxy{
			Name: "first",
			Nodes: specs.NodeList{
				{
					ID: "first",
				},
				{
					ID: "second",
					DependsOn: map[string]*specs.Node{
						"first": nil,
					},
				},
				{
					ID: "third",
					DependsOn: map[string]*specs.Node{
						"first":  nil,
						"second": nil,
					},
				},
			},
		},
		&specs.Proxy{
			Name: "second",
			Nodes: specs.NodeList{
				{
					ID: "first",
				},
				{
					ID: "second",
					DependsOn: map[string]*specs.Node{
						"first": nil,
					},
				},
				{
					ID: "third",
					DependsOn: map[string]*specs.Node{
						"first":  nil,
						"second": nil,
					},
				},
			},
		},
		&specs.Proxy{
			Name: "third",
			Nodes: specs.NodeList{
				{
					ID: "first",
				},
				{
					ID: "second",
					DependsOn: map[string]*specs.Node{
						"first": nil,
					},
				},
				{
					ID: "third",
					DependsOn: map[string]*specs.Node{
						"first":  nil,
						"second": nil,
					},
				},
			},
		},
	}

	ctx := logger.WithLogger(broker.NewBackground())
	err := ResolveFlows(ctx, flows)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	for _, flow := range flows {
		for _, call := range flow.GetNodes() {
			if len(call.DependsOn) == 0 {
				continue
			}

			for key, val := range call.DependsOn {
				if val == nil {
					t.Fatalf("call dependency not resolved %s.%s", call.ID, key)
				}
			}
		}
	}
}

func TestResolveCallDependencies(t *testing.T) {
	flow := &specs.Flow{
		Nodes: specs.NodeList{
			{
				ID: "first",
			},
			{
				ID: "second",
				DependsOn: map[string]*specs.Node{
					"first": nil,
				},
			},
			{
				ID: "third",
				DependsOn: map[string]*specs.Node{
					"first":  nil,
					"second": nil,
				},
			},
		},
	}

	tests := specs.NodeList{
		flow.Nodes[1],
		flow.Nodes[2],
	}

	for _, input := range tests {
		err := Resolve(flow, input.DependsOn, input.ID, make(Unresolved))
		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}

		for key, val := range input.DependsOn {
			if val == nil {
				t.Fatalf("dependency not resolved %s.%s", input.ID, key)
			}
		}
	}
}

func TestCallCircularDependenciesDetection(t *testing.T) {
	flow := &specs.Flow{
		Nodes: specs.NodeList{
			{
				ID: "first",
				DependsOn: map[string]*specs.Node{
					"second": nil,
				},
			},
			{
				ID: "second",
				DependsOn: map[string]*specs.Node{
					"first": nil,
				},
			},
		},
	}

	tests := map[string]*specs.Node{
		"first":  flow.Nodes[0],
		"second": flow.Nodes[1],
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			err := Resolve(flow, input.DependsOn, input.ID, make(Unresolved))
			if err == nil {
				t.Fatalf("unexpected pass %s", input.ID)
			}
		})
	}
}

func TestFlowCircularDependenciesDetection(t *testing.T) {
	flows := map[string]specs.FlowListInterface{
		"simple": {
			&specs.Flow{
				Nodes: specs.NodeList{
					{
						ID: "first",
						DependsOn: map[string]*specs.Node{
							"second": nil,
						},
					},
					{
						ID: "second",
						DependsOn: map[string]*specs.Node{
							"first": nil,
						},
					},
				},
			},
		},
		"intermediate": {
			&specs.Flow{
				Nodes: specs.NodeList{
					{
						ID: "first",
						Intermediate: &specs.ParameterMap{
							DependsOn: map[string]*specs.Node{
								"second": nil,
							},
						},
					},
					{
						ID: "second",
						Intermediate: &specs.ParameterMap{
							DependsOn: map[string]*specs.Node{
								"first": nil,
							},
						},
					},
				},
			},
		},
		"output": {
			&specs.Flow{
				Nodes: specs.NodeList{
					{
						ID: "first",
						Intermediate: &specs.ParameterMap{
							DependsOn: map[string]*specs.Node{
								template.OutputResource: nil,
							},
						},
					},
				},
				Output: &specs.ParameterMap{
					DependsOn: map[string]*specs.Node{
						"first": nil,
					},
				},
			},
		},
		"call": {
			&specs.Flow{
				Nodes: specs.NodeList{
					{
						ID: "first",
						DependsOn: map[string]*specs.Node{
							"second": nil,
						},
					},
					{
						ID: "second",
						Call: &specs.Call{
							Request: &specs.ParameterMap{
								DependsOn: map[string]*specs.Node{
									"first": nil,
								},
							},
						},
					},
				},
			},
		},
	}

	for name, input := range flows {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			err := ResolveFlows(ctx, input)
			if err == nil {
				t.Fatalf("unexpected pass")
			}
		})
	}
}

func TestSelfDependencyDetection(t *testing.T) {
	flow := &specs.Flow{
		Nodes: specs.NodeList{
			{
				ID: "first",
				DependsOn: map[string]*specs.Node{
					"first": nil,
				},
			},
			{
				ID: "second",
				DependsOn: map[string]*specs.Node{
					"second": nil,
				},
			},
		},
	}

	tests := specs.NodeList{
		flow.Nodes[0],
	}

	for _, input := range tests {
		err := Resolve(flow, input.DependsOn, input.ID, make(Unresolved))
		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}

		if len(input.DependsOn) > 0 {
			t.Fatalf("unexpted remaining dependencies, expected dependencies to be empty: %+v", input.DependsOn)
		}
	}
}
