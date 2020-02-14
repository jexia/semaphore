package specs

import (
	"testing"

	"github.com/jexia/maestro/utils"
)

func TestDuplicateManifests(t *testing.T) {
	tests := []*Manifest{
		{
			Services: []*Service{
				{
					Alias: "dup",
				},
				{
					Alias: "dup",
				},
			},
		},
		{
			Callers: []*Caller{
				{
					Name: "dup",
				},
				{
					Name: "dup",
				},
			},
		},
		{
			Endpoints: []*Endpoint{
				{
					Flow: "dup",
				},
				{
					Flow: "dup",
				},
			},
		},
		{
			Flows: []*Flow{
				{
					Name: "dup",
				},
				{
					Name: "dup",
				},
			},
		},
		{
			Flows: []*Flow{
				{
					Name: "first",
					Calls: []*Call{
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

	for _, input := range tests {
		err := CheckManifestDuplicates(input)
		if err == nil {
			t.Fatal("unexpected pass", input)
		}
	}
}

func TestDuplicateFlow(t *testing.T) {
	manifest := &Manifest{
		File: utils.FileInfo{
			Path: "test",
		},
	}

	tests := []*Flow{
		{
			Name: "first",
			Calls: []*Call{
				{
					Name: "dup",
				},
				{
					Name: "dup",
				},
			},
		},
	}

	for _, input := range tests {
		err := CheckFlowDuplicates(manifest, input)
		if err == nil {
			t.Fatal("unexpected pass", input)
		}
	}
}
