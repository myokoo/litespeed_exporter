package collector

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/myokoo/litespeed_exporter/pkg/rtreport"
)

var (
	vhostLabels = []string{"vhost"}
	vName       = "virtual_host"
)

type virtualHost struct{}

func (v virtualHost) scrape(ch chan<- prometheus.Metric, report *rtreport.LiteSpeedReport) {
	for vhost, valueMap := range report.VirtualHostReport {
		for key, value := range valueMap {
			switch key {
			case rtreport.VHostReportKeyProcessing:
				ch <- newMetric(
					namespace, vName, "running_processe",
					"The number of running process by vhost.",
					vhostLabels, prometheus.GaugeValue, value, vhost,
				)
			case rtreport.VhostReportKeyReqPerSec:
				ch <- newMetric(
					namespace, vName, "requests_per_sec",
					"The total requests per sec by vhost.",
					vhostLabels, prometheus.GaugeValue, value, vhost,
				)
			case rtreport.VHostReportKeyReqTotal:
				ch <- newMetric(
					namespace, vName, "requests_total",
					"The total requests by vhost.",
					vhostLabels, prometheus.GaugeValue, value, vhost,
				)
			case rtreport.VHostReportKeyStaticHits:
				ch <- newMetric(
					namespace, vName, "hists_total",
					"The number of static requests by vhost.",
					vhostLabels, prometheus.GaugeValue, value, vhost,
				)
			case rtreport.VHostReportKeyPubCacheHits:
				ch <- newMetric(
					namespace, vName, "public_cache_hists_total",
					"The number of public cache hits by vhost.",
					vhostLabels, prometheus.GaugeValue, value, vhost,
				)
			case rtreport.VHostReportKeyPteCacheHits:
				ch <- newMetric(
					namespace, vName, "private_cache_hists_total",
					"The number of private cache hits by vhost.",
					vhostLabels, prometheus.GaugeValue, value, vhost,
				)
			}
		}
	}
}
