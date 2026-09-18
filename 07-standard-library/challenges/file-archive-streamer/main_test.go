package streamer

import (
	"bytes"
	"strings"
	"testing"
)

func TestStreamCompressAndDecompressRoundtrip(t *testing.T) {
	// Sample repetitious data that compresses well
	input := strings.Repeat("Production Go standard library io streaming pipelines.\n", 200)

	var compressedBuf bytes.Buffer
	rawCount, compressedCount, err := StreamCompress(&compressedBuf, strings.NewReader(input))
	if err != nil {
		t.Fatalf("StreamCompress failed: %v", err)
	}

	if rawCount != int64(len(input)) {
		t.Errorf("rawCount = %d, want %d", rawCount, len(input))
	}

	if compressedCount <= 0 || compressedCount >= rawCount {
		t.Errorf("compressedCount = %d, expected strictly less than %d", compressedCount, rawCount)
	}

	// Decompress roundtrip
	var decompressedBuf bytes.Buffer
	decompressedCount, err := StreamDecompress(&decompressedBuf, &compressedBuf)
	if err != nil {
		t.Fatalf("StreamDecompress failed: %v", err)
	}

	if decompressedCount != rawCount {
		t.Errorf("decompressedCount = %d, want %d", decompressedCount, rawCount)
	}

	if decompressedBuf.String() != input {
		t.Errorf("decompressed output does not match original input")
	}
}
