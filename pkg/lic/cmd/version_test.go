package cmd

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/tehcyx/lic/pkg/lic/core"
)

func TestNewVersionOptions(t *testing.T) {
	type args struct {
		o *core.Options
	}
	tests := []struct {
		name string
		args args
		want *VersionOptions
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVersionOptions(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVersionOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewVersionCmd(t *testing.T) {
	type args struct {
		o *VersionOptions
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
			if got := NewVersionCmd(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVersionCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionOptions_Run(t *testing.T) {
	type fields struct {
		Options *core.Options
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &VersionOptions{
				Options: tt.fields.Options,
			}
			if err := o.Run(); (err != nil) != tt.wantErr {
				t.Errorf("VersionOptions.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
