package services

import (
	"context"
	"io"

	"github.com/jexia/maestro/codec"
)

// Call represents a caller which could be called
type Call func(context.Context, io.Reader) (io.Reader, error)

// Collection represents a collection of services
type Collection map[string]*Service

// Get attempts to fetch the given service by name
func (collection Collection) Get(name string) *Service {
	return collection[name]
}

// Service represents a flow service
type Service struct {
	Codec    codec.Manager
	Call     Call
	Rollback Call
}
