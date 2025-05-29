package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func main() {
	workspaceHome := flag.String("workspace-home", "", "Path to the workspace home directory")

	flag.Parse()

	if *workspaceHome == "" {
		fmt.Fprintln(os.Stderr, "Error: --workspace-home must be provided")
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.ReadFile("/home/andrius/Documents/obsidian-cabinet/resources/usenet/indexers.md")
	if err != nil {
		log.Fatalf("could not read file: %v", err)
	}
	reader := text.NewReader(f)
	doc := goldmark.DefaultParser().Parse(reader)

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch node := n.(type) {
		case *ast.Heading:
			text := extractText(node, f)
			fmt.Printf("Heading (level %d): %s\n", node.Level, text)

		case *ast.Link:
			dest := string(node.Destination)
			label := extractText(node, f)
			fmt.Printf("Link: [%s](%s)\n", label, dest)

		case *ast.Paragraph:
			text := extractText(node, f)
			fmt.Printf("Paragraph: %s\n", text)
		}

		return ast.WalkContinue, nil
	})
}

func listFiles(dir string) []string {
	root := os.DirFS(dir)

	mdFiles, err := fs.Glob(root, "*.md")
	if err != nil {
		log.Fatalf("could not list files: %v", err)
	}

	var files []string
	for _, file := range mdFiles {
		files = append(files, path.Join(dir, file))
	}

	return files
}
