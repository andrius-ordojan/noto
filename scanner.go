package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"

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

func (n *Note) getAllURLs() ([]string, error) {
	if n.AST == nil {
		panic("AST is nil, please parse the document first")
	}
	if n.Content == nil {
		panic("Content is nil, it's needed to extract text from nodes")
	}

	urls := make(map[string]struct{})
	err := ast.Walk(n.AST, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		var currentURLString string

		switch ConcreteNode := node.(type) {
		case *ast.Link:
			currentURLString = string(ConcreteNode.Destination)
		case *ast.Image:
			currentURLString = string(ConcreteNode.Destination)
		case *ast.AutoLink:
			if ConcreteNode.AutoLinkType != ast.AutoLinkURL {
				return ast.WalkContinue, nil
			}
			currentURLString = string(ConcreteNode.URL(n.Content))
		}

		if currentURLString != "" {
			urls[currentURLString] = struct{}{}
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking the AST: %w", err)
	}

	resultURLs := make([]string, 0, len(urls))
	for u := range urls {
		resultURLs = append(resultURLs, u)
	}

	return resultURLs, nil
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
			extension.GFM,
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

			document.Dump(f, 0)
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
