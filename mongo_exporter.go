package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"
	"time"

	"github.com/nevins-b/commgo"
	"github.com/nevins-b/mongo_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const subsystem = "exporter"

var (
	listenAddress     = flag.String("web.listen-address", ":9107", "Address to listen on for web interface and telemetry.")
	metricsPath       = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	mongoServer       = flag.String("mongo.server", "localhost:27017", "Address of the Mongo server to monitor")
	enabledCollectors = flag.String("collectors.enabled", "connection,operations,network,lock,cursor", "Comma-separated list of collectors to use.")

	collectorLabelNames = []string{"collector", "result"}

	scrapeDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: collector.Namespace,
			Subsystem: subsystem,
			Name:      "scrape_duration_seconds",
			Help:      "mongo_exporter: Duration of a scrape job.",
		},
		collectorLabelNames,
	)
)

type MongoCollector struct {
	info       *mgo.DialInfo
	collectors map[string]collector.Collector
}

// Implements Collector.
func (m MongoCollector) Describe(ch chan<- *prometheus.Desc) {
	scrapeDurations.Describe(ch)
}

// Implements Collector.
func (m MongoCollector) Collect(ch chan<- prometheus.Metric) {
	session, err := mgo.DialWithInfo(m.info)
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	defer session.Close()
	cmd := &bson.M{
		"serverStatus": 1,
	}

	status := &commgo.ServerStatus{}

	if err := session.DB("local").Run(&cmd, &status); err != nil {
		log.Errorf("%v", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(m.collectors))
	for name, c := range m.collectors {
		go func(name string, c collector.Collector, session *mgo.Session) {
			Execute(name, c, ch, status)
			wg.Done()
		}(name, c, session)
	}
	wg.Wait()
	scrapeDurations.Collect(ch)
}

func Execute(name string, c collector.Collector, ch chan<- prometheus.Metric, status *commgo.ServerStatus) {
	begin := time.Now()
	err := c.Update(ch, status)
	duration := time.Since(begin)
	var result string

	if err != nil {
		log.Infof("ERROR: %s failed after %fs: %s", name, duration.Seconds(), err)
		result = "error"
	} else {
		log.Infof("OK: %s success after %fs.", name, duration.Seconds())
		result = "success"
	}
	scrapeDurations.WithLabelValues(name, result).Observe(duration.Seconds())
}

func loadCollectors() (map[string]collector.Collector, error) {
	collectors := map[string]collector.Collector{}
	for _, name := range strings.Split(*enabledCollectors, ",") {
		fn, ok := collector.Factories[name]
		if !ok {
			return nil, fmt.Errorf("collector '%s' not available", name)
		}
		c, err := fn()
		if err != nil {
			return nil, err
		}
		collectors[name] = c
	}
	return collectors, nil
}

func main() {
	flag.Parse()

	collectors, err := loadCollectors()
	if err != nil {
		log.Fatalf("Couldn't load collectors: %s", err)
	}

	mongoCollector := MongoCollector{
		collectors: collectors,
		info: &mgo.DialInfo{
			Addrs:  []string{*mongoServer},
			Direct: true,
		}}

	prometheus.MustRegister(mongoCollector)

	log.Infof("Starting Server: %s", *listenAddress)
	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Mongo Exporter</title></head>
             <body>
             <h1>Mongo Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
