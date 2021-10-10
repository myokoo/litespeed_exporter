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
				"The maximum http connections value of server.",
				connectionLabel, prometheus.GaugeValue, value, "http",
			)
		case rtreport.ConnectionReportKeyMaxConnSsl:
			ch <- newMetric(
				namespace, cName, "max",
				"The maximum https connections value of server.",
				connectionLabel, prometheus.GaugeValue, value, "https",
			)
		case rtreport.ConnectionReportKeyIdleConn:
			ch <- newMetric(
				namespace, cName, "idle",
				"The current idle connections value of server.",
				nil, prometheus.GaugeValue, value,
			)
		case rtreport.ConnectionReportKeyUsedConn:
			ch <- newMetric(
				namespace, cName, "used",
				"The current number of used http connections to server.",
				connectionLabel, prometheus.GaugeValue, value, "http",
			)
		case rtreport.ConnectionReportKeyUsedConnSsl:
			ch <- newMetric(
				namespace, cName, "used",
				"The current number of used https connections to server.",
				connectionLabel, prometheus.GaugeValue, value, "https",
			)
		}
	}
}
