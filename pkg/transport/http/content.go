package http

// ContentTypeHeaderKey represents the HTTP header content type key
const ContentTypeHeaderKey = "Content-Type"

// AcceptHeaderKey represents the HTTP header accept key
const AcceptHeaderKey = "Accept"

// ContentType represents a supported content type
type ContentType string

// Available content types
const (
	ApplicationJSON              ContentType = "application/json"
	ApplicationWWWFormURLEncoded ContentType = "application/x-www-form-urlencoded"
	ApplicationProtobuf          ContentType = "application/protobuf"
)

// ContentTypes represents a lists of available codec content types and their Content-Type value
var ContentTypes = map[string]string{
	"json":            string(ApplicationJSON),
	"form-urlencoded": string(ApplicationWWWFormURLEncoded),
	"proto":           string(ApplicationProtobuf),
}
