package collector

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/myokoo/litespeed_exporter/pkg/rtreport"
)

var (
	networkLabel = []string{"scheme", "stream"}
	nName        = "network"
)

type network struct{}

func (n network) scrape(ch chan<- prometheus.Metric, report *rtreport.LiteSpeedReport) {
	for key, value := range report.NetworkReport {
		switch key {
		case rtreport.NetworkReportKeyBpsIn:
			ch <- newMetric(
				namespace, nName, "throughput",
				"Current ingress network throughput (http).",
				networkLabel, prometheus.GaugeValue, value, "http", "in",
			)
		case rtreport.NetworkReportKeyBpsOut:
			ch <- newMetric(
				namespace, nName, "throughput",
				"Current egress network throughput (http).",
				networkLabel, prometheus.GaugeValue, value, "http", "out",
			)
		case rtreport.NetworkReportKeySslBpsIn:
			ch <- newMetric(
				namespace, nName, "throughput",
				"Current ingress network throughput (https).",
				networkLabel, prometheus.GaugeValue, value, "https", "in",
			)
		case rtreport.NetworkReportKeySslBpsOut:
			ch <- newMetric(
				namespace, nName, "throughput",
				"Current egress network throughput (https).",
				networkLabel, prometheus.GaugeValue, value, "https", "out",
			)
		}
	}
}
