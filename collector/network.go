package collector

import (
	"github.com/nevins-b/commgo"
	"github.com/prometheus/client_golang/prometheus"
)

type networkCollector struct {
	in, out, requests prometheus.Gauge
}

func init() {
	Factories["network"] = NewNetworkCollector
}

func NewNetworkCollector() (Collector, error) {
	return &networkCollector{
		in: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "network_bytes_in",
			Help:      "Mongo Network Bytes in",
		}),
		out: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "network_bytes_out",
			Help:      "Mongo Network Bytes in",
		}),
		requests: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "network_requests",
			Help:      "Mongo Network requests",
		}),
	}, nil
}

func (c *networkCollector) Update(ch chan<- prometheus.Metric, status *commgo.ServerStatus) (err error) {

	c.in.Set(float64(status.Network.BytesIn))
	c.out.Set(float64(status.Network.BytesOut))
	c.requests.Set(float64(status.Network.NumRequests))

	c.in.Collect(ch)
	c.out.Collect(ch)
	c.requests.Collect(ch)
	return err
}
