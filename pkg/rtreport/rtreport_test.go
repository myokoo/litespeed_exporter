package rtreport

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_sum(t *testing.T) {
	type args struct {
		a *LiteSpeedReport
		b *LiteSpeedReport
	}
	tests := []struct {
		name string
		args args
		want *LiteSpeedReport
	}{
		{
			name: "ok",
			args: args{
				a: &LiteSpeedReport{
					Uptime:           123,
					Version:          5.4,
					NetworkReport:    map[string]float64{"BPS_IN": 123, "BPS_OUT": 713819, "SSL_BPS_IN": 136, "SSL_BPS_OUT": 891290},
					ConnectionReport: map[string]float64{"MAXCONN": 10000, "MAXSSL_CONN": 5000, "PLAINCONN": 2331, "AVAILCONN": 7669, "IDLECONN": 0, "SSLCONN": 5, "AVAILSSL": 4995},
					RequestReports: map[string]map[string]float64{
						"Server": {"REQ_PROCESSING": 3, "REQ_PER_SEC": 3.5, "TOT_REQS": 1533, "PUB_CACHE_HITS_PER_SEC": 0.0, "TOTAL_PUB_CACHE_HITS": 0, "PRIVATE_CACHE_HITS_PER_SEC": 1.1,
							"TOTAL_PRIVATE_CACHE_HITS": 123, "STATIC_HITS_PER_SEC": 4.4, "TOTAL_STATIC_HITS": 49},
						"hoge.com": {"REQ_PROCESSING": 1, "REQ_PER_SEC": 1.5, "TOT_REQS": 133, "PUB_CACHE_HITS_PER_SEC": 2.1,
							"TOTAL_PUB_CACHE_HITS": 345, "PRIVATE_CACHE_HITS_PER_SEC": 4.3, "TOTAL_PRIVATE_CACHE_HITS": 345, "STATIC_HITS_PER_SEC": 5.5, "TOTAL_STATIC_HITS": 813},
					},
					// EXTAPP [xxxx] [xxxx] [xxxx]: CMAXCONN: 1000, EMAXCONN: 1000, POOL_SIZE: 1, INUSE_CONN: 1, IDLE_CONN: 0, WAITQUE_DEPTH: 0, REQ_PER_SEC: 0.0, TOT_REQS: 0
					ExtAppReports: map[string]map[string]map[string]map[string]float64{
						"LSAPI": {"hoge.com": {"hoge.com_php7.3": {"CMAXCONN": 1000, "EMAXCONN": 1000, "POOL_SIZE": 1,
							"INUSE_CONN": 1, "IDLE_CONN": 0, "WAITQUE_DEPTH": 0, "REQ_PER_SEC": 0.0, "TOT_REQS": 0}}},
					},
				},
				b: &LiteSpeedReport{
					Uptime:           123,
					Version:          5.4,
					NetworkReport:    map[string]float64{"BPS_IN": 21213, "BPS_OUT": 343819, "SSL_BPS_IN": 123363, "SSL_BPS_OUT": 913290},
					ConnectionReport: map[string]float64{"MAXCONN": 10000, "MAXSSL_CONN": 5000, "PLAINCONN": 1000, "AVAILCONN": 9000, "IDLECONN": 1, "SSLCONN": 100, "AVAILSSL": 4900},
					RequestReports: map[string]map[string]float64{
						"Server": {"REQ_PROCESSING": 5, "REQ_PER_SEC": 6.6, "TOT_REQS": 903, "PUB_CACHE_HITS_PER_SEC": 3.8, "TOTAL_PUB_CACHE_HITS": 1100, "PRIVATE_CACHE_HITS_PER_SEC": 5.3,
							"TOTAL_PRIVATE_CACHE_HITS": 9393, "STATIC_HITS_PER_SEC": 7.9, "TOTAL_STATIC_HITS": 3939},
					},
					ExtAppReports: map[string]map[string]map[string]map[string]float64{
						"LSAPI": {"hoge.com": {"hoge.com_php7.3": {"CMAXCONN": 1000, "EMAXCONN": 1000, "POOL_SIZE": 2,
							"INUSE_CONN": 4, "IDLE_CONN": 3, "WAITQUE_DEPTH": 2, "REQ_PER_SEC": 1.2, "TOT_REQS": 3}}},
						"CGI": {"Server": {"lscgid": {"CMAXCONN": 1000, "EMAXCONN": 1000, "POOL_SIZE": 1,
							"INUSE_CONN": 1, "IDLE_CONN": 0, "WAITQUE_DEPTH": 0, "REQ_PER_SEC": 0.0, "TOT_REQS": 0}}},
					},
				},
			},
			want: &LiteSpeedReport{
				Uptime:           123,
				Version:          5.4,
				NetworkReport:    map[string]float64{"BPS_IN": 21336, "BPS_OUT": 1057638, "SSL_BPS_IN": 123499, "SSL_BPS_OUT": 1804580},
				ConnectionReport: map[string]float64{"MAXCONN": 20000, "MAXSSL_CONN": 10000, "PLAINCONN": 3331, "AVAILCONN": 16669, "IDLECONN": 1, "SSLCONN": 105, "AVAILSSL": 9895},
				RequestReports: map[string]map[string]float64{
					"Server": {"REQ_PROCESSING": 8, "REQ_PER_SEC": 10.1, "TOT_REQS": 2436, "PUB_CACHE_HITS_PER_SEC": 3.8, "TOTAL_PUB_CACHE_HITS": 1100, "PRIVATE_CACHE_HITS_PER_SEC": 6.4,
						"TOTAL_PRIVATE_CACHE_HITS": 9516, "STATIC_HITS_PER_SEC": 12.3, "TOTAL_STATIC_HITS": 3988},
					"hoge.com": {"REQ_PROCESSING": 1, "REQ_PER_SEC": 1.5, "TOT_REQS": 133, "PUB_CACHE_HITS_PER_SEC": 2.1,
						"TOTAL_PUB_CACHE_HITS": 345, "PRIVATE_CACHE_HITS_PER_SEC": 4.3, "TOTAL_PRIVATE_CACHE_HITS": 345, "STATIC_HITS_PER_SEC": 5.5, "TOTAL_STATIC_HITS": 813},
				},
				ExtAppReports: map[string]map[string]map[string]map[string]float64{
					"LSAPI": {"hoge.com": {"hoge.com_php7.3": {"CMAXCONN": 2000, "EMAXCONN": 2000, "POOL_SIZE": 3,
						"INUSE_CONN": 5, "IDLE_CONN": 3, "WAITQUE_DEPTH": 2, "REQ_PER_SEC": 1.2, "TOT_REQS": 3}}},
					"CGI": {"Server": {"lscgid": {"CMAXCONN": 1000, "EMAXCONN": 1000, "POOL_SIZE": 1,
						"INUSE_CONN": 1, "IDLE_CONN": 0, "WAITQUE_DEPTH": 0, "REQ_PER_SEC": 0.0, "TOT_REQS": 0}}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sum(tt.args.a, tt.args.b)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_load(t *testing.T) {
	tests := []struct {
		name string
		args string
		want *LiteSpeedReport
	}{
		{
			name: "ok",
			args: "../test/data/load/.rtreport",
			want: &LiteSpeedReport{
				Version:          5.4,
				Uptime:           56070,
				NetworkReport:    map[string]float64{"BPS_IN": 1, "BPS_OUT": 2, "SSL_BPS_IN": 3, "SSL_BPS_OUT": 4},
				ConnectionReport: map[string]float64{"MAXCONN": 10000, "MAXSSL_CONN": 5000, "PLAINCONN": 0, "AVAILCONN": 10000, "IDLECONN": 0, "SSLCONN": 0, "AVAILSSL": 5000},
				RequestReports: map[string]map[string]float64{
					"Server": {"REQ_PROCESSING": 0, "REQ_PER_SEC": 0.1, "TOT_REQS": 448, "PUB_CACHE_HITS_PER_SEC": 0.0, "TOTAL_PUB_CACHE_HITS": 0, "PRIVATE_CACHE_HITS_PER_SEC": 0.0,
						"TOTAL_PRIVATE_CACHE_HITS": 0, "STATIC_HITS_PER_SEC": 0.1, "TOTAL_STATIC_HITS": 133},
					"hoge.jp": {"REQ_PROCESSING": 3, "REQ_PER_SEC": 2.1, "TOT_REQS": 121, "PUB_CACHE_HITS_PER_SEC": 4.0,
						"TOTAL_PUB_CACHE_HITS": 345, "PRIVATE_CACHE_HITS_PER_SEC": 4.3, "TOTAL_PRIVATE_CACHE_HITS": 345, "STATIC_HITS_PER_SEC": 5.5, "TOTAL_STATIC_HITS": 813},
				},
				ExtAppReports: make(map[string]map[string]map[string]map[string]float64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := load(tt.args)
			if !cmp.Equal(got, tt.want, cmp.AllowUnexported(LiteSpeedReport{})) {
				t.Errorf("load() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    *LiteSpeedReport
		wantErr bool
	}{
		{
			name: "ok",
			args: "../test/data/new",
			want: &LiteSpeedReport{
				Version:          5.4,
				Uptime:           56070,
				NetworkReport:    map[string]float64{"BPS_IN": 2, "BPS_OUT": 4, "SSL_BPS_IN": 6, "SSL_BPS_OUT": 8},
				ConnectionReport: map[string]float64{"MAXCONN": 20000, "MAXSSL_CONN": 10000, "PLAINCONN": 0, "AVAILCONN": 20000, "IDLECONN": 0, "SSLCONN": 0, "AVAILSSL": 10000},
				RequestReports: map[string]map[string]float64{
					"Server": {"REQ_PROCESSING": 0, "REQ_PER_SEC": 0.2, "TOT_REQS": 896, "PUB_CACHE_HITS_PER_SEC": 0.0, "TOTAL_PUB_CACHE_HITS": 0, "PRIVATE_CACHE_HITS_PER_SEC": 0.0,
						"TOTAL_PRIVATE_CACHE_HITS": 0, "STATIC_HITS_PER_SEC": 0.2, "TOTAL_STATIC_HITS": 266},
					"hoge.jp": {"REQ_PROCESSING": 6, "REQ_PER_SEC": 4.2, "TOT_REQS": 242, "PUB_CACHE_HITS_PER_SEC": 8.0,
						"TOTAL_PUB_CACHE_HITS": 690, "PRIVATE_CACHE_HITS_PER_SEC": 8.6, "TOTAL_PRIVATE_CACHE_HITS": 690, "STATIC_HITS_PER_SEC": 11.0, "TOTAL_STATIC_HITS": 1626},
				},
				ExtAppReports: map[string]map[string]map[string]map[string]float64{
					"LSAPI": {"hoge.jp": {"hoge.jp_php73": {"CMAXCONN": 1000, "EMAXCONN": 1000, "POOL_SIZE": 1,
						"INUSE_CONN": 1, "IDLE_CONN": 0, "WAITQUE_DEPTH": 0, "REQ_PER_SEC": 0.0, "TOT_REQS": 0}}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmp.AllowUnexported(LiteSpeedReport{})) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}
