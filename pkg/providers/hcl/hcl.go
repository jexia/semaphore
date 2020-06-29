package hcl

import "github.com/jexia/maestro/internal/definitions/hcl"

// FlowsResolver constructs a resource resolver for the given path
var FlowsResolver = hcl.FlowsResolver

// ServicesResolver constructs a schema resolver for the given path.
// The HCL schema resolver relies on other schema registries.
// Those need to be resolved before the HCL schemas are resolved.
var ServicesResolver = hcl.ServicesResolver

// EndpointsResolver constructs a resource resolver for the given path
var EndpointsResolver = hcl.EndpointsResolver

// GetOptions returns the defined options inside the given path
var GetOptions = hcl.GetOptions
