package grpc

import (
	"strings"

	"github.com/jexia/maestro/pkg/metadata"
	rpcMeta "google.golang.org/grpc/metadata"
)

// CopyRPCMD copies the given grpc metadata into a maestro metadata
func CopyRPCMD(source rpcMeta.MD) metadata.MD {
	result := metadata.MD{}
	for key, vals := range source {
		result[key] = strings.Join(vals, ";")
	}

	return result
}

// CopyMD copies the given maestro metadata into a grpc metadata
func CopyMD(source metadata.MD) rpcMeta.MD {
	result := rpcMeta.MD{}
	for key, vals := range source {
		result[key] = []string{vals}
	}

	return result
}
