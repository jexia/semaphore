package flow

import (
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/modules/codec/metadata"
	"github.com/jexia/semaphore/pkg/specs"
)

// NewOnError constructs a new error for the given codec and header manager
func NewOnError(stack functions.Stack, codec codec.Manager, metadata *metadata.Manager, err *specs.OnError) *OnError {
	result := &OnError{
		stack:    stack,
		codec:    codec,
		metadata: metadata,
	}

	if err != nil {
		result.status = err.Status
		result.message = err.Message
	}

	return result
}

// OnError represents a error codec and metadata manager
type OnError struct {
	stack    functions.Stack
	codec    codec.Manager
	metadata *metadata.Manager
	status   *specs.Property
	message  *specs.Property
}
