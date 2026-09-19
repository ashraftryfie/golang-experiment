# Project 1: `go-notes` (CLI Notes Manager)

A robust, test-driven Command Line Interface (CLI) note-taking application written in pure Go.

## Features
- **JSON File Persistence**: Notes are atomically stored and updated in structured JSON format.
- **Search & Filter**: Keyword search across note titles, contents, and tags.
- **Flag & Argument Parsing**: Standard Go `flag` package handling with subcommands (`add`, `list`, `search`, `delete`).
- **Comprehensive Unit Testing**: 100% test coverage with temporary test directories.

## Usage
```powershell
# Add a note
go run ./projects/01-cli-tool add -title "Go Concurrency" -content "Remember: Don't communicate by sharing memory" -tags "go,concurrency"

# List all notes
go run ./projects/01-cli-tool list

# Search notes
go run ./projects/01-cli-tool search -query "concurrency"

# Delete a note
go run ./projects/01-cli-tool delete -id 1
```
