package definitions

import (
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/specs"
)

// FlowsResolver when called collects the available flow(s) with the configured configuration
type FlowsResolver func(instance.Context) (*specs.FlowsManifest, error)

// EndpointsResolver when called collects the available endpoint(s) with the configured configuration
type EndpointsResolver func(instance.Context) (*specs.EndpointsManifest, error)

// ServicesResolver when called collects the available service(s) with the configured configuration
type ServicesResolver func(instance.Context) (*specs.ServicesManifest, error)

// SchemaResolver when called collects the available service(s) with the configured configuration
type SchemaResolver func(instance.Context) (*specs.SchemaManifest, error)
