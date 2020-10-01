package grpc

import (
	"strings"

	"github.com/jexia/semaphore/pkg/codec/metadata"
	rpcMeta "google.golang.org/grpc/metadata"
)

// CopyRPCMD copies the given grpc metadata into a semaphore metadata
func CopyRPCMD(source rpcMeta.MD) metadata.MD {
	result := metadata.MD{}
	for key, vals := range source {
		result[strings.ToLower(key)] = strings.Join(vals, ";")
	}

	return result
}

// CopyMD copies the given semaphore metadata into a grpc metadata
func CopyMD(source metadata.MD) rpcMeta.MD {
	result := rpcMeta.MD{}
	for key, vals := range source {
		result[strings.ToLower(key)] = []string{vals}
	}

	return result
}
