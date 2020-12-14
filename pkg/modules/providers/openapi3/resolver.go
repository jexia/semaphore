package openapi3

import (
	"encoding/json"
	"fmt"

	openapi "github.com/getkin/kin-openapi/openapi3"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/specs"
	"go.uber.org/zap"
)

const (
	// XPackageExtensionField is the name of an info property which is used to define package name.
	// Example:
	//
	// info:
	//   x-semaphore-package: com.semaphore
	// components:
	//   <the rest of file>
	XPackageExtensionField = "x-semaphore-package"

	// XModelName is the name of an optional extension property included into response and request objects,
	// and defines a custom name for the object.
	XModelName = "x-semaphore-model"
)

// a dictionary of file names to swagger documents.
type swaggers map[string]*openapi.Swagger

// collect all the swagger documents from the paths.
// imports is a collection of all the files. The file path might include a mask.
// Example: []string{"/etc/schemas/user.yml", "/etc/schemas/animal_*.yml"} and so on.
func collect(ctx *broker.Context, imports []string) (swaggers, error) {
	var (
		docs = swaggers{} // used to collect all the loaded & parsed swagger files
	)

	loader := openapi.NewSwaggerLoader()

	for _, path := range imports {
		files, err := providers.ResolvePath(ctx, []string{}, path)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve path %s: %w", path, err)
		}

		// iterate over all the files matched by the single import path
		for _, file := range files {
			doc, err := loader.LoadSwaggerFromFile(file.Path)

			if err != nil {
				return nil, fmt.Errorf("failed to parse openapi file %s: %w", file.Path, err)
			}

			docs[file.Path] = doc
		}
	}

	return docs, nil
}

// return canonical name based on the package name (optional extra field info.x-semaphore-package) and the given name
//
// Example:
//
// info:
//   x-semaphore-package: com.semaphore
// paths:
//   ...
//
// the canonical name for User is `com.semaphore.User`.
// If the package is not defined, the name is `User`.
func getCanonicalName(doc *openapi.Swagger, name string) string {
	if doc.Info == nil {
		return name
	}

	prop := doc.Info.Extensions[XPackageExtensionField]
	if prop == nil {
		return name
	}

	raw, ok := prop.(json.RawMessage)
	if !ok {
		return name
	}

	var pkg string
	if err := json.Unmarshal(raw, &pkg); err != nil {
		return name
	}

	return fmt.Sprintf("%s.%s", pkg, name)
}

// SchemaResolver returns a new schemaRef resolver for the given openapi collection
func SchemaResolver(paths []string) providers.SchemaResolver {
	return func(ctx *broker.Context) (specs.Schemas, error) {
		logger.Debug(ctx, "resolving openapi schemas", zap.Strings("paths", paths))

		docs, err := collect(ctx, paths)
		if err != nil {
			return nil, err
		}

		schemas, err := newSchemas(docs)

		if err != nil {
			return nil, fmt.Errorf("failed to build schema: %w", err)
		}

		return schemas, nil
	}
}
