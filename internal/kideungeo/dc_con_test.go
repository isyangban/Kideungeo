package kideungeo

import (
	"testing"
)

func TestSearchDCCon(t *testing.T) {
	dcCons := searchDCCon("붕괴")

	if len(dcCons) == 0 {
		t.Error("Search Result should not be empty")
	}
}
