// Exporter is a prometheus exporter using multiple Factories to collect and export system metrics.
package collector

import (
	"github.com/nevins-b/commgo"
	"github.com/prometheus/client_golang/prometheus"
)

const Namespace = "mongo"

var Factories = make(map[string]func() (Collector, error))

// Interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update(ch chan<- prometheus.Metric, status *commgo.ServerStatus) (err error)
}

// TODO: Instead of periodically call Update, a Collector could be implemented
// as a real prometheus.Collector that only gathers metrics when
// scraped. (However, for metric gathering that takes very long, it might
// actually be better to do them proactively before scraping to minimize scrape
// time.)
