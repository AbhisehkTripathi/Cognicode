package services

import (
	"fmt"
	"strings"

	"github.com/gomutex/godocx"
)

type DocxGenerator struct{}

func NewDocxGenerator() *DocxGenerator {
	return &DocxGenerator{}
}

// Generate formatted .docx from structured text input
func (g *DocxGenerator) GenerateDocumentation(docText string, outputPath string) error {
	doc, err := godocx.NewDocument()
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	lines := strings.Split(docText, "\n")
	inCodeBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		switch {
		case trimmed == "":
			doc.AddEmptyParagraph()

		case strings.HasPrefix(trimmed, "```"):
			inCodeBlock = !inCodeBlock
			if inCodeBlock {
				p := doc.AddParagraph("Code Example:")
				p.AddText("Code Example:").Bold(true)
			} else {
				doc.AddEmptyParagraph()
			}

		case inCodeBlock:
			p := doc.AddParagraph("")
			p.AddText(trimmed)

		case strings.HasPrefix(trimmed, "# "):
			title := strings.TrimPrefix(trimmed, "# ")
			p := doc.AddParagraph(title)
			p.Style("Heading 1")

		case strings.HasPrefix(trimmed, "## "):
			subtitle := strings.TrimPrefix(trimmed, "## ")
			p := doc.AddParagraph(subtitle)
			p.Style("Heading 2")

		case strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* "):
			content := trimmed[2:]
			p := doc.AddParagraph(content)
			p.Style("List Bullet")

		default:
			doc.AddParagraph(trimmed)
		}
	}

	if err := doc.SaveTo(outputPath); err != nil {
		return fmt.Errorf("failed to save docx: %w", err)
	}

	return nil
}
