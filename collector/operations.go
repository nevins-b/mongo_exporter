package collector

import (
	"github.com/nevins-b/commgo"
	"github.com/prometheus/client_golang/prometheus"
)

type operationCollector struct {
	command, delete, getmore, insert, query, update prometheus.Gauge
}

func init() {
	Factories["operations"] = NewOperationCollector
}

func NewOperationCollector() (Collector, error) {
	return &operationCollector{
		command: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "operation_command",
			Help:      "Mongo OpcounterStats available",
		}),
		delete: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "operation_current",
			Help:      "Mongo OpcounterStats currently used",
		}),
		getmore: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "operation_getmore",
			Help:      "Mongo OpcounterStats created",
		}),
		insert: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "operation_insert",
			Help:      "Mongo OpcounterStats available",
		}),
		query: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "operation_query",
			Help:      "Mongo OpcounterStats available",
		}),
		update: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "operation_update",
			Help:      "Mongo OpcounterStats available",
		}),
	}, nil
}

func (c *operationCollector) Update(ch chan<- prometheus.Metric, status *commgo.ServerStatus) (err error) {
	c.command.Set(float64(status.Opcounters.Command))
	c.delete.Set(float64(status.Opcounters.Delete))
	c.getmore.Set(float64(status.Opcounters.Getmore))
	c.insert.Set(float64(status.Opcounters.Insert))
	c.query.Set(float64(status.Opcounters.Query))
	c.update.Set(float64(status.Opcounters.Update))

	c.command.Collect(ch)
	c.delete.Collect(ch)
	c.getmore.Collect(ch)
	c.insert.Collect(ch)
	c.query.Collect(ch)
	c.update.Collect(ch)
	return err
}
