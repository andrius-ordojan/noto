package main

import (
	"fmt"
	"io/fs"
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
	Title        string
	RelVaultPath string
	AST          ast.Node
	Content      []byte
}

func (n *Note) Fontmatter() string {
	var fontmatter string

	for k, v := range n.AST.OwnerDocument().Meta() {
		if v == nil {
			v = ""
		}
		fontmatter += fmt.Sprintf("%s: %v\n", k, v)
	}

	if fontmatter == "" {
		return ""
	}

	return fmt.Sprintf("---\n%s---", fontmatter)
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
	VaultRootPath string
}

func NewScanner(vaultRootPath string) *Scanner {
	if vaultRootPath == "" {
		panic("VaultRootPath cannot be empty")
	}

	return &Scanner{
		VaultRootPath: vaultRootPath,
	}
}

func (s *Scanner) Scan() ([]Note, error) {
	return ScanFS(os.DirFS(s.VaultRootPath), ".")
}

func parseNote(path string, content []byte, parser goldmark.Markdown) (*Note, error) {
	reader := text.NewReader(content)
	document := parser.Parser().Parse(reader)
	metaData := document.OwnerDocument().Meta()

	marked, ok := metaData["noto"]
	if !ok {
		return nil, nil // not a valid note
	}

	if marked != nil && marked.(string) == "done" {
		return nil, nil // skip notes that are marked as done
	}

	base := filepath.Base(path)
	ext := filepath.Ext(base)
	title := strings.TrimSuffix(base, ext)

	return &Note{
		Title:        title,
		RelVaultPath: path,
		AST:          document,
		Content:      content,
	}, nil
}

func ScanFS(vault fs.FS, root string) ([]Note, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.New(meta.WithStoresInDocument()),
			extension.GFM,
		),
	)

	var notes []Note
	err := fs.WalkDir(vault, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("could not access %s: %w", path, err)
		}
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}

		fileData, err := fs.ReadFile(vault, path)
		if err != nil {
			return fmt.Errorf("could not read %s: %w", path, err)
		}

		note, err := parseNote(path, fileData, md)
		if err != nil {
			return err
		}
		if note != nil {
			notes = append(notes, *note)
		}
		return nil
	})

	return notes, err
}
