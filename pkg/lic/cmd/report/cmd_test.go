// Package report implements the `lic report` command.
package report

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/tehcyx/lic/pkg/lic/core"
)

func TestNewReportOptions(t *testing.T) {
	type args struct {
		o *core.Options
	}
	tests := []struct {
		name string
		args args
		want *Options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReportOptions(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReportOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewReportCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReportCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReportCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
