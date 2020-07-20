package config

import (
	"github.com/jexia/semaphore/pkg/providers/hcl"
)

// New constructs a new global config
func New() *Semaphore {
	return &Semaphore{
		HTTP:         HTTP{},
		GraphQL:      GraphQL{},
		GRPC:         GRPC{},
		Protobuffers: []string{},
		Files:        []string{},
	}
}

// Parse parses the given HCL options definition
func Parse(options *hcl.Options, target *Semaphore) {
	if target.LogLevel == "" && options.LogLevel != "" {
		target.LogLevel = options.LogLevel
	}

	if len(options.Protobuffers) > 0 {
		target.Protobuffers = append(target.Protobuffers, options.Protobuffers...)
	}

	if options.GraphQL != nil && target.GraphQL.Address == "" {
		target.GraphQL = GraphQL{
			Address: options.GraphQL.Address,
		}
	}

	if options.HTTP != nil && target.HTTP.Address == "" {
		target.HTTP = HTTP{
			Address: options.HTTP.Address,
		}
	}

	if options.GRPC != nil && target.GRPC.Address == "" {
		target.GRPC = GRPC{
			Address: options.GRPC.Address,
		}
	}

	if options.Prometheus != nil && target.Prometheus.Address == "" {
		target.Prometheus = Prometheus{
			Address: options.Prometheus.Address,
		}
	}
}

// Semaphore configurations
type Semaphore struct {
	LogLevel     string
	HTTP         HTTP
	GraphQL      GraphQL
	GRPC         GRPC
	Prometheus   Prometheus
	Protobuffers []string
	Files        []string
}

// Prometheus configurations
type Prometheus struct {
	Address string
}

// HTTP configurations
type HTTP struct {
	Address string
}

// GRPC configurations
type GRPC struct {
	Address string
}

// GraphQL configurations
type GraphQL struct {
	Address string
}
