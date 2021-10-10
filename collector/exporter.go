package collector

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/myokoo/litespeed_exporter/pkg/rtreport"
)

const (
	namespace = "litespeed" // For Prometheus metrics.
)

var (
	errorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Whether the realtime report could be read", nil, nil,
	)
	upteimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "uptime_seconds_total"),
		"Current uptime in seconds.", nil, nil,
	)
)

type Exporter struct {
	mutex      sync.Mutex
	reportPath *string
	scrapers   []Scraper
}

func New(path *string) *Exporter {
	return &Exporter{
		reportPath: path,
		scrapers: []Scraper{
			connection{},
			network{},
			virtualHost{},
			extApp{},
		},
	}
}

// Describe implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- upteimeDesc
}

// Collect implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	report, err := rtreport.New(*e.reportPath)
	if err != nil {
		ch <- metricsIsLitespeedUp(float64(0))
		return
	}
	ch <- metricsIsLitespeedUp(float64(1))
	ch <- prometheus.MustNewConstMetric(upteimeDesc, prometheus.CounterValue, report.Uptime)

	for _, scraper := range e.scrapers {
		scraper.scrape(ch, report)
	}
}

func metricsIsLitespeedUp(i float64) prometheus.Metric {
	return prometheus.MustNewConstMetric(errorDesc, prometheus.GaugeValue, i)
}

func newMetric(namespace, subsystem, name, help string, label []string, metricType prometheus.ValueType, value float64, labelValues ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, name), help, label, nil), metricType, value, labelValues...)
}
