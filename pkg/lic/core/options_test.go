package core

import (
	"reflect"
	"testing"
)

func TestNewOptions(t *testing.T) {
	tests := []struct {
		name string
		want *Options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
