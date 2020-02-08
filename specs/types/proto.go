package types

import "github.com/jhump/protoreflect/desc"

// FindService attempts to find the given service by its fully qualified name
func FindService(descriptors []*desc.FileDescriptor, fqn string) *desc.ServiceDescriptor {
	for _, descriptor := range descriptors {
		service := descriptor.FindService(fqn)
		if service == nil {
			continue
		}

		return service
	}

	return nil
}

// FindMethod attempts to find the given method by its fully qualified name
func FindMethod(service *desc.ServiceDescriptor, fqn string) *desc.MethodDescriptor {
	for _, method := range service.GetMethods() {
		if method.GetName() != fqn {
			continue
		}

		return method
	}

	return nil
}
