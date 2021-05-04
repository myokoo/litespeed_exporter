package rtreport

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewLineParser(t *testing.T) {
	tests := []struct {
		name string
		args string
		want LineParser
	}{
		{
			name: "ok_versionLine",
			args: "VERSION: LiteSpeed Web Server/Enterprise/5.8",
			want: versionLine("VERSION: LiteSpeed Web Server/Enterprise/5.8"),
		},
		{
			name: "ok_uptimeLine",
			args: "UPTIME: 03:02:01",
			want: uptimeLine("UPTIME: 03:02:01"),
		},
		{
			name: "ok_networkLine",
			args: "BPS_IN: 1, BPS_OUT: 2, SSL_BPS_IN: 3, SSL_BPS_OUT: 4",
			want: networkLine("BPS_IN: 1, BPS_OUT: 2, SSL_BPS_IN: 3, SSL_BPS_OUT: 4"),
		},
		{
			name: "ok_connectionLine",
			args: "MAXCONN: 10000, MAXSSL_CONN: 5000, PLAINCONN: 1, AVAILCONN: 9999, IDLECONN: 0, SSLCONN: 0, AVAILSSL: 5000",
			want: connectionLine("MAXCONN: 10000, MAXSSL_CONN: 5000, PLAINCONN: 1, AVAILCONN: 9999, IDLECONN: 0, SSLCONN: 0, AVAILSSL: 5000"),
		},
		{
			name: "ok_requestLine",
			args: "REQ_RATE [hoge.com]: REQ_PROCESSING: 1, REQ_PER_SEC: 0.1, TOT_REQS: 152, PUB_CACHE_HITS_PER_SEC: 0.0, TOTAL_PUB_CACHE_HITS: 0, " +
				"PRIVATE_CACHE_HITS_PER_SEC: 0.0, TOTAL_PRIVATE_CACHE_HITS: 0, STATIC_HITS_PER_SEC: 0.0, TOTAL_STATIC_HITS: 47",
			want: virtualHostLine("REQ_RATE [hoge.com]: REQ_PROCESSING: 1, REQ_PER_SEC: 0.1, TOT_REQS: 152, PUB_CACHE_HITS_PER_SEC: 0.0, " +
				"TOTAL_PUB_CACHE_HITS: 0, PRIVATE_CACHE_HITS_PER_SEC: 0.0, TOTAL_PRIVATE_CACHE_HITS: 0, STATIC_HITS_PER_SEC: 0.0, TOTAL_STATIC_HITS: 47"),
		},
		{
			name: "ok_extAppLine",
			args: "EXTAPP [LSAPI] [hoge.com] [hoge.com_php7.3]: CMAXCONN: 1000, EMAXCONN: 1000, POOL_SIZE: 1, " +
				"INUSE_CONN: 1, IDLE_CONN: 0, WAITQUE_DEPTH: 0, REQ_PER_SEC: 0.0, TOT_REQS: 0",
			want: extAppLine("EXTAPP [LSAPI] [hoge.com] [hoge.com_php7.3]: CMAXCONN: 1000, EMAXCONN: 1000, POOL_SIZE: 1, " +
				"INUSE_CONN: 1, IDLE_CONN: 0, WAITQUE_DEPTH: 0, REQ_PER_SEC: 0.0, TOT_REQS: 0"),
		},
		{
			name: "ok_ignoreLine",
			args: "BLOCKED_IP:",
			want: ignoreLine("BLOCKED_IP:"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLineParser(tt.args); cmp.Equal(got, tt.want) != true {
				t.Errorf("NewLineParser() = %v, want %v, type: %v, %v", got, tt.want, reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
}

func Test_versionLine_parse(t *testing.T) {
	tests := []struct {
		name    string
		v       versionLine
		args    LiteSpeedReport
		want    string
		wantErr bool
	}{
		{
			name:    "ok",
			v:       versionLine("VERSION: LiteSpeed Web Server/Enterprise/5.8.1"),
			args:    LiteSpeedReport{},
			want:    "5.8.1",
			wantErr: false,
		},
		{
			name:    "ng",
			v:       versionLine("5.8"),
			args:    LiteSpeedReport{},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.parse(&tt.args)
			if (tt.args.error != nil) != tt.wantErr {
				t.Errorf("(versionLine)parse() error = %v, wantErr %v", tt.args.error, tt.wantErr)
			}
			if !tt.wantErr && tt.args.Version != tt.want {
				t.Errorf("(versionLine)parse() does not match. got = %v, want = %v", tt.args.Version, tt.want)
			}
		})
	}
}

func Test_uptimeLine_parse(t *testing.T) {
	tests := []struct {
		name    string
		u       uptimeLine
		args    LiteSpeedReport
		want    float64
		wantErr bool
	}{
		{
			name:    "ok",
			u:       uptimeLine("UPTIME: 03:02:01"),
			args:    LiteSpeedReport{},
			want:    10921,
			wantErr: false,
		},
		{
			name:    "ng",
			u:       uptimeLine("03:02:01"),
			args:    LiteSpeedReport{},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.parse(&tt.args)
			if (tt.args.error != nil) != tt.wantErr {
				t.Errorf("(uptimeLine)parse() error = %v, wantErr %v", tt.args.error, tt.wantErr)
			}
			if !tt.wantErr && tt.args.Uptime != tt.want {
				t.Errorf("(uptimeLine)parse() does not match. got = %v, want = %v", tt.args.Uptime, tt.want)
			}
		})
	}
}

func Test_networkLine_parse(t *testing.T) {
	tests := []struct {
		name    string
		n       networkLine
		args    LiteSpeedReport
		want    LiteSpeedReport
		wantErr bool
	}{
		{
			name: "ok",
			n:    networkLine("BPS_IN: 2, BPS_OUT: 1954, SSL_BPS_IN: 5, SSL_BPS_OUT: 3332"),
			args: LiteSpeedReport{},
			want: LiteSpeedReport{
				NetworkReport: map[string]float64{
					"BPS_IN":      2,
					"BPS_OUT":     1954,
					"SSL_BPS_IN":  5,
					"SSL_BPS_OUT": 3332,
				},
			},
			wantErr: false,
		},
		{
			name:    "ng",
			n:       networkLine("BPS_IN: 2 BPS_OUT: 1954"),
			args:    LiteSpeedReport{},
			want:    LiteSpeedReport{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.parse(&tt.args)
			if (tt.args.error != nil) != tt.wantErr {
				t.Errorf("(networkLine)parse() error = %v, wantErr %v", tt.args.error, tt.wantErr)
			}
			if !tt.wantErr && !cmp.Equal(tt.args.NetworkReport, tt.want.NetworkReport) {
				t.Errorf("(networkLine)parse() does not match. got = %v, want = %v", tt.args.NetworkReport, tt.want.NetworkReport)
			}
		})
	}
}

func Test_connectionLine_parse(t *testing.T) {
	tests := []struct {
		name    string
		c       connectionLine
		args    LiteSpeedReport
		want    LiteSpeedReport
		wantErr bool
	}{
		{
			name: "ok",
			c:    connectionLine("MAXCONN: 10000, MAXSSL_CONN: 5000, PLAINCONN: 100, AVAILCONN: 200, IDLECONN: 1, SSLCONN: 2, AVAILSSL: 3"),
			args: LiteSpeedReport{},
			want: LiteSpeedReport{
				ConnectionReport: map[string]float64{
					"MAXCONN":     10000,
					"MAXSSL_CONN": 5000,
					"PLAINCONN":   100,
					"AVAILCONN":   200,
					"IDLECONN":    1,
					"SSLCONN":     2,
					"AVAILSSL":    3,
				},
			},
			wantErr: false,
		},
		{
			name:    "ng",
			c:       connectionLine("MAXCONN: 10000 MAXSSL_CONN: 5000, PLAINCONN: 100, AVAILCONN: 200, IDLECONN: 1, SSLCONN: 2"),
			args:    LiteSpeedReport{},
			want:    LiteSpeedReport{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.parse(&tt.args)
			if (tt.args.error != nil) != tt.wantErr {
				t.Errorf("(connectionLine)parse() error = %v, wantErr %v", tt.args.error, tt.wantErr)
			}
			if !tt.wantErr && !cmp.Equal(tt.args.ConnectionReport, tt.want.ConnectionReport) {
				t.Errorf("(connectionLine)parse() does not match. got = %v, want = %v", tt.args.ConnectionReport, tt.want.ConnectionReport)
			}
		})
	}
}

func Test_ignoreLine_parse(t *testing.T) {
	tests := []struct {
		name    string
		i       ignoreLine
		args    LiteSpeedReport
		want    LiteSpeedReport
		wantErr bool
	}{
		{
			name:    "ok. not store",
			i:       "BLOCKED_IP:",
			args:    LiteSpeedReport{},
			want:    LiteSpeedReport{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.i.parse(&tt.args)
			if (tt.args.error != nil) != tt.wantErr {
				t.Errorf("(ignoreLine)parse() error = %v, wantErr %v", tt.args.error, tt.wantErr)
			}
			if !tt.wantErr && !cmp.Equal(tt.args.ConnectionReport, tt.want.ConnectionReport) {
				t.Errorf("(ignoreLine)parse() does not match. got = %v, want = %v", tt.args, tt.want)
			}
		})
	}
}

func Test_requestLine_parse(t *testing.T) {
	tests := []struct {
		name    string
		r       virtualHostLine
		args    LiteSpeedReport
		want    LiteSpeedReport
		wantErr bool
	}{
		{
			name: "ok_vhost",
			r: virtualHostLine("REQ_RATE [hoge.jp]: REQ_PROCESSING: 1, REQ_PER_SEC: 0.1, TOT_REQS: 2, PUB_CACHE_HITS_PER_SEC: 0.2, TOTAL_PUB_CACHE_HITS: 3, " +
				"PRIVATE_CACHE_HITS_PER_SEC: 0.3, TOTAL_PRIVATE_CACHE_HITS: 4, STATIC_HITS_PER_SEC: 0.4, TOTAL_STATIC_HITS: 5"),
			args: LiteSpeedReport{},
			want: LiteSpeedReport{
				VirtualHostReport: map[string]map[string]float64{
					"hoge.jp": {
						"REQ_PROCESSING":             1,
						"REQ_PER_SEC":                0.1,
						"TOT_REQS":                   2,
						"PUB_CACHE_HITS_PER_SEC":     0.2,
						"TOTAL_PUB_CACHE_HITS":       3,
						"PRIVATE_CACHE_HITS_PER_SEC": 0.3,
						"TOTAL_PRIVATE_CACHE_HITS":   4,
						"STATIC_HITS_PER_SEC":        0.4,
						"TOTAL_STATIC_HITS":          5,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "ok_vhost_and_port",
			r: virtualHostLine("REQ_RATE [hoge.jp:80]: REQ_PROCESSING: 1, REQ_PER_SEC: 0.1, TOT_REQS: 2, PUB_CACHE_HITS_PER_SEC: 0.2, TOTAL_PUB_CACHE_HITS: 3, " +
				"PRIVATE_CACHE_HITS_PER_SEC: 0.3, TOTAL_PRIVATE_CACHE_HITS: 4, STATIC_HITS_PER_SEC: 0.4, TOTAL_STATIC_HITS: 5"),
			args: LiteSpeedReport{},
			want: LiteSpeedReport{
				VirtualHostReport: map[string]map[string]float64{
					"hoge.jp:80": {
						"REQ_PROCESSING":             1,
						"REQ_PER_SEC":                0.1,
						"TOT_REQS":                   2,
						"PUB_CACHE_HITS_PER_SEC":     0.2,
						"TOTAL_PUB_CACHE_HITS":       3,
						"PRIVATE_CACHE_HITS_PER_SEC": 0.3,
						"TOTAL_PRIVATE_CACHE_HITS":   4,
						"STATIC_HITS_PER_SEC":        0.4,
						"TOTAL_STATIC_HITS":          5,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "ok_Server",
			r: virtualHostLine("REQ_RATE []: REQ_PROCESSING: 2, REQ_PER_SEC: 0.3, TOT_REQS: 5, PUB_CACHE_HITS_PER_SEC: 0.6, TOTAL_PUB_CACHE_HITS: 3, " +
				"PRIVATE_CACHE_HITS_PER_SEC: 0.3, TOTAL_PRIVATE_CACHE_HITS: 4, STATIC_HITS_PER_SEC: 0.4, TOTAL_STATIC_HITS: 5"),
			args: LiteSpeedReport{},
			want: LiteSpeedReport{
				VirtualHostReport: map[string]map[string]float64{
					"Server": {
						"REQ_PROCESSING":             2,
						"REQ_PER_SEC":                0.3,
						"TOT_REQS":                   5,
						"PUB_CACHE_HITS_PER_SEC":     0.6,
						"TOTAL_PUB_CACHE_HITS":       3,
						"PRIVATE_CACHE_HITS_PER_SEC": 0.3,
						"TOTAL_PRIVATE_CACHE_HITS":   4,
						"STATIC_HITS_PER_SEC":        0.4,
						"TOTAL_STATIC_HITS":          5,
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "ng",
			r:       virtualHostLine("REQ_RATE [hoge.jp]:: REQ_PROCESSING: 1, REQ_PER_SEC: 0.1, TOT_REQS: 2, PUB_CACHE_HITS_PER_SEC: 0.2, TOTAL_PUB_CACHE_HITS: 3"),
			args:    LiteSpeedReport{},
			want:    LiteSpeedReport{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.VirtualHostReport = make(map[string]map[string]float64)
			tt.r.parse(&tt.args)
			if (tt.args.error != nil) != tt.wantErr {
				t.Errorf("(virtualHostLine)parse() error = %v, wantErr %v", tt.args.error, tt.wantErr)
			}
			if !tt.wantErr && !cmp.Equal(tt.args.VirtualHostReport, tt.want.VirtualHostReport) {
				t.Errorf("(virtualHostLine)parse() does not match. got = %v, want = %v", tt.args.VirtualHostReport, tt.want.VirtualHostReport)
			}
		})
	}
}

func Test_extAppLine_parse(t *testing.T) {
	tests := []struct {
		name    string
		e       extAppLine
		args    LiteSpeedReport
		want    LiteSpeedReport
		wantErr bool
	}{
		{
			name: "ok_vhost",
			e:    extAppLine("EXTAPP [LSAPI] [fuga.com] [fuga.com_php73]: CMAXCONN: 1, EMAXCONN: 2, POOL_SIZE: 3, INUSE_CONN: 4, IDLE_CONN: 5, WAITQUE_DEPTH: 6, REQ_PER_SEC: 0.7, TOT_REQS: 8"),
			args: LiteSpeedReport{},
			want: LiteSpeedReport{
				ExtAppReports: map[string]map[string]map[string]map[string]float64{
					"LSAPI": {
						"fuga.com": {
							"fuga.com_php73": {
								"CMAXCONN":      1,
								"EMAXCONN":      2,
								"POOL_SIZE":     3,
								"INUSE_CONN":    4,
								"IDLE_CONN":     5,
								"WAITQUE_DEPTH": 6,
								"REQ_PER_SEC":   0.7,
								"TOT_REQS":      8,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "ok_vhost_port",
			e:    extAppLine("EXTAPP [LSAPI] [fuga.com:80] [fuga.com_php73]: CMAXCONN: 1, EMAXCONN: 2, POOL_SIZE: 3, INUSE_CONN: 4, IDLE_CONN: 5, WAITQUE_DEPTH: 6, REQ_PER_SEC: 0.7, TOT_REQS: 8"),
			args: LiteSpeedReport{},
			want: LiteSpeedReport{
				ExtAppReports: map[string]map[string]map[string]map[string]float64{
					"LSAPI": {
						"fuga.com:80": {
							"fuga.com_php73": {
								"CMAXCONN":      1,
								"EMAXCONN":      2,
								"POOL_SIZE":     3,
								"INUSE_CONN":    4,
								"IDLE_CONN":     5,
								"WAITQUE_DEPTH": 6,
								"REQ_PER_SEC":   0.7,
								"TOT_REQS":      8,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "ok_cgi",
			e:    extAppLine("EXTAPP [CGI] [] [lscgid]: CMAXCONN: 2, EMAXCONN: 3, POOL_SIZE: 4, INUSE_CONN: 5, IDLE_CONN: 6, WAITQUE_DEPTH: 7, REQ_PER_SEC: 0.8, TOT_REQS: 9"),
			args: LiteSpeedReport{},
			want: LiteSpeedReport{
				ExtAppReports: map[string]map[string]map[string]map[string]float64{
					"CGI": {
						"Server": {
							"lscgid": {
								"CMAXCONN":      2,
								"EMAXCONN":      3,
								"POOL_SIZE":     4,
								"INUSE_CONN":    5,
								"IDLE_CONN":     6,
								"WAITQUE_DEPTH": 7,
								"REQ_PER_SEC":   0.8,
								"TOT_REQS":      9,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "ng",
			e:       extAppLine("EXTAPP [LSA:PI] [fuga.com]: CMAXCONN: 2, EMAXCONN: 3, POOL_SIZE: 4, INUSE_CONN: 5, IDLE_CONN: 6, WAITQUE_DEPTH: 7, REQ_PER_SEC: 0.8, TOT_REQS: 9"),
			args:    LiteSpeedReport{},
			want:    LiteSpeedReport{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.ExtAppReports = make(map[string]map[string]map[string]map[string]float64)
			tt.e.parse(&tt.args)
			if (tt.args.error != nil) != tt.wantErr {
				t.Errorf("(extAppLine)parse() error = %v, wantErr %v", tt.args.error, tt.wantErr)
			}
			if !tt.wantErr && !cmp.Equal(tt.args.ExtAppReports, tt.want.ExtAppReports) {
				t.Errorf("(extAppLine)parse() does not match. got = %v, want = %v", tt.args.ExtAppReports, tt.want.ExtAppReports)
			}
		})
	}
}

func Test_pickUpStringName(t *testing.T) {
	tests := []struct {
		name string
		args string
		want []string
	}{
		{
			name: "ok_single",
			args: "REQ_RATE [hoge.com]: REQ_PROCESSING: 1",
			want: []string{"hoge.com"},
		},
		{
			name: "ok_multi",
			args: "EXTAPP [LSAPI] [hoge.com] [hoge.com_php7.3]: CMAXCONN: 1000,",
			want: []string{"LSAPI", "hoge.com", "hoge.com_php7.3"},
		},
		{
			name: "ok_vhost_and_port",
			args: "EXTAPP [LSAPI] [hoge.com:80] [hoge.com_php7.3]: CMAXCONN: 1000,",
			want: []string{"LSAPI", "hoge.com:80", "hoge.com_php7.3"},
		},
		{
			name: "ok_trim_space",
			args: "EXTAPP [LSAPI] [ hoge.com] [hoge.com_php7.3]: CMAXCONN: 1000,",
			want: []string{"LSAPI", "hoge.com", "hoge.com_php7.3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pickUpStringName(tt.args); !cmp.Equal(got, tt.want) {
				t.Errorf("pickUpStringName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertStringToMap(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    map[string]float64
		wantErr bool
	}{
		{
			name:    "ok_single",
			args:    "xxxx: 1234",
			want:    map[string]float64{"xxxx": 1234},
			wantErr: false,
		},
		{
			name:    "ok_multi",
			args:    "xxxx: 1234, oooo: 432.1",
			want:    map[string]float64{"xxxx": 1234, "oooo": 432.1},
			wantErr: false,
		},
		{
			name:    "ng",
			args:    "xxxx: 1234 oooo: 432.1",
			want:    map[string]float64{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertStringToMap(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertStringToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !cmp.Equal(got, tt.want) {
				t.Errorf("convertStringToMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}
