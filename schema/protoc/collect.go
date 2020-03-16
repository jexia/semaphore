package protoc

import (
	"os"
	"path/filepath"

	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
)

// Collect attempts to collect all the available proto files inside the given path and parses them to resources
func Collect(imports []string, path string) (schema.Resolver, error) {
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

	files, err := utils.ResolvePath(path)
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

	collection := NewCollection(descriptors)
	return SchemaResolver(collection), nil
}

// SchemaResolver returns a new schema resolver for the given protoc collection
func SchemaResolver(collection schema.Collection) schema.Resolver {
	return func(schemas *schema.Store) error {
		schemas.Add(collection)
		return nil
	}
}

// UnmarshalFiles attempts to parse the given HCL files to intermediate resources.
// Files are parsed based from the given import paths
func UnmarshalFiles(imports []string, files []*utils.FileInfo) ([]*desc.FileDescriptor, error) {
	parser := &protoparse.Parser{
		ImportPaths:           imports,
		IncludeSourceCodeInfo: true,
	}

	results := []*desc.FileDescriptor{}

	for _, file := range files {
		descs, err := parser.ParseFiles(file.AbsolutePath)
		if err != nil {
			return nil, err
		}

		results = append(results, descs...)
	}

	return results, nil
}
