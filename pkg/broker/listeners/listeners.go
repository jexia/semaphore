package listeners

import (
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/transport"
	"go.uber.org/zap"
)

// Apply collects the given endpoints and applies them to the configured listeners.
// A error is returned if a listener is not found or a listener returned a error.
func Apply(ctx *broker.Context, codec codec.Constructors, listeners transport.ListenerList, endpoints transport.EndpointList) error {
	collections := make(map[string][]*transport.Endpoint, len(listeners))

	logger.Debug(ctx, "constructing listeners")

	for _, endpoint := range endpoints {
		if endpoint == nil {
			continue
		}

		logger.Info(ctx, "preparing endpoint", zap.String("flow", endpoint.Flow.GetName()), zap.String("listener", endpoint.Listener))
		collections[endpoint.Listener] = append(collections[endpoint.Listener], endpoint)
	}

	for key, collection := range collections {
		logger.Debug(ctx, "applying listener handles", zap.String("listener", key))

		listener := listeners.Get(key)
		if listener == nil {
			logger.Error(ctx, "listener not found", zap.String("listener", key))
			return ErrNoListener{Listener: key}
		}

		err := listener.Handle(ctx, collection, codec)
		if err != nil {
			logger.Error(ctx, "listener returned an error", zap.String("listener", listener.Name()), zap.Error(err))
			return err
		}
	}

	return nil
}
