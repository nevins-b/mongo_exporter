package collector

import (
	"github.com/nevins-b/commgo"
	"github.com/prometheus/client_golang/prometheus"
)

type lockCollector struct {
	total, time prometheus.Gauge
}

func init() {
	Factories["lock"] = NewLockCollector
}

func NewLockCollector() (Collector, error) {
	return &lockCollector{
		total: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "lock_total_time",
			Help:      "Mongo lock total time",
		}),
		time: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "lock_time",
			Help:      "Mongo lock time",
		}),
	}, nil
}

func (c *lockCollector) Update(ch chan<- prometheus.Metric, status *commgo.ServerStatus) (err error) {

	c.total.Set(float64(status.GlobalLock.TotalTime))
	c.time.Set(float64(status.GlobalLock.LockTime))

	c.total.Collect(ch)
	c.time.Collect(ch)
	return err
}
