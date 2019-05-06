package godep

import (
	"testing"

	"github.com/tehcyx/lic/internal/licensereport"
)

func TestCollect(t *testing.T) {
	type args struct {
		proj    *licensereport.Project
		prjPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Collect(tt.args.proj, tt.args.prjPath); (err != nil) != tt.wantErr {
				t.Errorf("Collect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
