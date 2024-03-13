package bins

import (
	"reflect"
	"testing"
)

func TestNormalization(t *testing.T) {
	bins := []Bin{{1, 1, 0.1}, {0, 2, 0.2}, {1, 1, 0.8}}
	wantNBins := []NBin{{true, 0.5, 0.125}, {false, 1, 0.25}, {true, 0.5, 1}}

	nBins := make([]NBin, len(bins))
	Normilize(bins, nBins)

	if !reflect.DeepEqual(nBins, wantNBins) {
		t.Fatalf("invalid normalize result: want %v, got %v", wantNBins, nBins)
	}
}
