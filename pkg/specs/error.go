package specs

// ErrorHandle represents a error handle object
type ErrorHandle interface {
	GetError() *ParameterMap
	GetOnError() *OnError
}
