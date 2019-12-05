package rtreport

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type LineParser interface {
	parse(report *LiteSpeedReport)
}

func NewLineParser(lineTxt string) LineParser {
	var v LineParser
	switch {
	case strings.HasPrefix(lineTxt, "VERSION:"):
		v = versionLine(lineTxt)
	case strings.HasPrefix(lineTxt, "UPTIME:"):
		v = uptimeLine(lineTxt)
	case strings.HasPrefix(lineTxt, "BPS_IN:"):
		v = networkLine(lineTxt)
	case strings.HasPrefix(lineTxt, "MAXCONN:"):
		v = connectionLine(lineTxt)
	case strings.HasPrefix(lineTxt, "REQ_RATE"):
		v = virtualHostLine(lineTxt)
	case strings.HasPrefix(lineTxt, "EXTAPP"):
		v = extAppLine(lineTxt)
	default:
		v = ignoreLine(lineTxt)
	}
	return v
}

type versionLine string

// VERSION: LiteSpeed Web Server/Enterprise/x.x.x to x.x.x
func (v versionLine) parse(report *LiteSpeedReport) {
	lineText := string(v)
	if len(lineText) < 10 {
		report.error = createTooShortParseLineError(lineText)
		return
	}
	s := strings.Split(lineText, "/")
	report.Version = s[len(s)-1]
}

type uptimeLine string

// parse UPTIME: xx:xx:xx
func (u uptimeLine) parse(report *LiteSpeedReport) {
	lineText := string(u)
	if len(lineText) < 16 {
		report.error = createTooShortParseLineError(lineText)
		return
	}
	v := strings.Split(lineText[8:], ":")
	if len(v) != 3 {
		report.error = errors.New(fmt.Sprintf("Doesn't match split count. string: %s", lineText[8:]))
		return
	}
	h, _ := strconv.ParseUint(v[0], 10, 64)
	m, _ := strconv.ParseUint(v[1], 10, 64)
	s, _ := strconv.ParseUint(v[2], 10, 64)
	report.Uptime = float64((h * 60 * 60) + (m * 60) + s)
	return
}

type networkLine string

// parse BPS_IN: x, BPS_OUT: x, SSL_BPS_IN: x, SSL_BPS_OUT: x
func (n networkLine) parse(report *LiteSpeedReport) {
	report.NetworkReport, report.error = convertStringToMap(string(n))
}

type connectionLine string

// parse MAXCONN: 10000, MAXSSL_CONN: 5000, PLAINCONN: 1, AVAILCONN: 9999, IDLECONN: 0, SSLCONN: 0, AVAILSSL: 5000
func (c connectionLine) parse(report *LiteSpeedReport) {
	report.ConnectionReport, report.error = convertStringToMap(string(c))
}

type virtualHostLine string

// parse REQ_RATE [xxxx]: REQ_PROCESSING: 1, REQ_PER_SEC: 0.1, TOT_REQS: 152, PUB_CACHE_HITS_PER_SEC: 0.0, TOTAL_PUB_CACHE_HITS: 0, PRIVATE_CACHE_HITS_PER_SEC: 0.0, TOTAL_PRIVATE_CACHE_HITS: 0, STATIC_HITS_PER_SEC: 0.0, TOTAL_STATIC_HITS: 47
func (r virtualHostLine) parse(report *LiteSpeedReport) {
	lineText := string(r)

	// pick up vhostName
	s := pickUpStringName(lineText)
	if len(s) < 1 {
		report.error = errors.New(fmt.Sprintf("Cann't Parse VirtualHostName. string: %s", lineText))
		return
	}
	vhName := s[0]
	if vhName == "" {
		vhName = "Server"
	}

	i := strings.Index(lineText, ":")
	report.VirtualHostReport[vhName], report.error = convertStringToMap(lineText[i+1:])
}

type extAppLine string

// parse EXTAPP [xxxx] [xxxx] [xxxx]: CMAXCONN: 1000, EMAXCONN: 1000, POOL_SIZE: 1, INUSE_CONN: 1, IDLE_CONN: 0, WAITQUE_DEPTH: 0, REQ_PER_SEC: 0.0, TOT_REQS: 0
func (e extAppLine) parse(report *LiteSpeedReport) {
	lineText := string(e)

	// pick up ExtAppType, vhostName, ExtAppName
	s := pickUpStringName(lineText)
	if len(s) < 3 {
		report.error = errors.New(fmt.Sprintf("Cann't Parse ExtAppType, VirtualHostName, ExtAppName. string: %s", lineText))
		return
	}
	vhostName := s[1]
	if vhostName == "" {
		vhostName = "Server"
	}
	i := strings.Index(lineText, ":")
	var m map[string]float64
	m, report.error = convertStringToMap(lineText[i+1:])

	if _, exist := report.ExtAppReports[s[0]]; !exist {
		report.ExtAppReports[s[0]] = map[string]map[string]map[string]float64{vhostName: {s[2]: m}}
	} else if _, exit := report.ExtAppReports[s[0]][vhostName]; !exit {
		report.ExtAppReports[s[0]][vhostName] = map[string]map[string]float64{s[2]: m}
	} else {
		report.ExtAppReports[s[0]][vhostName][s[2]] = m
	}
}

type ignoreLine string

// does not parse.
func (i ignoreLine) parse(report *LiteSpeedReport) {
}

// convert "xxxx: 1234, oooo: 4321" strings to map[string]float64{"xxxx":1234, "oooo":4321}
func convertStringToMap(lineText string) (map[string]float64, error) {
	m := make(map[string]float64)
	keyValues := strings.Split(lineText, ",")

	for _, keyValue := range keyValues {
		s := strings.Split(keyValue, ":")
		if len(s) < 2 {
			return nil, errors.New(fmt.Sprintf("Cann't split key value. string: %s", keyValue))
		}
		var err error
		if m[strings.TrimSpace(s[0])], err = strconv.ParseFloat(strings.TrimSpace(s[1]), 64); err != nil {
			return nil, errors.New(fmt.Sprintf("Cann't convert value string to float64. string: %s", strings.TrimSpace(s[1])))
		}
	}
	return m, nil
}

// pick up "[]string{"oooo", "oooo"}" from "XXXX [oooo] [oooo]: xxxxx"
func pickUpStringName(lineText string) []string {
	var s []string
	for i := strings.Count(lineText, "["); i > 0; i-- {
		startIndex := strings.Index(lineText, "[")
		endIndex := strings.Index(lineText, "]")
		s = append(s, strings.TrimSpace(lineText[startIndex+1:endIndex]))
		lineText = lineText[endIndex+1:]
	}
	return s
}

// create parse line too short error.
func createTooShortParseLineError(s string) error {
	return errors.New(fmt.Sprintf("Parse line too short. string: %s", s))
}
