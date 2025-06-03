package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/go-shiori/go-readability"
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

func runr() error {
	// u := "https://betterstack.com/community/guides/logging/how-to-start-logging-with-python/#logging-errors-in-python"
	u := "https://last9.io/blog/python-logging-exceptions/"
	resp, err := http.Get(u)
	if err != nil {
		log.Fatalf("failed to download %s: %v\n", u, err)
	}
	defer resp.Body.Close()

	parsedURL, err := url.Parse(u)
	if err != nil {
		log.Fatalf("error parsing url")
	}

	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		log.Fatalf("failed to parse %s: %v\n", u, err)
	}

	// fmt.Printf("URL     : %s\n", u)
	// fmt.Printf("Title   : %s\n", article.Title)
	// fmt.Printf("Excerpt : %s\n", article.Excerpt)
	// fmt.Printf("SiteName: %s\n", article.Content)
	// fmt.Println("Content :")

	markdown, err := htmltomarkdown.ConvertString(article.Content)
	if err != nil {
		log.Fatalf("failed to convert html to markdown: %v", err)
	}
	fmt.Println(markdown)
	return nil
}

func main() {
	err := runr()
	// err := runt()
	// err := run()
	if err != nil {
		log.Fatal(err)
	}
}
