package core

import (
	"github.com/jexia/semaphore/pkg/core/api"
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/jexia/semaphore/pkg/core/trace"
	"github.com/jexia/semaphore/pkg/transport"
	"github.com/sirupsen/logrus"
)

// NewListeners constructs the listeners from the given collection of endpoints
func NewListeners(endpoints []*transport.Endpoint, options api.Options) error {
	collections := make(map[string][]*transport.Endpoint, len(options.Listeners))

	options.Ctx.Logger(logger.Core).WithField("endpoints", endpoints).Debug("constructing listeners")

	for _, endpoint := range endpoints {
		if endpoint == nil {
			continue
		}

		options.Ctx.Logger(logger.Core).WithFields(logrus.Fields{
			"flow":     endpoint.Flow.GetName(),
			"listener": endpoint.Listener,
		}).Info("Preparing endpoint")

		listener := options.Listeners.Get(endpoint.Listener)
		if listener == nil {
			options.Ctx.Logger(logger.Core).WithFields(logrus.Fields{
				"listener": endpoint.Listener,
			}).Error("Listener not found")

			return trace.New(trace.WithMessage("unknown listener %s", endpoint.Listener))
		}

		collections[endpoint.Listener] = append(collections[endpoint.Listener], endpoint)
	}

	for key, collection := range collections {
		options.Ctx.Logger(logger.Core).WithField("listener", key).Debug("applying listener handles")

		listener := options.Listeners.Get(key)
		err := listener.Handle(options.Ctx, collection, options.Codec)
		if err != nil {
			options.Ctx.Logger(logger.Core).WithFields(logrus.Fields{
				"listener": listener.Name(),
				"err":      err,
			}).Error("Listener returned an error")

			return err
		}
	}

	return nil
}
