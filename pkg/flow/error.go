package flow

import (
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/metadata"
	"github.com/jexia/semaphore/pkg/specs"
)

// NewOnError constructs a new error for the given codec and header manager
func NewOnError(functions functions.Stack, codec codec.Manager, metadata *metadata.Manager, status, message *specs.Property) *OnError {
	return &OnError{
		functions: functions,
		codec:     codec,
		metadata:  metadata,
		status:    status,
		message:   message,
	}
}

// OnError represents a error codec and metadata manager
type OnError struct {
	functions functions.Stack
	codec     codec.Manager
	metadata  *metadata.Manager
	status    *specs.Property
	message   *specs.Property
}
