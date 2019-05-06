package gopath

import (
	"testing"

	"github.com/tehcyx/lic/internal/licensereport"
)

func TestReadImports(t *testing.T) {
	type args struct {
		proj     *licensereport.Project
		filePath string
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
			if err := ReadImports(tt.args.proj, tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("ReadImports() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
