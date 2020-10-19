package prometheus

import (
	"context"
	"net/http"
	"time"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// New constructs a new prometheus middleware instance
func New(addr string) semaphore.Middleware {
	return semaphore.MiddlewareFunc(func(parent *broker.Context) ([]semaphore.Option, error) {
		ctx := logger.WithLogger(broker.WithModule(parent, "prometheus"))
		logger.Info(ctx, "setting up", zap.String("addr", addr))

		collector, err := NewCollector()
		if err != nil {
			return nil, err
		}

		server := http.Server{
			Addr:    addr,
			Handler: promhttp.HandlerFor(collector.Registry(), promhttp.HandlerOpts{}),
		}

		go server.ListenAndServe()

		handles := semaphore.NewCollection(
			semaphore.BeforeManagerDo(collector.BeforeDo),
			semaphore.AfterManagerDo(collector.AfterDo),
			semaphore.BeforeNodeDo(collector.BeforeNode),
			semaphore.BeforeNodeRollback(collector.BeforeNode),
			semaphore.AfterNodeDo(collector.AfterNode),
			semaphore.AfterNodeRollback(collector.AfterNode),
			semaphore.BeforeManagerRollback(collector.BeforeRollback),
			semaphore.AfterManagerRollback(collector.AfterRollback),
		)

		promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{},
		)

		return handles, nil
	})
}

// CtxKey context key type
type CtxKey string

var (
	// StartTimeCtx context time value
	StartTimeCtx = CtxKey("start-time")
)

// NewCollector constructs a new prometheus collector
func NewCollector() (Collector, error) {
	registry := prometheus.NewRegistry()
	collector := &collector{
		registry: registry,
		flowsDo: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "flow_total",
			Help: "Total amount of times a flow has been called",
		}, []string{"flow"}),
		flowDo: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name: "flow_duration_seconds",
			Help: "Avarage flow execution duration in seconds",
		}, []string{"flow"}),
		flowsRollback: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "flow_rollback_total",
			Help: "Total amount of times a flow has been rolled back",
		}, []string{"flow"}),
		flowRollback: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name: "flow_rollback_duration_seconds",
			Help: "Avarage rollback execution duration in seconds",
		}, []string{"flow"}),
		nodes: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "node_total",
			Help: "Total amount of times a node has been called",
		}, []string{"flow", "node"}),
		node: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name: "node_duration_seconds",
			Help: "Avarage node execution duration in seconds",
		}, []string{"flow", "node"}),
	}

	err := registry.Register(collector.flowsDo)
	if err != nil {
		return nil, err
	}

	err = registry.Register(collector.flowDo)
	if err != nil {
		return nil, err
	}

	err = registry.Register(collector.flowsRollback)
	if err != nil {
		return nil, err
	}

	err = registry.Register(collector.flowRollback)
	if err != nil {
		return nil, err
	}

	err = registry.Register(collector.nodes)
	if err != nil {
		return nil, err
	}

	err = registry.Register(collector.node)
	if err != nil {
		return nil, err
	}

	return collector, nil
}

// Collector represents a middleware collector
type Collector interface {
	Registry() *prometheus.Registry
	BeforeDo(next flow.BeforeManager) flow.BeforeManager
	AfterDo(next flow.AfterManager) flow.AfterManager
	BeforeNode(next flow.BeforeNode) flow.BeforeNode
	AfterNode(next flow.AfterNode) flow.AfterNode
	BeforeRollback(next flow.BeforeManager) flow.BeforeManager
	AfterRollback(next flow.AfterManager) flow.AfterManager
}

// Collector collects data from middleware calls and exposes them for prometheus to consume
type collector struct {
	registry      *prometheus.Registry
	flowsDo       *prometheus.CounterVec
	flowDo        *prometheus.SummaryVec
	flowsRollback *prometheus.CounterVec
	flowRollback  *prometheus.SummaryVec
	nodes         *prometheus.CounterVec
	node          *prometheus.SummaryVec
}

func (collector *collector) Registry() *prometheus.Registry {
	return collector.registry
}

// BeforeDo gets called before a flow gets executed
func (collector *collector) BeforeDo(next flow.BeforeManager) flow.BeforeManager {
	return func(ctx context.Context, manager *flow.Manager, store references.Store) (context.Context, error) {
		req := collector.flowsDo.With(prometheus.Labels{
			"flow": manager.Name,
		})

		req.Inc()

		now := time.Now()
		ctx = context.WithValue(ctx, StartTimeCtx, now)

		return next(ctx, manager, store)
	}
}

// AfterDo gets called after a flow is executed
func (collector *collector) AfterDo(next flow.AfterManager) flow.AfterManager {
	return func(ctx context.Context, manager *flow.Manager, store references.Store) (context.Context, error) {
		value := ctx.Value(StartTimeCtx)
		if value != nil {
			duration := collector.flowDo.With(prometheus.Labels{
				"flow": manager.Name,
			})

			start := value.(time.Time)
			diff := time.Since(start)

			duration.Observe(diff.Seconds())
		}

		return next(ctx, manager, store)
	}
}

// BeforeRollback gets called before a flow rollback gets executed
func (collector *collector) BeforeRollback(next flow.BeforeManager) flow.BeforeManager {
	return func(ctx context.Context, manager *flow.Manager, store references.Store) (context.Context, error) {
		req := collector.flowsRollback.With(prometheus.Labels{
			"flow": manager.Name,
		})

		req.Inc()

		now := time.Now()
		ctx = context.WithValue(ctx, StartTimeCtx, now)

		return next(ctx, manager, store)
	}
}

// AfterRollback gets called before a flow rollback is executed
func (collector *collector) AfterRollback(next flow.AfterManager) flow.AfterManager {
	return func(ctx context.Context, manager *flow.Manager, store references.Store) (context.Context, error) {
		value := ctx.Value(StartTimeCtx)
		if value != nil {
			duration := collector.flowRollback.With(prometheus.Labels{
				"flow": manager.Name,
			})

			start := value.(time.Time)
			diff := time.Since(start)

			duration.Observe(diff.Seconds())
		}

		return next(ctx, manager, store)
	}
}

// BeforeDo gets called before a node gets executed
func (collector *collector) BeforeNode(next flow.BeforeNode) flow.BeforeNode {
	return func(ctx context.Context, node *flow.Node, tracker flow.Tracker, processes *flow.Processes, store references.Store) (context.Context, error) {
		req := collector.nodes.With(prometheus.Labels{
			"flow": tracker.Flow(),
			"node": node.Name,
		})

		req.Inc()

		now := time.Now()
		ctx = context.WithValue(ctx, StartTimeCtx, now)

		return next(ctx, node, tracker, processes, store)
	}
}

// AfterDo gets called after a node is executed
func (collector *collector) AfterNode(next flow.AfterNode) flow.AfterNode {
	return func(ctx context.Context, node *flow.Node, tracker flow.Tracker, processes *flow.Processes, store references.Store) (context.Context, error) {
		value := ctx.Value(StartTimeCtx)
		if value != nil {
			duration := collector.node.With(prometheus.Labels{
				"flow": tracker.Flow(),
				"node": node.Name,
			})

			start := value.(time.Time)
			diff := time.Since(start)

			duration.Observe(diff.Seconds())
		}

		return next(ctx, node, tracker, processes, store)
	}
}
