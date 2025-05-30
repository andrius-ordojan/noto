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
		fmt.Printf("Title: %s, Path: %s, Directive: %s\n", note.Title, note.Path, note.Directive)
		note.parse()
	}

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
