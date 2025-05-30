package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"

	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/text"
)

type Note struct {
	Title     string
	Path      string
	Directive string
	AST       ast.Node
	Content   []byte
}

func (n *Note) parse() ([]string, error) {
	if n.AST == nil {
		panic("AST is nil, please parse the document first")
	}

	var paragraphs []string
	var err error

	// find links in text and link nodes
	err = ast.Walk(n.AST, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if para, ok := node.(*ast.Paragraph); ok {
			var buf bytes.Buffer

			for c := para.FirstChild(); c != nil; c = c.NextSibling() {
				if t, ok := c.(*ast.Text); ok {
					buf.Write(t.Segment.Value(n.Content))
				}
				// You might want to handle other inline elements like Emphasis, Strong, Link, etc.
				// to reconstruct the full paragraph text accurately.
				// For simplicity, this example just extracts raw text segments.
			}
			paragraphs = append(paragraphs, buf.String())
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking the AST: %w", err)
	}

	return paragraphs, nil
	// ast.Walk(n.AST, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
	// 	if entering {
	// 		switch n := node.(type) {
	// 		case *ast.Link:
	// 			fmt.Printf("Found link: %s\n", n.Destination)
	// 		case *ast.Text:
	// 			fmt.Printf("Found text: %s\n", n.Text(source))
	// 		}
	// 	}
	// 	return ast.WalkContinue, nil
	// })
	// f, err := os.ReadFile(n.Path) if err != nil {
	// 	log.Fatalf("could not read file: %v", err)
	// }
	// n.AST.Dump(f, 2)
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
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.New(
				meta.WithStoresInDocument(),
			),
		),
	).Parser()

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
			document := md.Parse(reader)
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
					AST:       document,
					Content:   f,
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
