package listeners

import (
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/core/api"
	"github.com/jexia/semaphore/pkg/core/trace"
	"github.com/jexia/semaphore/pkg/transport"
	"go.uber.org/zap"
)

// Apply collects the given endpoints and applies them to the configured listeners.
// A error is returned if a listener is not found or a listener returned a error.
func Apply(endpoints []*transport.Endpoint, options api.Options) error {
	collections := make(map[string][]*transport.Endpoint, len(options.Listeners))

	logger.Debug(options.Ctx, "constructing listeners")

	for _, endpoint := range endpoints {
		if endpoint == nil {
			continue
		}

		logger.Info(options.Ctx, "preparing endpoint", zap.String("flow", endpoint.Flow.GetName()), zap.String("listener", endpoint.Listener))

		listener := options.Listeners.Get(endpoint.Listener)
		if listener == nil {
			logger.Error(options.Ctx, "listener not found", zap.String("listener", endpoint.Listener))
			return trace.New(trace.WithMessage("unknown listener %s", endpoint.Listener))
		}

		collections[endpoint.Listener] = append(collections[endpoint.Listener], endpoint)
	}

	for key, collection := range collections {
		logger.Debug(options.Ctx, "applying listener handles", zap.String("listener", key))

		listener := options.Listeners.Get(key)
		err := listener.Handle(options.Ctx, collection, options.Codec)
		if err != nil {
			logger.Error(options.Ctx, "listener returned an error", zap.String("listener", listener.Name()), zap.Error(err))
			return err
		}
	}

	return nil
}
