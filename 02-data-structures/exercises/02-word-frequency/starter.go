package wordfreq

// WordCount holds a word and its observed count.
type WordCount struct {
	Word  string
	Count int
}

// CountFrequencies returns lowercase, punctuation-free word counts.
func CountFrequencies(text string) map[string]int {
	// TODO: Clean text, split words, count occurrences in map
	return nil
}

// TopKWords returns the top k most frequent words sorted descending by count,
// with ties broken alphabetically in ascending order.
func TopKWords(counts map[string]int, k int) []WordCount {
	// TODO: Extract slice from map, sort using sort.Slice, truncate to k
	return nil
}
