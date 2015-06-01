package collector

import (
	"github.com/nevins-b/commgo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type connectionCollector struct {
	available, current, total prometheus.Gauge
}

func init() {
	Factories["connection"] = NewConnectionsCollector
}

func NewConnectionsCollector() (Collector, error) {
	return &connectionCollector{
		available: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "connections_available",
			Help:      "Mongo connections available",
		}),
		current: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "connections_current",
			Help:      "Mongo connections currently used",
		}),
		total: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "connections_total",
			Help:      "Mongo connections created",
		}),
	}, nil
}

func (c *connectionCollector) Update(ch chan<- prometheus.Metric, session *mgo.Session) (err error) {

	cmd := &bson.M{
		"serverStatus": 1,
	}

	result := &commgo.ServerStatus{}

	if err := session.DB("local").Run(&cmd, &result); err != nil {
		log.Errorf("%v", err)
		return err
	}

	c.available.Set(float64(result.Connections.Available))
	c.current.Set(float64(result.Connections.Current))
	c.total.Set(float64(result.Connections.TotalCreated))

	c.available.Collect(ch)
	c.current.Collect(ch)
	c.total.Collect(ch)
	return err
}
