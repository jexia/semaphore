package protobuffers

import (
	"os"
	"path/filepath"

	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
)

// Collect attempts to collect all the available proto files inside the given path and parses them to resources
func Collect(ctx instance.Context, paths []string, path string) ([]*desc.FileDescriptor, error) {
	imports := make([]string, len(paths))
	for index, path := range paths {
		imports[index] = path
	}

	ctx.Logger(logger.Core).WithField("path", path).Debug("Collect available proto")
	ctx.Logger(logger.Core).WithField("imports", paths).Debug("Collect available proto with imports")

	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	ctx.Logger(logger.Core).WithField("path", path).Debug("Absolute proto path")

	for index, path := range imports {
		path, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}

		imports[index] = path
	}

	ctx.Logger(logger.Core).WithField("imports", paths).Debug("Absolute proto imports")

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
	return func(ctx instance.Context) ([]*specs.ServicesManifest, error) {
		ctx.Logger(logger.Core).WithField("path", path).Debug("Resolving proto services")

		files, err := Collect(ctx, imports, path)
		if err != nil {
			return nil, err
		}

		return NewServices(files), nil
	}
}

// SchemaResolver returns a new schema resolver for the given protoc collection
func SchemaResolver(imports []string, path string) providers.SchemaResolver {
	return func(ctx instance.Context) (specs.Objects, error) {
		ctx.Logger(logger.Core).WithField("path", path).Debug("Resolving proto schemas")

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
