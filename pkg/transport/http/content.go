package http

// ContentTypeHeaderKey represents the HTTP header content type key
const ContentTypeHeaderKey = "Content-Type"

// ContentType represents a supported content type
type ContentType string

// Available content types
const (
	ApplicationJSON ContentType = "application/json"
)

// ContentTypes represents a lists of available codec content types and their Content-Type value
var ContentTypes = map[string]string{
	"json": string(ApplicationJSON),
}
