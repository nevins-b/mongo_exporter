package collector

import (
	"github.com/nevins-b/commgo"
	"github.com/prometheus/client_golang/prometheus"
)

type cursorCollector struct {
	open, size, timedOut prometheus.Gauge
}

func init() {
	Factories["cursor"] = NewCursorCollector
}

func NewCursorCollector() (Collector, error) {
	return &cursorCollector{
		open: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "cursor_totalopen",
			Help:      "Mongo cursor total time",
		}),
		size: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "cursor_size",
			Help:      "Mongo cursor time",
		}),
		timedOut: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "cursor_timed_out",
			Help:      "Mongo cursor time",
		}),
	}, nil
}

func (c *cursorCollector) Update(ch chan<- prometheus.Metric, status *commgo.ServerStatus) (err error) {

	c.open.Set(float64(status.Cursors.TotalOpen))
	c.size.Set(float64(status.Cursors.ClientCursorSize))
	c.timedOut.Set(float64(status.Cursors.TimedOut))

	c.open.Collect(ch)
	c.size.Collect(ch)
	c.timedOut.Collect(ch)
	return err
}
