package streamer

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
)

// CountingWriter tracks total bytes written through it.
type CountingWriter struct {
	w     io.Writer
	count int64
}

func NewCountingWriter(w io.Writer) *CountingWriter {
	return &CountingWriter{w: w}
}

func (cw *CountingWriter) Write(p []byte) (n int, err error) {
	n, err = cw.w.Write(p)
	cw.count += int64(n)
	return n, err
}

func (cw *CountingWriter) BytesWritten() int64 {
	return cw.count
}

// StreamCompress compresses src into dst using gzip compression, returning (rawBytes, compressedBytes, error).
func StreamCompress(dst io.Writer, src io.Reader) (int64, int64, error) {
	bufWriter := bufio.NewWriter(dst)
	countWriter := NewCountingWriter(bufWriter)
	gzWriter := gzip.NewWriter(countWriter)

	// Stream src into gzip writer
	rawBytes, err := io.Copy(gzWriter, src)
	if err != nil {
		_ = gzWriter.Close()
		return 0, 0, fmt.Errorf("compress stream error: %w", err)
	}

	// gzip.Writer MUST be closed to flush trailer & footer checksum
	if err := gzWriter.Close(); err != nil {
		return 0, 0, fmt.Errorf("gzip close error: %w", err)
	}

	// bufio.Writer MUST be flushed to write any remaining buffered bytes to dst
	if err := bufWriter.Flush(); err != nil {
		return 0, 0, fmt.Errorf("bufio flush error: %w", err)
	}

	return rawBytes, countWriter.BytesWritten(), nil
}

// StreamDecompress unpacks gzip-compressed data from src and writes uncompressed bytes to dst.
func StreamDecompress(dst io.Writer, src io.Reader) (int64, error) {
	gzReader, err := gzip.NewReader(src)
	if err != nil {
		return 0, fmt.Errorf("gzip reader error: %w", err)
	}
	defer gzReader.Close()

	n, err := io.Copy(dst, gzReader)
	if err != nil {
		return 0, fmt.Errorf("decompress copy error: %w", err)
	}

	return n, nil
}
