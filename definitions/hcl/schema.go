package hcl

import (
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs/trace"
	log "github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
)

type collection struct {
	services []schema.Service
}

func (collection *collection) GetService(name string) schema.Service {
	for _, service := range collection.services {
		if service.GetName() == name {
			return service
		}
	}

	return nil
}

func (collection *collection) GetServices() []schema.Service {
	return collection.services
}

func (collection *collection) GetMessage(name string) schema.Property {
	return nil
}

func (collection *collection) GetMessages() []schema.Property {
	return make([]schema.Property, 0)
}

// ParseSchema parses the given intermediate manifest to a schema
func ParseSchema(manifest Manifest, schemas schema.Collection) (schema.Collection, error) {
	log.Info("Parsing intermediate manifest to schema")

	result := &collection{
		services: make([]schema.Service, len(manifest.Services)),
	}

	for index, intermediate := range manifest.Services {
		service, err := ParseIntermediateService(intermediate, schemas)
		if err != nil {
			return nil, err
		}

		result.services[index] = service
	}

	return result, nil
}

// service represents a schema service
type service struct {
	name          string
	documentation string
	host          string
	protocol      string
	codec         string
	methods       []schema.Method
	options       schema.Options
}

// GetName returns the service name
func (service *service) GetName() string {
	return service.name
}

// GetComment returns the service documentation
func (service *service) GetComment() string {
	return service.documentation
}

// GetHost returns the service host
func (service *service) GetHost() string {
	return service.host
}

// GetProtocol returns the service protocol
func (service *service) GetProtocol() string {
	return service.protocol
}

// GetCodec returns the service codec
func (service *service) GetCodec() string {
	return service.codec
}

// GetOptions returns the service options
func (service *service) GetOptions() schema.Options {
	return service.options
}

// GetMethod attempts to find a method with the given name
func (service *service) GetMethod(name string) schema.Method {
	for _, method := range service.methods {
		if method.GetName() == name {
			return method
		}
	}

	return nil
}

// GetMethods returns the available methods within the given service
func (service *service) GetMethods() schema.Methods {
	return service.methods
}

// ParseIntermediateService parses the given intermediate service to a specs service
func ParseIntermediateService(manifest Service, collection schema.Collection) (schema.Service, error) {
	log.WithField("service", manifest.Name).Debug("Parsing intermediate service to schema")

	methods, err := ParseIntermediateMethods(manifest.Methods, collection)
	if err != nil {
		return nil, err
	}

	result := &service{
		name:     manifest.Name,
		protocol: manifest.Protocol,
		host:     manifest.Host,
		codec:    manifest.Codec,
		methods:  methods,
		options:  ParseIntermediateSchemaOptions(manifest.Options),
	}

	return result, nil
}

type method struct {
	name          string
	documentation string
	request       schema.Property
	response      schema.Property
	options       schema.Options
}

func (method *method) GetName() string {
	return method.name
}

func (method *method) GetComment() string {
	return method.documentation
}

func (method *method) GetInput() schema.Property {
	return method.request
}

func (method *method) GetOutput() schema.Property {
	return method.response
}

func (method *method) GetOptions() schema.Options {
	return method.options
}

// ParseIntermediateMethods parses the given methods for the given service
func ParseIntermediateMethods(methods []Method, collection schema.Collection) ([]schema.Method, error) {
	result := make([]schema.Method, len(methods))

	for index, manifest := range methods {
		log.WithFields(log.Fields{
			"method": manifest.Name,
		}).Debug("Parsing intermediate method to schema")

		request := collection.GetMessage(manifest.Request)
		if request == nil && manifest.Request != "" {
			return nil, trace.New(trace.WithMessage("undefined request method '%s' inside schema collection", manifest.Request))
		}

		response := collection.GetMessage(manifest.Response)
		if response == nil && manifest.Response != "" {
			return nil, trace.New(trace.WithMessage("undefined response method '%s' inside schema collection", manifest.Response))
		}

		result[index] = &method{
			name:     manifest.Name,
			request:  request,
			response: response,
			options:  ParseIntermediateSchemaOptions(manifest.Options),
		}
	}

	return result, nil
}

// ParseIntermediateSchemaOptions parses the given intermediate options to a schema options
func ParseIntermediateSchemaOptions(options *Options) schema.Options {
	if options == nil {
		return schema.Options{}
	}

	result := schema.Options{}
	attrs, _ := options.Body.JustAttributes()

	for key, val := range attrs {
		val, _ := val.Expr.Value(nil)
		if val.Type() != cty.String {
			continue
		}

		result[key] = val.AsString()
	}

	return result
}
