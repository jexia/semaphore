package http

// ContentTypeHeaderKey represents the HTTP header content type key
const ContentTypeHeaderKey = "Content-Type"

// ContentTypes represents a lists of available codec content types and their Content-Type value
var ContentTypes = map[string]string{
	"json": "application/json",
}
