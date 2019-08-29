package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/myokoo/litespeed_exporter/pkg/rtreport"
)

const (
	namespace = "litespeed" // For Prometheus metrics.
)

var (
	listenAddress = kingpin.Flag(
		"web.listen-address",
		"Address to listen on for web interface and telemetry.",
	).Default(":9104").String()
	metricPath = kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
	).Default("/metrics").String()
	reportPath = kingpin.Flag(
		"lsws.report-path",
		"Path under which to exist lsws real-time statistics report.",
	).Default(rtreport.DefaultReportPath).String()
	insecure = kingpin.Flag(
		"insecure",
		"Ignore server certificate if using https.",
	).Default("false").Bool()

	extAppLabels = []string{"type", "vhost", "extapp_name"}
	vhostLabel   = []string{"vhost"}
	schemeLabel  = []string{"scheme"}
)

type Exporter struct {
	mutex sync.Mutex
	path  string

	up                            *prometheus.Desc
	uptime                        prometheus.Gauge
	networkThroughput             *prometheus.GaugeVec
	serverConnectionsMax          *prometheus.GaugeVec
	serverConnectionsUsed         *prometheus.GaugeVec
	serverConnectionsIdle         prometheus.Gauge
	vhostRunningProcesses         *prometheus.GaugeVec
	vhostRequestsTotal            *prometheus.GaugeVec
	vhostStaticHitsTotal          *prometheus.GaugeVec
	vhostPrivateCacheHitsTotal    *prometheus.GaugeVec
	vhostPublicCacheHitsTotal     *prometheus.GaugeVec
	extAppMaxConnections          *prometheus.GaugeVec
	extAppEffectiveMaxConnections *prometheus.GaugeVec
	extAppRequestsTotal           *prometheus.GaugeVec
	extAppPoolSize                *prometheus.GaugeVec
	extAppConnectionUsed          *prometheus.GaugeVec
	extAppConnectionIdle          *prometheus.GaugeVec
	extAppConnectionWaitQueue     *prometheus.GaugeVec
}

func NewExporter(path string) *Exporter {
	return &Exporter{
		path: path,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Could the realtime report be readed",
			nil,
			nil,
		),
		uptime: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "uptime_seconds_total",
				Help:      "Current uptime in seconds (*)",
			},
		),
		networkThroughput: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "network_throughput",
				Help:      "Current network throughput by schema(http or https)",
			},
			[]string{"scheme", "stream"},
		),
		serverConnectionsMax: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "server_connection_max",
				Help:      "The max connection value of server.",
			},
			schemeLabel,
		),
		serverConnectionsUsed: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "server_connections_used",
				Help:      "The number of using connections to server.",
			},
			schemeLabel,
		),
		serverConnectionsIdle: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "server_connections_idle",
				Help:      "The idle connection values of server.",
			},
		),
		vhostRunningProcesses: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "vhost_running_processes",
				Help:      "The number of running process by vhost.",
			},
			vhostLabel,
		),
		vhostRequestsTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "vhost_requests_total",
				Help:      "The total requests by vhost.",
			},
			vhostLabel,
		),
		vhostStaticHitsTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "vhost_static_hists_total",
				Help:      "The number of static requests by vhost.",
			},
			vhostLabel,
		),
		vhostPrivateCacheHitsTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "vhost_private_cache_hists_total",
				Help:      "The number of private cache hits by vhost.",
			},
			vhostLabel,
		),
		vhostPublicCacheHitsTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "vhost_public_cache_hists_total",
				Help:      "The number of public cache hits by vhost.",
			},
			vhostLabel,
		),
		extAppMaxConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "external_application_max_connections",
				Help:      "The max connection value of external application.",
			},
			extAppLabels,
		),
		extAppEffectiveMaxConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "external_application_effective_max_connections",
				Help:      "The max effective connection value of external application.",
			},
			extAppLabels,
		),
		extAppRequestsTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "external_application_requests_total",
				Help:      "The total requests by external application.",
			},
			extAppLabels,
		),
		extAppPoolSize: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "external_application_pool_size",
				Help:      "The number of pool size by external application.",
			},
			extAppLabels,
		),
		extAppConnectionUsed: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "external_application_connection_used",
				Help:      "The number of using connections by external application.",
			},
			extAppLabels,
		),
		extAppConnectionIdle: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "external_application_connection_idles",
				Help:      "The number of idle connections by external application.",
			},
			extAppLabels,
		),
		extAppConnectionWaitQueue: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "external_application_connection_wait_queues",
				Help:      "The number of wait queues by external application.",
			},
			extAppLabels,
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up
	e.uptime.Describe(ch)
	e.networkThroughput.Describe(ch)
	e.serverConnectionsMax.Describe(ch)
	e.serverConnectionsUsed.Describe(ch)
	e.serverConnectionsIdle.Describe(ch)
	e.vhostRunningProcesses.Describe(ch)
	e.vhostRequestsTotal.Describe(ch)
	e.vhostStaticHitsTotal.Describe(ch)
	e.vhostPrivateCacheHitsTotal.Describe(ch)
	e.vhostPublicCacheHitsTotal.Describe(ch)
	e.extAppMaxConnections.Describe(ch)
	e.extAppEffectiveMaxConnections.Describe(ch)
	e.extAppRequestsTotal.Describe(ch)
	e.extAppPoolSize.Describe(ch)
	e.extAppConnectionUsed.Describe(ch)
	e.extAppConnectionIdle.Describe(ch)
	e.extAppConnectionWaitQueue.Describe(ch)
}

