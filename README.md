# litespeed_exporter
Prometheus exporter for LiteSpeed server metrics.

## Building and running
### Compatibility
- Go 1.14+

### Build

```bash
make build
or
make cross_build
```

## usage

```bash
usage: litespeed_exporter [<flags>]

Flags:
  -h, --help              Show context-sensitive help (also try --help-long and --help-man).
      --web.listen-address=":9104"
                          Listen address for web interface and telemetry.
      --web.telemetry-path="/metrics"
                          URL path under which to expose metrics.
      --lsws.report-path="/tmp/lshttpd"
                          Filesystem path under which exist lsws real-time statistics reports.
      --log.level="info"  Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]
      --log.format="logger:stderr"
                          Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true"
      --version           Show application version.

```

## author
@myokoo

