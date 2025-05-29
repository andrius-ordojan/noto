package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/text"
)

type Note struct {
	Title     string
	Path      string
	Directive string
}
type Scanner struct {
	RootPath string
}

func NewScanner(rootPath string) *Scanner {
	return &Scanner{
		RootPath: rootPath,
	}
}

func (s *Scanner) Scan() ([]Note, error) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.New(
				meta.WithStoresInDocument(),
			),
		),
	)

	notes := []Note{}
	err := filepath.WalkDir(s.RootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("could not access path %q: %w", path, err)
		}

		if d.IsDir() {
			return nil // Skip directories
		}

		if strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			f, err := os.ReadFile(path)
			if err != nil {
				log.Fatalf("could not read file: %v", err)
			}

			reader := text.NewReader(f)
			document := markdown.Parser().Parse(reader)
			metaData := document.OwnerDocument().Meta()

			if marked, ok := metaData["noto"]; ok {
				base := filepath.Base(path)
				ext := filepath.Ext(base)
				title := strings.TrimSuffix(base, ext)

				var directive string
				if marked != nil {
					directive = strings.TrimSpace(marked.(string))
				} else {
					directive = "default-directive" // TODO: fill in with default directive if nil
				}

				notes = append(notes, Note{
					Title:     title,
					Path:      path,
					Directive: directive,
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not walk the path %q: %w", s.RootPath, err)
	}

	return notes, nil
}
