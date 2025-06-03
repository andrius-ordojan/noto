package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	nurl "net/url"
	"time"

	html2markdown "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/markusmobius/go-trafilatura"
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

var httpClient = &http.Client{Timeout: 30 * time.Second}

func runt() error {
	url := "https://scalegrid.io/blog/mongodb-rollback/"
	url = "https://news.ycombinator.com/item?id=44157177"
	url = "https://betterstack.com/community/guides/logging/how-to-start-logging-with-python/#logging-errors-in-python"
	parsedURL, err := nurl.ParseRequestURI(url)
	if err != nil {
		log.Fatalf("failed  to parse url: %v", err)
	}

	// Fetch article
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Fatalf("failed to fetch the page: %v", err)
	}
	defer resp.Body.Close()

	// Extract content
	opts := trafilatura.Options{
		IncludeImages: false,
		OriginalURL:   parsedURL,
	}

	result, err := trafilatura.Extract(resp.Body, opts)
	if err != nil {
		log.Fatalf("failed to extract: %v", err)
	}

	doc := trafilatura.CreateReadableDocument(result)
	fmt.Println(result.Content)

	htmlContent := []byte(result.Content)
	mdBuf := new(bytes.Buffer)
	converter := html2markdown.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(htmlContent)
	if err != nil {
		log.Fatalf("failed to convert html to markdown: %v", err)
	}
	fmt.Println(markdown)
	return nil
}

func main() {
	err := runt()
	// err := run()
	if err != nil {
		log.Fatal(err)
	}
}
