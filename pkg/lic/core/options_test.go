package core

import (
	"testing"
)

func TestNewOptions(t *testing.T) {
	got := NewOptions()

	if got == nil {
		t.Fatal("NewOptions() returned nil")
	}

	// Verify default values are set correctly
	if got.Verbose {
		t.Error("NewOptions() Verbose should default to false")
	}
}
