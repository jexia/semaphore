package openapi3

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
)

func Test_scanPaths(t *testing.T) {
	type testEndpoint struct {
		path, method string
		objects      []string
	}

	loader := openapi3.NewSwaggerLoader()

	doc, err := loader.LoadSwaggerFromFile("./fixtures/petstore.yml")
	if err != nil {
		t.Fatalf("unexpected error on loading openapi document: %v", err)
	}

	endpoints, err := scanPaths(doc.Paths)
	assert.Nil(t, err)
	assert.Len(t, endpoints, 3, "should have 3 registered paths")

	getPets := findEndpoint(endpoints, "/pets", "GET")
	assert.NotNil(t, getPets)
	assert.NotEmpty(t, getPets.objects()["GET:/pets:Response[application/json][200]"])
	assert.NotEmpty(t, getPets.objects()["GET:/pets:Response[application/json][default]"])

	createPet := findEndpoint(endpoints, "/pets", "POST")
	assert.NotNil(t, createPet)
	assert.NotEmpty(t, createPet.objects()["POST:/pets:Response[application/json][default]"])
	assert.NotEmpty(t, createPet.objects()["POST:/pets:Request[application/json]"])

	getPet := findEndpoint(endpoints, "/pets/{petId}", "GET")
	assert.NotNil(t, getPet)
	assert.NotEmpty(t, getPet.objects()["GET:/pets/{petId}:Response[application/json][200]"])
	assert.NotEmpty(t, getPet.objects()["GET:/pets/{petId}:Response[application/json][default]"])
}

func findEndpoint(endpoints []*endpointRef, path, method string) *endpointRef {
	for _, e := range endpoints {
		if e.path == path && e.method == method {
			return e
		}
	}

	return nil
}
