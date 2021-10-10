package collector

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/myokoo/litespeed_exporter/pkg/rtreport"
)

// Scraper is a minimal interface that allows you to add new prometheus metrics to litespeed_exporter.
type Scraper interface {
	scrape(ch chan<- prometheus.Metric, report *rtreport.LiteSpeedReport)
}
