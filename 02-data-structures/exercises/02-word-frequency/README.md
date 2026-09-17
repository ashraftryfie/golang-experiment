# Exercise 02: Word Frequency Counter (Tier 2 - Medium)

## 🎯 Problem Statement
Implement a word frequency analyzer that processes text documents, counts word occurrences using a map, and extracts the Top-K most frequent words sorted deterministically.

### Requirements
1. `CountFrequencies(text string) map[string]int`:
   - Normalizes all words to lowercase.
   - Strips non-alphanumeric punctuation (e.g. commas, periods, exclamation points).
   - Splits on whitespace.
   - Returns a map of words to their counts.
2. `TopKWords(counts map[string]int, k int) []WordCount`:
   - Returns the top `k` most frequent words in descending order of count.
   - **Tie-breaking rule**: If two words have the same frequency, sort them alphabetically in ascending order.
   - If `k <= 0`, return an empty slice.
   - If `k > len(counts)`, return all words sorted.

```go
type WordCount struct {
    Word  string
    Count int
}
```

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `CountFrequencies` and `TopKWords`.
3. Use `strings.ToLower()`, `strings.Fields()`, and `sort.Slice()`.
4. Run tests:
   ```powershell
   go test -v ./02-data-structures/exercises/02-word-frequency
   ```
