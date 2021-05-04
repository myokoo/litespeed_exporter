package rtreport

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_mergeSingleMap(t *testing.T) {
	type args struct {
		a map[string]float64
		b map[string]float64
	}
	tests := []struct {
		name string
		args args
		want map[string]float64
	}{
		{
			name: "ok",
			args: args{
				a: map[string]float64{"hoge": 10, "fuga": 0.1},
				b: map[string]float64{"hoge": 3, "aaa": 123},
			},
			want: map[string]float64{"hoge": 13, "fuga": 0.1, "aaa": 123},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergeSingleMap(tt.args.a, tt.args.b)
			if !cmp.Equal(tt.args.a, tt.want) {
				t.Errorf("mergeSingleMap() does not match. got = %v, want = %v", tt.args.a, tt.want)
			}
		})
	}
}

func Test_mergeDoubleMap(t *testing.T) {
	type args struct {
		a map[string]map[string]float64
		b map[string]map[string]float64
	}
	tests := []struct {
		name string
		args args
		want map[string]map[string]float64
	}{
		{
			name: "ok",
			args: args{
				a: map[string]map[string]float64{
					"hoge": {"fuga": 100, "xxxx": 123},
					"aaaa": {"abc": 1.1},
				},
				b: map[string]map[string]float64{
					"hoge": {"fuga": 321, "cccc": 10},
					"bbb":  {"dddd": 3456},
				},
			},
			want: map[string]map[string]float64{
				"hoge": {"fuga": 421, "xxxx": 123, "cccc": 10},
				"aaaa": {"abc": 1.1},
				"bbb":  {"dddd": 3456},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergeDoubleMap(tt.args.a, tt.args.b)
			if !cmp.Equal(tt.args.a, tt.want) {
				t.Errorf("mergeDoubleMap() does not match. got = %v, want = %v", tt.args.a, tt.want)
			}
		})
	}
}

func Test_mergeTripleMap(t *testing.T) {
	type args struct {
		a map[string]map[string]map[string]float64
		b map[string]map[string]map[string]float64
	}
	tests := []struct {
		name string
		args args
		want map[string]map[string]map[string]float64
	}{
		{

			name: "ok",
			args: args{
				a: map[string]map[string]map[string]float64{"aaa": {"bb": {"a": 1000, "b": 202}}},
				b: map[string]map[string]map[string]float64{"aaa": {"bb": {"a": 21}}, "ccc": {"dd": {"e": 100}}},
			},
			want: map[string]map[string]map[string]float64{"aaa": {"bb": {"a": 1021, "b": 202}}, "ccc": {"dd": {"e": 100}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergeTripleMap(tt.args.a, tt.args.b)
			if !cmp.Equal(tt.args.a, tt.want) {
				t.Errorf("mergeTripleMap() does not match. got = %v, want = %v", tt.args.a, tt.want)
			}
		})
	}
}

func Test_mergeQuadrupleMap(t *testing.T) {
	type args struct {
		a map[string]map[string]map[string]map[string]float64
		b map[string]map[string]map[string]map[string]float64
	}
	tests := []struct {
		name string
		args args
		want map[string]map[string]map[string]map[string]float64
	}{
		{

			name: "ok",
			args: args{
				a: map[string]map[string]map[string]map[string]float64{"aaaa": {"bbb": {"aa": {"a": 1000, "h": 544}, "bb": {"a": 1345}}}},
				b: map[string]map[string]map[string]map[string]float64{"aaaa": {"bbb": {"aa": {"h": 21}}}, "cccc": {"ddd": {"ee": {"z:": 456}}}},
			},
			want: map[string]map[string]map[string]map[string]float64{"aaaa": {"bbb": {"aa": {"a": 1000, "h": 565}, "bb": {"a": 1345}}}, "cccc": {"ddd": {"ee": {"z:": 456}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergeQuadrupleMap(tt.args.a, tt.args.b)
			if !cmp.Equal(tt.args.a, tt.want) {
				t.Errorf("mergeTripleMap() does not match. got = %v, want = %v", tt.args.a, tt.want)
			}
		})
	}
}
