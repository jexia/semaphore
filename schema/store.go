package schema

import (
	"context"

	"github.com/jexia/maestro/logger"
)

// NewStore constructs a new schema store
func NewStore(ctx context.Context) *Store {
	return &Store{
		ctx:      ctx,
		services: make(map[string]Service),
		messages: make(map[string]Property),
	}
}

// Store represents a schema collection store
type Store struct {
	ctx      context.Context
	services map[string]Service
	messages map[string]Property
}

// GetService attempts to return a service with the given name
func (store *Store) GetService(name string) Service {
	return store.services[name]
}

// GetServices returns all available services within the given store
func (store *Store) GetServices() []Service {
	result := make([]Service, len(store.services))

	for _, service := range store.services {
		result = append(result, service)
	}

	return result
}

// GetMessage attempts to return a message with the given name
func (store *Store) GetMessage(name string) Property {
	return store.messages[name]
}

// GetMessages returns all available messages within the given store
func (store *Store) GetMessages() []Property {
	result := make([]Property, len(store.messages))

	for _, message := range store.messages {
		result = append(result, message)
	}

	return result
}

// Add appends the given collection to the existing collection
func (store *Store) Add(collection Collection) {
	if collection == nil {
		return
	}

	logger.FromCtx(store.ctx, logger.Core).WithField("collection", collection).Debug("Appending schema collection to schema store")

	for _, service := range collection.GetServices() {
		if service == nil {
			continue
		}

		logger.FromCtx(store.ctx, logger.Core).WithField("service", service.GetName()).Debug("Appending service to schema store")
		store.services[service.GetFullyQualifiedName()] = service
	}

	for _, message := range collection.GetMessages() {
		if message == nil {
			continue
		}

		logger.FromCtx(store.ctx, logger.Core).WithField("message", message.GetName()).Debug("Appending message to schema store")
		store.messages[message.GetName()] = message
	}
}
