package main

import (
	"testing"
	"testing/fstest"
)

func TestScanFS(t *testing.T) {
	mockFS := fstest.MapFS{
		"folder2/valid_note.md": &fstest.MapFile{
			Data: []byte(`---
noto: test
---

A [link](https://example.com) and image ![img](https://example.com/image.jpg).
<https://autolink.com>`),
		},
		"invalid_note.md": &fstest.MapFile{
			Data: []byte(`# This file has no noto frontmatter`),
		},
	}

	notes, err := ScanFS(mockFS, ".")
	if err != nil {
		t.Fatalf("ScanFS failed: %v", err)
	}

	if len(notes) != 1 {
		t.Fatalf("expected 1 valid note, got %d", len(notes))
	}

	note := notes[0]
	if note.Title != "valid_note" {
		t.Errorf("expected title 'valid_note', got %q", note.Title)
	}

	if note.Directive != "test" {
		t.Errorf("expected directive 'test', got %q", note.Directive)
	}

	if note.RelVaultPath != "folder2/valid_note.md" {
		t.Errorf("expected RelVaultPath to be 'valid_note.md', got %q", note.RelVaultPath)
	}
}

func TestNote_getAllURLs(t *testing.T) {
	note := Note{
		Title:        "example",
		RelVaultPath: "example.md",
		Directive:    "test",
		Content: []byte(`---
noto: test
---

[Link](https://a.com)
<https://c.com>`),
	}

	notes, err := ScanFS(fstest.MapFS{
		"example.md": &fstest.MapFile{Data: note.Content},
	}, ".")
	if err != nil {
		t.Fatalf("ScanFS failed: %v", err)
	}

	note = notes[0]
	urls, err := note.getAllURLs()
	if err != nil {
		t.Fatalf("getAllURLs failed: %v", err)
	}

	expected := map[string]bool{
		"https://a.com": false,
		"https://c.com": false,
	}

	for _, u := range urls {
		if _, ok := expected[u]; ok {
			expected[u] = true
		} else {
			t.Errorf("unexpected URL found: %s", u)
		}
	}

	for u, found := range expected {
		if !found {
			t.Errorf("expected URL not found: %s", u)
		}
	}
}
