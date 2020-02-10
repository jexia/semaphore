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
func Collect(imports []string, path string) ([]*desc.FileDescriptor, error) {
	files, err := utils.ReadDir(path, true, ProtoExt)
	if err != nil {
		return nil, err
	}

	descriptors, err := UnmarshalFiles(imports, path, files)
	if err != nil {
		return nil, err
	}

	return descriptors, nil
}

// UnmarshalFiles attempts to parse the given HCL files to intermediate resources.
// Files are parsed based from the given import paths
func UnmarshalFiles(imports []string, path string, files []utils.FileInfo) ([]*desc.FileDescriptor, error) {
	parser := &protoparse.Parser{
		ImportPaths: imports,
	}

	results := []*desc.FileDescriptor{}

	for _, file := range files {
		descs, err := parser.ParseFiles(utils.RelativePath(path, filepath.Join(file.Path, file.Name())))
		if err != nil {
			return nil, err
		}

		results = append(results, descs...)
	}

	return results, nil
}

// GetService attempts to find the given service by its fully qualified name
func GetService(descriptors []*desc.FileDescriptor, service string) *desc.ServiceDescriptor {
	for _, descriptor := range descriptors {
		service := descriptor.FindService(service)
		if service == nil {
			continue
		}

		return service
	}

	return nil
}

// GetMethod attempts to find the given method by its name
func GetMethod(service *desc.ServiceDescriptor, name string) *desc.MethodDescriptor {
	for _, method := range service.GetMethods() {
		if method.GetName() != name {
			continue
		}

		return method
	}

	return nil
}
