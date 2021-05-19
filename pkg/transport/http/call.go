package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/discovery"
	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/transport"
	"go.uber.org/zap"
)

// Call represents the HTTP caller implementation
type Call struct {
	ctx      *broker.Context
	service  string
	host     string
	methods  map[string]*Method
	proxy    *httputil.ReverseProxy
	resolver discovery.Resolver
}

func (call *Call) Address() string {
	addr, ok := call.resolver.Resolve()
	if !ok {
		return ""
	}
	return addr
}

// GetMethods returns the available methods within the HTTP caller
func (call *Call) GetMethods() []transport.Method {
	result := make([]transport.Method, 0, len(call.methods))

	for _, method := range call.methods {
		result = append(result, method)
	}

	return result
}

// GetMethod attempts to return a method matching the given name
func (call *Call) GetMethod(name string) transport.Method {
	for _, method := range call.methods {
		if method.GetName() == name {
			return method
		}
	}

	return nil
}

// SendMsg calls the configured host and attempts to call the given endpoint with the given headers and stream
func (call *Call) SendMsg(ctx context.Context, rw transport.ResponseWriter, pr *transport.Request, refs references.Store) error {
	request := http.MethodGet
	uri, err := url.Parse(call.Address())
	if err != nil {
		return err
	}

	if pr.Method != nil {
		method := call.methods[pr.Method.GetName()]
		if method == nil {
			return ErrUnknownMethod{
				Method:  pr.Method.GetName(),
				Service: call.service,
			}
		}

		endpoint := LookupEndpointReferences(method, refs)
		if endpoint != "" {
			endpointURI, err := url.Parse(endpoint)
			if err != nil {
				return fmt.Errorf("failed to parse endpoint: %w", err)
			}

			uri.Path = endpointURI.Path
			uri.RawQuery = endpointURI.RawQuery
		}

		request = method.request
	}

	logger.Debug(call.ctx, "calling HTTP caller",
		zap.String("uri", uri.String()),
		zap.String("method", request),
	)

	req, err := http.NewRequestWithContext(ctx, request, uri.String(), pr.Body)
	if err != nil {
		return err
	}

	req.Header = CopyMetadataHeader(pr.Header)
	req.Header.Add(ContentTypeHeaderKey, ContentTypes[pr.RequestCodec])
	req.Header.Add(AcceptHeaderKey, ContentTypes[pr.ResponseCodec])

	// TODO: configure http client
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	rw.HeaderStatus(res.StatusCode)
	rw.HeaderMessage(http.StatusText(res.StatusCode))
	AppendHTTPHeader(rw.Header(), res.Header)

	go func() {
		defer rw.Close()
		_, err = io.Copy(rw, res.Body)
		if err != nil && err != io.EOF {
			logger.Debug(call.ctx, "unexpected error while copying response body", zap.Error(err))
		}
	}()

	return nil
}

// Close closes the given caller
func (call *Call) Close() error {
	logger.Info(call.ctx, "closing HTTP caller", zap.String("host", call.host))
	return nil
}
