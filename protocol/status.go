package protocol

// Status codes (HTTP) as registered with IANA.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
const (
	StatusOK       = 200
	StatusCreated  = 201
	StatusAccepted = 202

	StatusBadRequest            = 400
	StatusUnauthorized          = 401
	StatusPaymentRequired       = 402
	StatusForbidden             = 403
	StatusNotFound              = 404
	StatusMethodNotAllowed      = 405
	StatusNotAcceptable         = 406
	StatusProxyAuthRequired     = 407
	StatusRequestTimeout        = 408
	StatusConflict              = 409
	StatusRequestEntityTooLarge = 413

	StatusInternalServerError = 500
	StatusNotImplemented      = 501
	StatusBadGateway          = 502
	StatusServiceUnavailable  = 503
	StatusGatewayTimeout      = 504
)

var statusText = map[int]string{
	StatusOK:       "OK",
	StatusCreated:  "Created",
	StatusAccepted: "Accepted",

	StatusBadRequest:            "Bad Request",
	StatusUnauthorized:          "Unauthorized",
	StatusPaymentRequired:       "Payment Required",
	StatusForbidden:             "Forbidden",
	StatusNotFound:              "Not Found",
	StatusNotAcceptable:         "Not Acceptable",
	StatusRequestTimeout:        "Request Timeout",
	StatusConflict:              "Conflict",
	StatusRequestEntityTooLarge: "Request Entity Too Large",

	StatusInternalServerError: "Internal Server Error",
	StatusNotImplemented:      "Not Implemented",
	StatusBadGateway:          "Bad Gateway",
	StatusServiceUnavailable:  "Service Unavailable",
	StatusGatewayTimeout:      "Gateway Timeout",
}

// StatusSuccess checks whether the given status code is a success
func StatusSuccess(code int) bool {
	if code > 200 && code < 300 {
		return true
	}

	return false
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	return statusText[code]
}
