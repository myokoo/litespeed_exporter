package collector

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/myokoo/litespeed_exporter/pkg/rtreport"
)

var (
	connectionLabel = []string{"scheme"}
	cName           = "server_connection"
)

type connection struct{}

func (c connection) scrape(ch chan<- prometheus.Metric, report *rtreport.LiteSpeedReport) {
	for key, value := range report.ConnectionReport {
		switch key {
		case rtreport.ConnectionReportKeyMaxConn:
			ch <- newMetric(
				namespace, cName, "max",
				"The max connection value of server.",
				connectionLabel, prometheus.GaugeValue, value, "http",
			)
		case rtreport.ConnectionReportKeyMaxConnSsl:
			ch <- newMetric(
				namespace, cName, "max",
				"The max connection value of server.",
				connectionLabel, prometheus.GaugeValue, value, "https",
			)
		case rtreport.ConnectionReportKeyIdleConn:
			ch <- newMetric(
				namespace, cName, "idle",
				"The idle connection values of server.",
				nil, prometheus.GaugeValue, value,
			)
		case rtreport.ConnectionReportKeyUsedConn:
			ch <- newMetric(
				namespace, cName, "used",
				"The number of using connections to server.",
				connectionLabel, prometheus.GaugeValue, value, "http",
			)
		case rtreport.ConnectionReportKeyUsedConnSsl:
			ch <- newMetric(
				namespace, cName, "used",
				"The number of using connections to server.",
				connectionLabel, prometheus.GaugeValue, value, "https",
			)
		}
	}
}
