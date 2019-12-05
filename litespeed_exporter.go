package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/myokoo/litespeed_exporter/collector"
	"github.com/myokoo/litespeed_exporter/pkg/rtreport"
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
)

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

	exporter := collector.New(reportPath)
	prometheus.MustRegister(exporter)
	prometheus.MustRegister(version.NewCollector("litespeed_exporter"))

	http.Handle(*metricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write(landingPage) })
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