func (e *Exporter) networkCollect(ch chan<- prometheus.Metric, report *rtreport.LiteSpeedReport) {
	e.networkThroughput.Reset()

	for key, value := range report.NetworkReport {
		switch key {
		case rtreport.NetworkReportKeyBpsIn:
			e.networkThroughput.WithLabelValues("http", "in").Set(value)
		case rtreport.NetworkReportKeyBpsOut:
			e.networkThroughput.WithLabelValues("http", "out").Set(value)
		case rtreport.NetworkReportKeySslBpsIn:
			e.networkThroughput.WithLabelValues("https", "in").Set(value)
		case rtreport.NetworkReportKeySslBpsOut:
			e.networkThroughput.WithLabelValues("https", "out").Set(value)
		}
	}
	
	e.networkThroughput.Collect(ch)
}

func (e *Exporter) connectionCollect(ch chan<- prometheus.Metric, report *rtreport.LiteSpeedReport) {
	e.serverConnectionsMax.Reset()
	e.serverConnectionsUsed.Reset()

	for key, value := range report.ConnectionReport {
		switch key {
		case rtreport.ConnectionReportKeyMaxConn:
			e.serverConnectionsMax.WithLabelValues("http").Set(value)
		case rtreport.ConnectionReportKeyMaxConnSsl:
			e.serverConnectionsMax.WithLabelValues("https").Set(value)
		case rtreport.ConnectionReportKeyIdleConn:
			e.serverConnectionsIdle.Set(value)
		case rtreport.ConnectionReportKeyUsedConn:
			e.serverConnectionsUsed.WithLabelValues("http").Set(value)
		case rtreport.ConnectionReportKeyUsedConnSsl:
			e.serverConnectionsUsed.WithLabelValues("https").Set(value)
		}
	}

	e.serverConnectionsMax.Collect(ch)
	e.serverConnectionsUsed.Collect(ch)
	e.serverConnectionsIdle.Collect(ch)
}

func (e *Exporter) vhostCollect(ch chan<- prometheus.Metric, report *rtreport.LiteSpeedReport) {
	e.vhostRunningProcesses.Reset()
	e.vhostRequestsTotal.Reset()
	e.vhostStaticHitsTotal.Reset()
	e.vhostPublicCacheHitsTotal.Reset()
	e.vhostPrivateCacheHitsTotal.Reset()

	for vhost, m := range report.RequestReports {
		for key, value := range m {
			switch key {
			case rtreport.RequestReportKeyProcessing:
				e.vhostRunningProcesses.WithLabelValues(vhost).Set(value)
			case rtreport.RequestReportKeyReqTotal:
				e.vhostRequestsTotal.WithLabelValues(vhost).Set(value)
			case rtreport.RequestReportKeyStaticHits:
				e.vhostStaticHitsTotal.WithLabelValues(vhost).Set(value)
			case rtreport.RequestReportKeyPubCacheHits:
				e.vhostPublicCacheHitsTotal.WithLabelValues(vhost).Set(value)
			case rtreport.RequestReportKeyPteCacheHits:
				e.vhostPrivateCacheHitsTotal.WithLabelValues(vhost).Set(value)
			}
		}
	}

	e.vhostRunningProcesses.Collect(ch)
	e.vhostRequestsTotal.Collect(ch)
	e.vhostStaticHitsTotal.Collect(ch)
	e.vhostPublicCacheHitsTotal.Collect(ch)
	e.vhostPrivateCacheHitsTotal.Collect(ch)
}

