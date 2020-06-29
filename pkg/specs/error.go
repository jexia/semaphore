package specs

// ErrorHandle represents a error handle object
type ErrorHandle interface {
	GetResponse() *ParameterMap
	GetStatusCode() *Property
	GetMessage() *Property
}
