package collector

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/myokoo/litespeed_exporter/pkg/rtreport"
)

// Scraper is minimal interface that let's you add new prometheus metrics to litespeed_exporter.
type Scraper interface {
	scrape(ch chan<- prometheus.Metric, report *rtreport.LiteSpeedReport)
}