func (e *Exporter) extAppCollect(ch chan<- prometheus.Metric, report *rtreport.LiteSpeedReport) {
	e.extAppMaxConnections.Reset()
	e.extAppEffectiveMaxConnections.Reset()
	e.extAppPoolSize.Reset()
	e.extAppConnectionUsed.Reset()
	e.extAppConnectionIdle.Reset()
	e.extAppConnectionWaitQueue.Reset()
	e.extAppRequestsTotal.Reset()

	for typeName, m := range report.ExtAppReports {
		for vhost, m2 := range m {
			for extAppName, m3 := range m2 {
				for key, value := range m3 {
					switch key {
					case rtreport.ExtAppKeyMaxConn:
						e.extAppMaxConnections.WithLabelValues(typeName, vhost, extAppName).Set(value)
					case rtreport.ExtAppKeyEffectiveMaxConn:
						e.extAppEffectiveMaxConnections.WithLabelValues(typeName, vhost, extAppName).Set(value)
					case rtreport.ExtAppKeyPoolSize:
						e.extAppPoolSize.WithLabelValues(typeName, vhost, extAppName).Set(value)
					case rtreport.ExtAppKeyInUseConn:
						e.extAppConnectionUsed.WithLabelValues(typeName, vhost, extAppName).Set(value)
					case rtreport.ExtAppKeyIdleConn:
						e.extAppConnectionIdle.WithLabelValues(typeName, vhost, extAppName).Set(value)
					case rtreport.ExtAppKeyWaitQueue:
						e.extAppConnectionWaitQueue.WithLabelValues(typeName, vhost, extAppName).Set(value)
					case rtreport.ExtAppKeyReqTotal:
						e.extAppRequestsTotal.WithLabelValues(typeName, vhost, extAppName).Set(value)
					}
				}
			}
		}
	}

	e.extAppMaxConnections.Collect(ch)
	e.extAppEffectiveMaxConnections.Collect(ch)
	e.extAppPoolSize.Collect(ch)
	e.extAppConnectionUsed.Collect(ch)
	e.extAppConnectionIdle.Collect(ch)
	e.extAppConnectionWaitQueue.Collect(ch)
	e.extAppRequestsTotal.Collect(ch)
}

func (e *Exporter) collect(ch chan<- prometheus.Metric) error {
	report, err := rtreport.New(e.path)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
		return fmt.Errorf("Error read real-time report: %v", err)
	}
	ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 1)
	e.uptime.Set(report.Uptime)
	e.uptime.Collect(ch)
	e.networkCollect(ch, report)
	e.connectionCollect(ch, report)
	e.vhostCollect(ch, report)
	e.extAppCollect(ch, report)
	return nil
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if err := e.collect(ch); err != nil {
		log.Errorf("Error scraping apache: %s", err)
	}
	return
}

func main() {
	// Parse flags.
	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("litespeed_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	// landingPage contains the HTML served at '/'.
	// TODO: Make this nicer and more informative.
	var landingPage = []byte(`<html>
<head><title>LiteSpeed exporter</title></head>
<body>
<h1>LiteSpeed exporter</h1>
<p><a href='` + *metricPath + `'>Metrics</a></p>
</body>
</html>
`)

	log.Infoln("Starting litespeed_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())
	log.Infoln("Listening on", *listenAddress)

	exporter := NewExporter(*reportPath)
	prometheus.MustRegister(exporter)
	prometheus.MustRegister(version.NewCollector("litespeed_exporter"))

	http.Handle(*metricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write(landingPage) })
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
