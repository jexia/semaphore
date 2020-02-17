package protoc

import (
	"path/filepath"

	"github.com/jexia/maestro/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
)

// ProtoExt file extension
var ProtoExt = ".proto"

// Collect attempts to collect all the available proto files inside the given path and parses them to resources
func Collect(imports []string, path string) (Collection, error) {
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

	files, err := utils.ReadDir(path, true, ProtoExt)
	if err != nil {
		return nil, err
	}

	descriptors, err := UnmarshalFiles(imports, files)
	if err != nil {
		return nil, err
	}

	return NewCollection(descriptors), nil
}

// UnmarshalFiles attempts to parse the given HCL files to intermediate resources.
// Files are parsed based from the given import paths
func UnmarshalFiles(imports []string, files []utils.FileInfo) ([]*desc.FileDescriptor, error) {
	parser := &protoparse.Parser{
		ImportPaths: imports,
	}

	results := []*desc.FileDescriptor{}

	for _, file := range files {
		descs, err := parser.ParseFiles(filepath.Join(file.Path, file.Name()))
		if err != nil {
			return nil, err
		}

		results = append(results, descs...)
	}

	return results, nil
}
