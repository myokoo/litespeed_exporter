module github.com/myokoo/litespeed_exporter

go 1.16

require (
	github.com/alecthomas/units v0.0.0-20210208195552-ff826a37aa15 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.5
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.29.0
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/prometheus/common v0.29.0 => github.com/prometheus/common v0.26.0
