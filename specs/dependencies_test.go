package specs

import (
	"testing"
)

func TestResolveManifestDependencies(t *testing.T) {
	manifest := &Manifest{
		Flows: []*Flow{
			{
				Name: "first",
				Calls: []*Call{
					{
						Name: "first",
					},
					{
						Name: "second",
						DependsOn: map[string]*Call{
							"first": nil,
						},
					},
					{
						Name: "third",
						DependsOn: map[string]*Call{
							"first":  nil,
							"second": nil,
						},
					},
				},
			},
			{
				Name: "second",
				DependsOn: map[string]*Flow{
					"first": nil,
				},
				Calls: []*Call{
					{
						Name: "first",
					},
					{
						Name: "second",
						DependsOn: map[string]*Call{
							"first": nil,
						},
					},
					{
						Name: "third",
						DependsOn: map[string]*Call{
							"first":  nil,
							"second": nil,
						},
					},
				},
			},
			{
				Name: "third",
				DependsOn: map[string]*Flow{
					"first":  nil,
					"second": nil,
				},
				Calls: []*Call{
					{
						Name: "first",
					},
					{
						Name: "second",
						DependsOn: map[string]*Call{
							"first": nil,
						},
					},
					{
						Name: "third",
						DependsOn: map[string]*Call{
							"first":  nil,
							"second": nil,
						},
					},
				},
			},
		},
	}

	err := ResolveManifestDependencies(manifest)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	for _, flow := range manifest.Flows {
		if len(flow.DependsOn) == 0 {
			continue
		}

		for key, val := range flow.DependsOn {
			if val == nil {
				t.Fatalf("flow dependency not resolved %s.%s", flow.Name, key)
			}
		}

		for _, call := range flow.Calls {
			if len(call.DependsOn) == 0 {
				continue
			}

			for key, val := range call.DependsOn {
				if val == nil {
					t.Fatalf("call dependency not resolved %s.%s", call.Name, key)
				}
			}
		}
	}
}

func TestResolveFlowDependencies(t *testing.T) {
	manifest := &Manifest{
		Flows: []*Flow{
			{
				Name: "first",
			},
			{
				Name: "second",
				DependsOn: map[string]*Flow{
					"first": nil,
				},
			},
			{
				Name: "third",
				DependsOn: map[string]*Flow{
					"first":  nil,
					"second": nil,
				},
			},
		},
	}

	tests := []*Flow{
		manifest.Flows[1],
		manifest.Flows[2],
	}

	for _, input := range tests {
		err := ResolveFlowManagerDependencies(manifest, input, make(map[string]FlowManager))
		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}

		for key, val := range input.DependsOn {
			if val == nil {
				t.Fatalf("dependency not resolved %s.%s", input.Name, key)
			}
		}
	}
}

func TestFlowCircularDependenciesDetection(t *testing.T) {
	manifest := &Manifest{
		Flows: []*Flow{
			&Flow{
				Name: "first",
				DependsOn: map[string]*Flow{
					"second": nil,
				},
			},
			&Flow{
				Name: "second",
				DependsOn: map[string]*Flow{
					"first": nil,
				},
			},
		},
	}

	tests := []*Flow{
		manifest.Flows[0],
		manifest.Flows[1],
	}

	for _, input := range tests {
		err := ResolveFlowManagerDependencies(manifest, input, make(map[string]FlowManager))
		if err == nil {
			t.Fatalf("unexpected pass %s", input.Name)
		}
	}
}

func TestResolveCallDependencies(t *testing.T) {
	flow := &Flow{
		Calls: []*Call{
			{
				Name: "first",
			},
			{
				Name: "second",
				DependsOn: map[string]*Call{
					"first": nil,
				},
			},
			{
				Name: "third",
				DependsOn: map[string]*Call{
					"first":  nil,
					"second": nil,
				},
			},
		},
	}

	tests := []*Call{
		flow.Calls[1],
		flow.Calls[2],
	}

	for _, input := range tests {
		err := ResolveCallDependencies(flow, input, make(map[string]*Call))
		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}

		for key, val := range input.DependsOn {
			if val == nil {
				t.Fatalf("dependency not resolved %s.%s", input.Name, key)
			}
		}
	}
}

func TestCallCircularDependenciesDetection(t *testing.T) {
	flow := &Flow{
		Calls: []*Call{
			{
				Name: "first",
				DependsOn: map[string]*Call{
					"second": nil,
				},
			},
			{
				Name: "second",
				DependsOn: map[string]*Call{
					"first": nil,
				},
			},
		},
	}

	tests := []*Call{
		flow.Calls[0],
		flow.Calls[1],
	}

	for _, input := range tests {
		err := ResolveCallDependencies(flow, input, make(map[string]*Call))
		if err == nil {
			t.Fatalf("unexpected pass %s", input.Name)
		}
	}
}
