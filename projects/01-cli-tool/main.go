package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dbPath := filepath.Join(home, ".go-notes.json")
	store, err := NewNoteStore(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing note store: %v\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		title := addCmd.String("title", "", "Note title")
		content := addCmd.String("content", "", "Note content")
		tagsStr := addCmd.String("tags", "", "Comma-separated tags")
		_ = addCmd.Parse(os.Args[2:])

		var tags []string
		if *tagsStr != "" {
			for _, t := range strings.Split(*tagsStr, ",") {
				tags = append(tags, strings.TrimSpace(t))
			}
		}

		note, err := store.Add(*title, *content, tags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error adding note: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created note #%d: %s\n", note.ID, note.Title)

	case "list":
		notes := store.List()
		if len(notes) == 0 {
			fmt.Println("No notes found.")
			return
		}
		for _, n := range notes {
			fmt.Printf("[%d] %s (Tags: %s)\n    %s\n", n.ID, n.Title, strings.Join(n.Tags, ", "), n.Content)
		}

	case "search":
		searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
		query := searchCmd.String("query", "", "Search query")
		_ = searchCmd.Parse(os.Args[2:])

		matches := store.Search(*query)
		fmt.Printf("Found %d note(s) matching '%s':\n", len(matches), *query)
		for _, n := range matches {
			fmt.Printf("[%d] %s\n", n.ID, n.Title)
		}

	case "delete":
		delCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		id := delCmd.Int("id", 0, "Note ID to delete")
		_ = delCmd.Parse(os.Args[2:])

		if err := store.Delete(*id); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting note: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted note #%d\n", *id)

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: go-notes <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  add    -title <string> -content <string> [-tags <csv>]")
	fmt.Println("  list")
	fmt.Println("  search -query <string>")
	fmt.Println("  delete -id <int>")
}
