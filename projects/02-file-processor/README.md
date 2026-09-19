# Project 2: `file-worker` (Concurrent File Analyzer)

A concurrent, bounded multi-threaded file analyzer designed to process files, calculate line and word counts, and aggregate character frequencies without exceeding system file descriptor limits.

## Features
- **Bounded Worker Pool**: Concurrency level dynamically configured to prevent thread and FD starvation.
- **Context Cancellation**: Halts processing instantly when context deadline or cancellation signal triggers.
- **Aggregated Results**: Thread-safe stats collection with line counts, word counts, and byte sizes.

## Usage
```powershell
go test -v ./projects/02-file-processor/...
```
