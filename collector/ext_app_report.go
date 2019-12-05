package collector

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/myokoo/litespeed_exporter/pkg/rtreport"
)

var (
	extAppLabels = []string{"type", "vhost", "extapp_name"}
	eName        = "external_application"
)

type extApp struct{}

func (e extApp) scrape(ch chan<- prometheus.Metric, report *rtreport.LiteSpeedReport) {
	for typeName, vhostMap := range report.ExtAppReports {
		for vhost, extAppMap := range vhostMap {
			for extAppName, valueMap := range extAppMap {
				for key, value := range valueMap {
					switch key {
					case rtreport.ExtAppKeyMaxConn:
						ch <- newMetric(
							namespace, eName, "max_connections",
							"The max connection value of external application.",
							extAppLabels, prometheus.GaugeValue, value, typeName, vhost, extAppName,
						)
					case rtreport.ExtAppKeyEffectiveMaxConn:
						ch <- newMetric(
							namespace, eName, "effective_max_connections",
							"The max effective connection value of external application.",
							extAppLabels, prometheus.GaugeValue, value, typeName, vhost, extAppName,
						)
					case rtreport.ExtAppKeyPoolSize:
						ch <- newMetric(
							namespace, eName, "pool_size",
							"The number of pool size by external application.",
							extAppLabels, prometheus.GaugeValue, value, typeName, vhost, extAppName,
						)
					case rtreport.ExtAppKeyInUseConn:
						ch <- newMetric(
							namespace, eName, "connection_used",
							"The number of using connections by external application.",
							extAppLabels, prometheus.GaugeValue, value, typeName, vhost, extAppName,
						)
					case rtreport.ExtAppKeyIdleConn:
						ch <- newMetric(
							namespace, eName, "connection_idles",
							"The number of idle connections by external application.",
							extAppLabels, prometheus.GaugeValue, value, typeName, vhost, extAppName,
						)
					case rtreport.ExtAppKeyWaitQueue:
						ch <- newMetric(
							namespace, eName, "wait_queues",
							"The number of wait queues by external application.",
							extAppLabels, prometheus.GaugeValue, value, typeName, vhost, extAppName,
						)
					case rtreport.ExtAppKeyReqTotal:
						ch <- newMetric(
							namespace, eName, "requests_total",
							"The total requests by external application.",
							extAppLabels, prometheus.GaugeValue, value, typeName, vhost, extAppName,
						)
					}
				}
			}
		}
	}
}
