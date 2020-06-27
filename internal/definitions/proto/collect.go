package proto

import (
	"os"
	"path/filepath"

	"github.com/jexia/maestro/internal/definitions"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
)

// Collect attempts to collect all the available proto files inside the given path and parses them to resources
func Collect(ctx instance.Context, paths []string, path string) ([]*desc.FileDescriptor, error) {
	imports := make([]string, len(paths))
	for index, path := range paths {
		imports[index] = path
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	for index, path := range imports {
		path, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}

		imports[index] = path
	}

	files, err := definitions.ResolvePath(ctx, []string{}, path)
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
func ServiceResolver(imports []string, path string) definitions.ServicesResolver {
	return func(ctx instance.Context) ([]*specs.ServicesManifest, error) {
		files, err := Collect(ctx, imports, path)
		if err != nil {
			return nil, err
		}

		return NewServices(files), nil
	}
}

// SchemaResolver returns a new schema resolver for the given protoc collection
func SchemaResolver(imports []string, path string) definitions.SchemaResolver {
	return func(ctx instance.Context) ([]*specs.SchemaManifest, error) {
		files, err := Collect(ctx, imports, path)
		if err != nil {
			return nil, err
		}

		return NewSchema(files), nil
	}
}

// UnmarshalFiles attempts to parse the given HCL files to intermediate resources.
// Files are parsed based from the given import paths
func UnmarshalFiles(imports []string, files []*definitions.FileInfo) ([]*desc.FileDescriptor, error) {
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
