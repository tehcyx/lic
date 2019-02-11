package licensereport

import "testing"

func TestInit(t *testing.T) {
	if len(Licenses) == 0 {
		t.Errorf("initialization of licenses failed")
	}

	if _, ok := Licenses["na"]; !ok {
		t.Errorf("Fallback not initialized, something went wrong")
	}
}
