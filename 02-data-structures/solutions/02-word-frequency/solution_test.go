package wordfreq

import (
	"reflect"
	"testing"
)

func TestSolutionWordFreq(t *testing.T) {
	text := "Go, go! Gopher loves Go."
	counts := CountFrequencies(text)

	if counts["go"] != 3 {
		t.Errorf("expected 'go' to be 3, got %d", counts["go"])
	}
	if counts["gopher"] != 1 {
		t.Errorf("expected 'gopher' to be 1, got %d", counts["gopher"])
	}

	top := TopKWords(counts, 2)
	want := []WordCount{
		{Word: "go", Count: 3},
		{Word: "gopher", Count: 1}, // gopher before loves
	}
	if !reflect.DeepEqual(top, want) {
		t.Errorf("got %v, want %v", top, want)
	}
}
