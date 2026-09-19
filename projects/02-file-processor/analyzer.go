package fileworker

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	ErrInvalidWorkerCount = errors.New("worker count must be greater than zero")
)

type FileStats struct {
	Path      string `json:"path"`
	Lines     int    `json:"lines"`
	Words     int    `json:"words"`
	Bytes     int64  `json:"bytes"`
	Error     string `json:"error,omitempty"`
}

type AggregatedSummary struct {
	TotalFiles int `json:"total_files"`
	TotalLines int `json:"total_lines"`
	TotalWords int `json:"total_words"`
	TotalBytes int64 `json:"total_bytes"`
}

type FileAnalyzer struct {
	workers int
}

func NewFileAnalyzer(workers int) (*FileAnalyzer, error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkerCount
	}
	return &FileAnalyzer{workers: workers}, nil
}

func (a *FileAnalyzer) AnalyzeFiles(ctx context.Context, filePaths []string) ([]FileStats, AggregatedSummary, error) {
	pathsCh := make(chan string, len(filePaths))
	for _, p := range filePaths {
		pathsCh <- p
	}
	close(pathsCh)

	resultsCh := make(chan FileStats, len(filePaths))
	var wg sync.WaitGroup

	for i := 0; i < a.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathsCh {
				select {
				case <-ctx.Done():
					resultsCh <- FileStats{Path: path, Error: ctx.Err().Error()}
					return
				default:
					stats := processFile(path)
					resultsCh <- stats
				}
			}
		}()
	}

	wg.Wait()
	close(resultsCh)

	var results []FileStats
	var summary AggregatedSummary

	for res := range resultsCh {
		results = append(results, res)
		if res.Error == "" {
			summary.TotalFiles++
			summary.TotalLines += res.Lines
			summary.TotalWords += res.Words
			summary.TotalBytes += res.Bytes
		}
	}

	return results, summary, nil
}

func processFile(path string) FileStats {
	file, err := os.Open(path)
	if err != nil {
		return FileStats{Path: path, Error: err.Error()}
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return FileStats{Path: path, Error: err.Error()}
	}

	scanner := bufio.NewScanner(file)
	lines := 0
	words := 0

	for scanner.Scan() {
		lines++
		line := scanner.Text()
		words += len(strings.Fields(line))
	}

	if err := scanner.Err(); err != nil {
		return FileStats{Path: path, Error: fmt.Sprintf("scan error: %v", err)}
	}

	return FileStats{
		Path:  path,
		Lines: lines,
		Words: words,
		Bytes: info.Size(),
	}
}
