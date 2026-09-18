# Challenge: Streaming Archive & Compression Pipeline

## 🎯 Objective
Build a streaming gzip compression and decompression pipeline using pure Go standard library primitives (`io.Reader`, `io.Writer`, `compress/gzip`, `bufio.Writer`).

## 📋 Requirements
1. **Streaming Compress**:
   ```go
   func StreamCompress(dst io.Writer, src io.Reader) (int64, int64, error)
   ```
   Reads from `src`, compresses on-the-fly with `gzip.Writer`, wraps `dst` with `bufio.Writer` for high throughput, ensures `gzip.Writer.Close()` and `bufio.Writer.Flush()` are executed in order, and returns `(uncompressedBytes, compressedBytes, error)`.

2. **Streaming Decompress**:
   ```go
   func StreamDecompress(dst io.Writer, src io.Reader) (int64, error)
   ```
   Reads compressed bytes from `src` using `gzip.NewReader`, streams uncompressed bytes to `dst`, and closes the gzip reader cleanly.

3. **Counting Wrapper**:
   Implement a lightweight `CountingWriter` or `CountingReader` using Go standard library interface composition.

4. **Zero intermediate file buffering**:
   Entire processing must stream chunk-by-chunk through standard library pipelines.

## 🧪 Verification
```powershell
go test -v ./07-standard-library/challenges/file-archive-streamer/...
```
