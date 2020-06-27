package middleware

import (
	"path/filepath"

	"github.com/jexia/maestro/internal/definitions/hcl"
	"github.com/jexia/maestro/pkg/core/api"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/core/logger"
	"github.com/jexia/maestro/pkg/specs/template"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
)

// ServiceSelector parses the HCL definition on the given path and manipulates the collected services after constructed
func ServiceSelector(path string) api.AfterConstructorHandler {
	return func(next api.AfterConstructor) api.AfterConstructor {
		return func(ctx instance.Context, collection *api.Collection) error {
			definitions, err := hcl.ResolvePath(ctx, []string{}, path)
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

							attrs, _ := selector.Options.JustAttributes()
							for _, attr := range attrs {
								value, _ := attr.Expr.Value(nil)
								if value.Type() != cty.String {
									continue
								}

								service.Options[attr.Name] = value.AsString()
							}
						}
					}
				}
			}

			return next(ctx, collection)
		}
	}
}
