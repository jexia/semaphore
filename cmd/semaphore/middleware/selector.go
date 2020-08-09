package middleware

import (
	"path/filepath"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/core/api"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
)

// ServiceSelector parses the HCL definition on the given path and manipulates the collected services after constructed
func ServiceSelector(path string) api.AfterConstructorHandler {
	return func(next api.AfterConstructor) api.AfterConstructor {
		return func(ctx *broker.Context, flows specs.FlowListInterface, endpoints specs.EndpointList, services specs.ServiceList, schemas specs.Schemas) error {
			definitions, err := hcl.ResolvePath(ctx, []string{}, path)
			if err != nil {
				return err
			}

			for _, definition := range definitions {
				for _, srvs := range definition.ServiceSelector {
					for _, selector := range srvs.Selectors {
						for _, service := range services {
							name := template.JoinPath(service.Package, service.Name)
							matched, err := filepath.Match(selector.Pattern, name)
							if err != nil {
								return err
							}

							logger.Debug(ctx, "pattern matching service",
								zap.String("pattern", selector.Pattern),
								zap.String("service", name),
								zap.Bool("matched", matched),
							)

							if !matched {
								continue
							}

							logger.Info(ctx, "overriding service configuration",
								zap.String("host", selector.Host),
								zap.String("pattern", selector.Pattern),
								zap.String("codec", selector.Codec),
								zap.String("transport", selector.Transport),
								zap.String("service", name),
							)

							if selector.Host != "" {
								service.Host = selector.Host
							}

							if selector.Transport != "" {
								service.Transport = selector.Transport
							}

							if selector.Codec != "" {
								service.Codec = selector.Codec
							}

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

			return next(ctx, flows, endpoints, services, schemas)
		}
	}
}
