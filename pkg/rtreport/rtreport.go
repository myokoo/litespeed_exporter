package rtreport

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Constants
const (
	DefaultReportPath    = "/tmp/lshttpd"
	reportFileNamePrefix = ".rtreport"
)

// MapKey
const (
	NetworkReportKeyBpsIn          = "BPS_IN"
	NetworkReportKeyBpsOut         = "BPS_OUT"
	NetworkReportKeySslBpsIn       = "SSL_BPS_IN"
	NetworkReportKeySslBpsOut      = "SSL_BPS_OUT"
	ConnectionReportKeyMaxConn     = "MAXCONN"
	ConnectionReportKeyMaxConnSsl  = "MAXSSL_CONN"
	ConnectionReportKeyUsedConn    = "PLAINCONN"
	ConnectionReportKeyIdleConn    = "IDLECONN"
	ConnectionReportKeyUsedConnSsl = "SSLCONN"
	RequestReportKeyProcessing     = "REQ_PROCESSING"
	RequestReportKeyReqTotal       = "TOT_REQS"
	RequestReportKeyPubCacheHits   = "TOTAL_PUB_CACHE_HITS"
	RequestReportKeyPteCacheHits   = "TOTAL_PRIVATE_CACHE_HITS"
	RequestReportKeyStaticHits     = "TOTAL_STATIC_HITS"
	ExtAppKeyMaxConn               = "CMAXCONN"
	ExtAppKeyEffectiveMaxConn      = "EMAXCONN"
	ExtAppKeyPoolSize              = "POOL_SIZE"
	ExtAppKeyInUseConn             = "INUSE_CONN"
	ExtAppKeyIdleConn              = "IDLE_CONN"
	ExtAppKeyWaitQueue             = "WAITQUE_DEPTH"
	ExtAppKeyReqTotal              = "TOT_REQS"
)

// LiteSpeedReport
type LiteSpeedReport struct {
	error            error
	Version          string
	Uptime           float64
	NetworkReport    map[string]float64
	ConnectionReport map[string]float64
	RequestReports   map[string]map[string]float64
	ExtAppReports    map[string]map[string]map[string]map[string]float64
}

// New return a new instance of real time report and error.
func New(path string) (*LiteSpeedReport, error) {
	reportFiles, err := searchReportFiles(path)
	if err != nil {
		return nil, err
	}

	counter := len(reportFiles)
	ch := make(chan *LiteSpeedReport, counter)
	defer close(ch)
	done := make(chan interface{})
	defer close(done)

	loadReportFiles(done, ch, reportFiles)
	sumReportData(done, ch, counter)

	r := <-ch
	return r, r.error
}

// Search Real TIme Report Files.
func searchReportFiles(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var reportFiles []string
	for _, file := range files {
		if file.IsDir() || !strings.HasPrefix(file.Name(), reportFileNamePrefix) {
			continue
		}
		reportFiles = append(reportFiles, filepath.Join(path, file.Name()))
	}
	return reportFiles, nil
}

func loadReportFiles(done <-chan interface{}, ch chan<- *LiteSpeedReport, reportFiles []string) {
	for _, reportFile := range reportFiles {
		go func(filePath string) {
			select {
			case <-done:
				return
			case ch <- load(filePath):
			}
		}(reportFile)
	}
}

func load(filePath string) *LiteSpeedReport {
	fp, err := os.Open(filePath)
	if err != nil {
		return &LiteSpeedReport{error: err}
	}
	defer fp.Close()

	v := &LiteSpeedReport{
		NetworkReport:    make(map[string]float64),
		ConnectionReport: make(map[string]float64),
		RequestReports:   make(map[string]map[string]float64),
		ExtAppReports:    make(map[string]map[string]map[string]map[string]float64),
	}
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		NewLineParser(scanner.Text()).parse(v)
		if v.error != nil {
			break
		}
	}
	return v
}

func sumReportData(done <-chan interface{}, ch chan *LiteSpeedReport, counter int) {
	for counter > 1 {
		report1 := <-ch
		counter--
		report2 := <-ch
		// be offset. counter-- ; counter++;

		go func(a, b *LiteSpeedReport) {
			select {
			case <-done:
				return
			case ch <- sum(a, b):
			}
		}(report1, report2)
	}
}

func sum(a, b *LiteSpeedReport) *LiteSpeedReport {
	// if error exist. return only error.
	if a.error != nil || b.error != nil {
		v := &LiteSpeedReport{}
		if a.error == nil {
			v.error = b.error
		} else if b.error == nil {
			v.error = a.error
		} else {
			v.error = errors.New(fmt.Sprintf("%s\n--------------%s", a.error.Error(), b.error.Error()))
		}
		return v
	}
	// merge map value.
	margeSingleMap(a.NetworkReport, b.NetworkReport)
	margeSingleMap(a.ConnectionReport, b.ConnectionReport)
	margeDoubleMap(a.RequestReports, b.RequestReports)
	margeQuadrupleMap(a.ExtAppReports, b.ExtAppReports)
	return a
}
