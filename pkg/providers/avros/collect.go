package avros

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/providers"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"go.uber.org/zap"
)

// Collect attempts to collect all the available avro files inside the given path and parses them to resources
func Collect(ctx *broker.Context, paths []string, path string) ([]*AvroSchema, error) {
	imports := make([]string, len(paths))
	for index, path := range paths {
		imports[index] = path
	}

	logger.Debug(ctx, "collect available avro", zap.String("path", path))
	logger.Debug(ctx, "collect available avro with imports", zap.Strings("imports", paths))

	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	logger.Debug(ctx, "absolute avro path", zap.String("path", path))

	for index, path := range imports {
		path, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}

		imports[index] = path
	}

	logger.Debug(ctx, "absolute avro imports", zap.Strings("imports", imports))

	files, err := providers.ResolvePath(ctx, []string{}, path)
	if err != nil {
		return nil, err
	}

	for index, path := range imports {
		stat, err := os.Stat(path)
		if err != nil {
			imports[index] = filepath.Dir(path)
			continue
		}

		if stat.IsDir() {
			imports[index] = path
			continue
		}

		imports[index] = filepath.Dir(path)
	}

	descriptors, err := UnmarshalFiles(imports, files)
	if err != nil {
		return nil, err
	}

	return descriptors, nil
}

// SchemaResolver returns a new schema resolver for the given avro collection
func SchemaResolver(imports []string, path string) providers.SchemaResolver {
	return func(ctx *broker.Context) (specs.Schemas, error) {
		logger.Debug(ctx, "resolving acro schemas", zap.String("path", path))

		files, err := Collect(ctx, imports, path)
		if err != nil {
			return nil, err
		}

		return NewSchema(files), nil
	}
}

// UnmarshalFiles attempts to parse the given HCL files to intermediate resources.
// Files are parsed based from the given import paths
func UnmarshalFiles(imports []string, files []*providers.FileInfo) ([]*AvroSchema, error) {
	results := make([]*AvroSchema, 0)
	for _, file := range files {
		schema, err := ioutil.ReadFile(file.Path)
		if err != nil {
			return nil, err
		}
		tempSchema := AvroSchema{}
		err = json.Unmarshal(schema, &tempSchema)

		results = append(results, &tempSchema)
	}

	return results, nil
}
