package flow

import (
	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/metadata"
)

// NewOnError constructs a new error for the given codec and header manager
func NewOnError(functions functions.Stack, codec codec.Manager, metadata *metadata.Manager) *OnError {
	return &OnError{
		functions: functions,
		codec:     codec,
		metadata:  metadata,
	}
}

// OnError represents a error codec and metadata manager
type OnError struct {
	functions functions.Stack
	codec     codec.Manager
	metadata  *metadata.Manager
}
