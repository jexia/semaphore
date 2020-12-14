package micro

import (
	"github.com/jexia/semaphore/pkg/modules/codec/metadata"
	micrometa "github.com/micro/go-micro/v2/metadata"
)

// CopyMetadataHeader copies the given metadata header to go micro metadata
func CopyMetadataHeader(md metadata.MD) micrometa.Metadata {
	result := make(micrometa.Metadata, len(md))

	for key, val := range md {
		result[key] = val
	}

	return result
}
