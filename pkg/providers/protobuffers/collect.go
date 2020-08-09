package protobuffers

import (
	"os"
	"path/filepath"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"go.uber.org/zap"
)

// Collect attempts to collect all the available proto files inside the given path and parses them to resources
func Collect(ctx *broker.Context, paths []string, path string) ([]*desc.FileDescriptor, error) {
	imports := make([]string, len(paths))
	for index, path := range paths {
		imports[index] = path
	}

	logger.Debug(ctx, "collect available proto", zap.String("path", path))
	logger.Debug(ctx, "collect available proto with imports", zap.Strings("imports", paths))

	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	logger.Debug(ctx, "absolute proto path", zap.String("path", path))

	for index, path := range imports {
		path, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}

		imports[index] = path
	}

	logger.Debug(ctx, "absolute proto imports", zap.Strings("imports", imports))

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

// ServiceResolver returns a new service(s) resolver for the given protoc collection
func ServiceResolver(imports []string, path string) providers.ServicesResolver {
	return func(ctx *broker.Context) (specs.ServiceList, error) {
		logger.Debug(ctx, "resolving proto services", zap.String("path", path))

		files, err := Collect(ctx, imports, path)
		if err != nil {
			return nil, err
		}

		return NewServices(files), nil
	}
}

// SchemaResolver returns a new schema resolver for the given protoc collection
func SchemaResolver(imports []string, path string) providers.SchemaResolver {
	return func(ctx *broker.Context) (specs.Schemas, error) {
		logger.Debug(ctx, "resolving proto schemas", zap.String("path", path))

		files, err := Collect(ctx, imports, path)
		if err != nil {
			return nil, err
		}

		return NewSchema(files), nil
	}
}

// UnmarshalFiles attempts to parse the given HCL files to intermediate resources.
// Files are parsed based from the given import paths
func UnmarshalFiles(imports []string, files []*providers.FileInfo) ([]*desc.FileDescriptor, error) {
	// NOTE: protoparser expects relative paths, we currently resolved this issue by including root as a import path
	parser := &protoparse.Parser{
		ImportPaths:           append(imports, "/"),
		IncludeSourceCodeInfo: true,
	}

	results := []*desc.FileDescriptor{}

	for _, file := range files {
		descs, err := parser.ParseFiles(file.Path)
		if err != nil {
			return nil, err
		}

		results = append(results, descs...)
	}

	return results, nil
}
