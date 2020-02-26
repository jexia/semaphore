package services

import (
	"context"

	"github.com/jexia/maestro/refs"
)

// Call represents a caller which could be called
type Call func(context.Context, *refs.Store) error

// Collection represents a collection of services
type Collection map[string]*Service

// Get attempts to fetch the given service by name
func (collection Collection) Get(name string) *Service {
	return collection[name]
}

// Service represents a flow service
type Service struct {
	Call     Call
	Rollback Call
}
