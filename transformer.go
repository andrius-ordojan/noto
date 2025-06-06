package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/genai"
)

type Transformer struct {
	VaultRootPath string
	client        *genai.Client
	ctx           context.Context
}

func NewTransformer(vaultRootPath, APIKey string) *Transformer {
	if vaultRootPath == "" {
		panic("VaultRootPath cannot be empty")
	}
	if APIKey == "" {
		panic("API Key cannot be empty")
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Transformer{
		VaultRootPath: vaultRootPath,
		client:        client,
		ctx:           ctx,
	}
}

func (t *Transformer) Process(note Note) error {
	if note.Title == "" {
		panic("Title is empty")
	}
	if note.RelVaultPath == "" {
		panic("RelVaultPath is empty")
	}
	if len(note.Content) == 0 {
		panic("note content must not be empty")
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(`
You are a technical summarizer.

Your task is to extract the most important concepts, steps, and code examples from the note and linked content below. Do NOT write an article. Instead:

- Use a bullet-point list format with clear sections and short explanations.
- Include key concepts, implementation steps, and relevant code samples.
- Keep it technical and concise. Avoid storytelling or introductions.
- Use Markdown formatting: headings (##), bullet points, and fenced code blocks.
- Target length: up to 1500 words. Prioritize clarity and density, but keep it as short as possible.

Guidance:
- Prioritize topics, technologies, or keywords explicitly mentioned in the note content.
- Expand in more detail on these prioritized topics, while still covering other relevant material.
- If multiple links are provided, merge overlapping content and avoid repeating basic explanations.
`, genai.RoleUser),
		Tools: []*genai.Tool{
			{
				URLContext: &genai.URLContext{},
			},
		},
	}

	prompt := genai.Text(fmt.Sprintf(`
Note metadata:
- Title: %s
- Path: %s

--- BEGIN NOTE CONTENT ---
%s
--- END NOTE CONTENT ---
`, note.Title, note.RelVaultPath, note.Content))

	result, err := t.client.Models.GenerateContent(
		t.ctx,
		"gemini-2.5-flash-preview-05-20",
		prompt,
		config,
	)
	if err != nil {
		return fmt.Errorf("could not generate content: %w", err)
	}

	note.AST.OwnerDocument().Meta()["noto"] = "done"
	newNote := fmt.Sprintf("%s\n%s\n\n+++++ original note bellow +++++\n%s", note.Fontmatter(), result.Text(), note.Content)

	notePath := filepath.Join(t.VaultRootPath, note.RelVaultPath)
	if err := os.WriteFile(notePath, []byte(newNote), 0644); err != nil {
		return fmt.Errorf("could not write note file %s: %w", notePath, err)
	}

	return nil
}
