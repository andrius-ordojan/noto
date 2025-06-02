package main

import (
	"fmt"
	"log"
)

func run() error {
	scanner := NewScanner("./testdata/")
	notes, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("could not scan notes: %w", err)
	}

	for _, note := range notes {
		// fmt.Printf("Title: %s, Path: %s, Directive: %s\n", note.Title, note.Path, note.Directive)
		urls, err := note.getAllURLs()
		if err != nil {
			return fmt.Errorf("could not get URLs for note %s: %w", note.Title, err)
		}
		fmt.Printf("Note: %s, URLs: %v\n", note.Title, urls)
	}

	// TODO: archive to archive box and wait for content to apopear and use the result link to replace the URL in the note

	return nil

	// document.Dump(f, 2)

	// fmt.Println("Listing files in the directory:")
	// fmt.Println(findMDFiles("/home/andrius/Documents/obsidian-cabinet/resources/"))
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
