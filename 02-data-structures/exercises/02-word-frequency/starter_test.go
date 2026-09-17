package wordfreq

import (
	"reflect"
	"testing"
)

func TestCountFrequencies(t *testing.T) {
	text := "Go is expressive, concise, clean, and efficient. Go is fun!"
	got := CountFrequencies(text)
	if got == nil {
		t.Skip("skipping: CountFrequencies is not yet implemented (implement in starter.go)")
	}

	if got["go"] != 2 {
		t.Errorf("expected count for 'go' to be 2, got %d", got["go"])
	}
	if got["clean"] != 1 {
		t.Errorf("expected count for 'clean' to be 1, got %d", got["clean"])
	}
}

func TestTopKWords(t *testing.T) {
	counts := map[string]int{
		"apple":  5,
		"banana": 5,
		"cherry": 2,
		"date":   10,
	}

	got := TopKWords(counts, 2)
	if got == nil {
		t.Skip("skipping: TopKWords is not yet implemented (implement in starter.go)")
	}

	want := []WordCount{
		{Word: "date", Count: 10},
		{Word: "apple", Count: 5}, // apple before banana due to alphabetical tie-breaking
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TopKWords() = %v, want %v", got, want)
	}
}
