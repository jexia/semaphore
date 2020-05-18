package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

// NewProxy constructs a new reverse proxy with the given options
func NewProxy(options *CallerOptions) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director:      func(*http.Request) {},
		FlushInterval: options.FlushInterval,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   options.Timeout,
				KeepAlive: options.KeepAlive,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          options.MaxIdleConns,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: options.Insecure,
			},
		},
	}
}
