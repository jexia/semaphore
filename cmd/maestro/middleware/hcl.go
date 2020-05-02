package middleware

import (
	"path/filepath"

	"github.com/jexia/maestro/pkg/constructor"
	"github.com/jexia/maestro/pkg/definitions/hcl"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/specs/template"
	"github.com/sirupsen/logrus"
)

// ServiceSelector parses the HCL definition on the given path and manipulates the collected services after constructed
func ServiceSelector(path string) constructor.AfterConstructorHandler {
	return func(next constructor.AfterConstructor) constructor.AfterConstructor {
		return func(ctx instance.Context, collection *constructor.Collection) error {
			definitions, err := hcl.ResolvePath(ctx, path)
			if err != nil {
				return err
			}

			for _, definition := range definitions {
				for _, services := range definition.ServiceSelector {
					for _, selector := range services.Selectors {
						for _, service := range collection.Services.Services {
							name := template.JoinPath(service.Package, service.Name)
							matched, err := filepath.Match(selector.Pattern, name)
							if err != nil {
								return err
							}

							ctx.Logger(logger.Core).WithFields(logrus.Fields{
								"pattern": selector.Pattern,
								"service": name,
								"matched": matched,
							}).Debug("pattern matching service")

							if !matched {
								continue
							}

							ctx.Logger(logger.Core).WithFields(logrus.Fields{
								"original": service.Host,
								"new":      selector.Host,
								"selector": selector.Pattern,
								"service":  name,
							}).Info("overriding service configuration")

							service.Host = selector.Host
						}
					}
				}
			}

			return next(ctx, collection)
		}
	}
}
