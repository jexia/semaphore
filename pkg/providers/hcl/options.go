package hcl

// Options represents the available options that could be defined inside a HCL definition
type Options struct {
	LogLevel     string
	Protobuffers []string
	Avro         []string
	Openapi3     []string
	GraphQL      *GraphQL
	HTTP         *HTTP
	GRPC         *GRPC
	Prometheus   *Prometheus
	Discovery    *Discovery
}
