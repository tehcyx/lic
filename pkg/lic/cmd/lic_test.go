package cmd

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/tehcyx/lic/pkg/lic/core"
)

func TestNewLicCmd(t *testing.T) {
	type args struct {
		o *core.Options
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLicCmd(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLicCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
