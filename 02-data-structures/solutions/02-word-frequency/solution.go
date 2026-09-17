package wordfreq

import (
	"sort"
	"strings"
	"unicode"
)

type WordCount struct {
	Word  string
	Count int
}

func sanitize(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return unicode.ToLower(r)
	}
	return ' '
}

// CountFrequencies returns lowercase, punctuation-free word counts.
func CountFrequencies(text string) map[string]int {
	cleaned := strings.Map(sanitize, text)
	words := strings.Fields(cleaned)

	counts := make(map[string]int, len(words))
	for _, w := range words {
		counts[w]++
	}
	return counts
}

// TopKWords returns the top k most frequent words sorted descending by count,
// with ties broken alphabetically in ascending order.
func TopKWords(counts map[string]int, k int) []WordCount {
	if k <= 0 || len(counts) == 0 {
		return []WordCount{}
	}

	items := make([]WordCount, 0, len(counts))
	for word, count := range counts {
		items = append(items, WordCount{Word: word, Count: count})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Count != items[j].Count {
			return items[i].Count > items[j].Count // Descending by count
		}
		return items[i].Word < items[j].Word // Ascending alphabetically
	})

	if k > len(items) {
		k = len(items)
	}

	return items[:k]
}
