package fileworker

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileAnalyzer_AnalyzeFiles(t *testing.T) {
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "file1.txt")
	_ = os.WriteFile(file1, []byte("hello world\nsecond line\n"), 0644)

	file2 := filepath.Join(tempDir, "file2.txt")
	_ = os.WriteFile(file2, []byte("the quick brown fox jumps\n"), 0644)

	analyzer, err := NewFileAnalyzer(2)
	if err != nil {
		t.Fatalf("failed to create analyzer: %v", err)
	}

	results, summary, err := analyzer.AnalyzeFiles(context.Background(), []string{file1, file2})
	if err != nil {
		t.Fatalf("AnalyzeFiles failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 file results, got %d", len(results))
	}
	if summary.TotalFiles != 2 {
		t.Fatalf("expected 2 total files, got %d", summary.TotalFiles)
	}
	if summary.TotalLines != 3 {
		t.Fatalf("expected 3 total lines, got %d", summary.TotalLines)
	}
	if summary.TotalWords != 9 {
		t.Fatalf("expected 9 total words (4 in file1 + 5 in file2), got %d", summary.TotalWords)
	}
}

func TestFileAnalyzer_Cancellation(t *testing.T) {
	analyzer, _ := NewFileAnalyzer(2)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // pre-cancel context

	results, _, err := analyzer.AnalyzeFiles(ctx, []string{"dummy1.txt", "dummy2.txt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, res := range results {
		if res.Error == "" {
			t.Fatal("expected cancelled error on results")
		}
	}
}
