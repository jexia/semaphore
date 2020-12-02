package json

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/specs"
	"go.uber.org/zap"
)

type collection struct {
	Flows     specs.FlowList     `json:"flows,omitempty"`
	Proxy     specs.ProxyList    `json:"proxy,omitempty"`
	Endpoints specs.EndpointList `json:"endpoints,omitempty"`
	Services  specs.ServiceList  `json:"services,omitempty"`
	Schemas   specs.Schemas      `json:"schemas,omitempty"`
}

func (collection *collection) Append(incoming collection) {
	collection.Flows.Append(incoming.Flows)
	collection.Proxy.Append(incoming.Proxy)
	collection.Endpoints.Append(incoming.Endpoints)
	collection.Services.Append(incoming.Services)
	collection.Schemas.Append(incoming.Schemas)
}

// ServicesResolver constructs a JSON service resolver for the given path.
func ServicesResolver(path string) providers.ServicesResolver {
	return func(ctx *broker.Context) (specs.ServiceList, error) {
		logger.Debug(ctx, "resolving JSON services", zap.String("path", path))

		collection, err := resolvePath(ctx, path)
		if err != nil {
			return nil, err
		}

		return collection.Services, nil
	}
}

// FlowsResolver constructs a resource resolver for the given path
func FlowsResolver(path string) providers.FlowsResolver {
	return func(ctx *broker.Context) (specs.FlowListInterface, error) {
		logger.Debug(ctx, "resolving JSON flows", zap.String("path", path))

		collection, err := resolvePath(ctx, path)
		if err != nil {
			return nil, err
		}

		result := make(specs.FlowListInterface, 0, len(collection.Flows)+len(collection.Proxy))

		for _, flow := range collection.Flows {
			result = append(result, flow)
		}

		for _, flow := range collection.Proxy {
			result = append(result, flow)
		}

		return result, nil
	}
}

// EndpointsResolver constructs a resource resolver for the given path
func EndpointsResolver(path string) providers.EndpointsResolver {
	return func(ctx *broker.Context) (specs.EndpointList, error) {
		logger.Debug(ctx, "resolving JSON endpoints", zap.String("path", path))

		collection, err := resolvePath(ctx, path)
		if err != nil {
			return nil, err
		}

		return collection.Endpoints, nil
	}
}

// SchemaResolver constructs a schema resolver for the given path
func SchemaResolver(path string) providers.SchemaResolver {
	return func(ctx *broker.Context) (specs.Schemas, error) {
		logger.Debug(ctx, "resolving JSON schemas", zap.String("path", path))

		collection, err := resolvePath(ctx, path)
		if err != nil {
			return nil, err
		}

		return collection.Schemas, nil
	}
}

// resolvePath resolves the given path and returns the available manifests.
// All defined includes are followed and their manifests are included
func resolvePath(ctx *broker.Context, path string) (collection, error) {
	logger.Debug(ctx, "resolving JSON path", zap.String("path", path))

	files, err := providers.ResolvePath(ctx, []string{}, path)
	if err != nil {
		return collection{}, err
	}

	if len(files) == 0 {
		return collection{}, ErrPathNotFound{
			Path: path,
		}
	}

	logger.Debug(ctx, "files found", zap.String("path", path), zap.Int("files", len(files)))
	result := collection{}

	for _, file := range files {
		logger.Debug(ctx, "resolving file", zap.String("path", file.Path))

		reader, err := os.Open(file.Path)
		if err != nil {
			return collection{}, err
		}

		bb, err := ioutil.ReadAll(reader)
		if err != nil {
			return collection{}, err
		}

		collection := collection{}
		json.Unmarshal(bb, &collection)

		result.Append(collection)
	}

	logger.Debug(ctx, "resolve path result", zap.String("path", path), zap.Int("manifests", len(files)))
	return result, nil
}
