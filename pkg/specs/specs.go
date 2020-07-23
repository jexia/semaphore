package specs

// Options represents a collection of options
type Options map[string]string

// Header represents a collection of key values
type Header map[string]*Property

// Collection represents a collection of flow, endpoint, services and schema manifests.
// This collection is used to define the flow managers and function definitions.
// The references and types defined within this collection are not type checked.
type Collection struct {
	*FlowsManifest
	*EndpointsManifest
	*ServicesManifest
}
