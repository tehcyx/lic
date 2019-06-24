package license

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	if len(Licenses) == 0 {
		t.Errorf("initialization of licenses failed")
	}

	if _, ok := Licenses["na"]; !ok {
		t.Errorf("Fallback not initialized, something went wrong")
	}
}

func TestGet(t *testing.T) {
	type args struct {
		name    string
		version string
		branch  string
		url     string
	}
	tests := []struct {
		name string
		args args
		want License
	}{
		{name: "Get success", args: args{name: "github.com/tehcyx/girc", version: "0.0.1", branch: "", url: "http://github.com/tehcyx/girc"}, want: Licenses["apache-2.0"]},
		{name: "Get no license", args: args{name: "noproject", version: "", branch: "", url: "noprojecturl"}, want: Licenses["na"]},
		{name: "Get no license", args: args{name: "github.com/nouser/noproject", version: "", branch: "", url: "noprojecturl"}, want: Licenses["na"]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.args.name, tt.args.version, tt.args.branch, tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
