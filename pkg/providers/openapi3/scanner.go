package openapi3

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
)

type endpointRef struct {
	path   string
	method string

	responses map[string]*objectRef
	requests  map[string]*objectRef
}

func newEndpointRef(method, path string) *endpointRef {
	return &endpointRef{
		method: method,
		path: path,
		responses: map[string]*objectRef{},
		requests: map[string]*objectRef{},
	}
}

// return all the objects
func (e *endpointRef) objects() map[string]*objectRef {
	var collected = make(map[string]*objectRef, len(e.responses) + len(e.requests))

	for name, obj := range e.responses {
		collected[name] = obj
	}

	for name, obj := range e.requests {
		collected[name] = obj
	}

	return collected
}

func (e *endpointRef) addRequest(obj *objectRef) error {
	var name string

	if obj.modelName != "" {
		name = obj.modelName
	} else {
		name = fmt.Sprintf("%s:%s:Request[%s]", e.method, e.path, obj.contentType)
	}

	e.requests[name] = obj
	return nil
}

func (e *endpointRef) addResponse(obj *objectRef) error {
	var name string

	if obj.modelName != "" {
		name = obj.modelName
	} else {
		name =  fmt.Sprintf("%s:%s:Response[%s][%s]", e.method, e.path, obj.contentType, obj.code)
	}

	e.responses[name] = obj
	return nil
}

type objectRef struct {
	code        string // is used by response object only
	contentType string
	ref         *openapi3.SchemaRef
	modelName   string
}

// scan paths object and extract every request body and response,
// naming them with a unique generated name based on the endpoint,
// verb, code and content type.
func scanPaths(paths openapi3.Paths) ([]*endpointRef, error) {
	var endpoints []*endpointRef

	for path, item := range paths {
		for method, op := range item.Operations() {
			e := newEndpointRef(method, path)

			for _, response := range scanResponses(op.Responses) {
				if err := e.addResponse(response); err != nil {
					return nil, fmt.Errorf("failed to process response schema %s:%s: %w", method, path, err)
				}
			}

			for _, request := range scanRequests(op.RequestBody) {
				if err := e.addRequest(request); err != nil {
					return nil, fmt.Errorf("failed to process request body %s:%s: %w", method, path, err)
				}
			}

			endpoints = append(endpoints, e)
		}
	}

	return endpoints, nil
}

func scanResponses(responses openapi3.Responses) []*objectRef {
	var collected []*objectRef

	for code, ref := range responses {
		if ref.Value == nil {
			continue
		}

		for contentType, media := range ref.Value.Content {
			if media == nil || media.Schema == nil {
				continue
			}

			r := &objectRef{
				contentType: contentType,
				ref:         media.Schema,
				code:        code,
			}

			collected = append(collected, r)
		}
	}

	return collected
}

func scanRequests(ref *openapi3.RequestBodyRef) []*objectRef {
	var collected []*objectRef

	if ref == nil || ref.Value == nil {
		return nil
	}

	for contentType, media := range ref.Value.Content {
		if media == nil || media.Schema == nil {
			continue
		}

		r := &objectRef{
			contentType: contentType,
			ref:         media.Schema,
		}

		collected = append(collected, r)
	}

	return collected
}
