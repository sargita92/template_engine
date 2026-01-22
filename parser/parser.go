package parser

import (
	"github.com/sargita92/template_engine/ast"
	"strings"
)

// Parse transforma um texto em AST.Document
func Parse(input string) ast.Document {
	doc := ast.Document{}

	// separa parágrafos por linha em branco
	blocks := strings.Split(input, "\n\n")

	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		p := ast.Paragraph{}
		p.Segments = parseInline(block)

		doc.Paragraphs = append(doc.Paragraphs, p)
	}

	return doc
}

func parseInline(text string) []ast.Segment {
	var segments []ast.Segment
	cur := ast.Segment{}

	for len(text) > 0 {
		switch {
		case strings.HasPrefix(text, "[b]"):
			flush(&segments, &cur)
			cur.Bold = true
			text = text[3:]

		case strings.HasPrefix(text, "[/b]"):
			flush(&segments, &cur)
			cur.Bold = false
			text = text[4:]

		case strings.HasPrefix(text, "[i]"):
			flush(&segments, &cur)
			cur.Italic = true
			text = text[3:]

		case strings.HasPrefix(text, "[/i]"):
			flush(&segments, &cur)
			cur.Italic = false
			text = text[4:]

		case strings.HasPrefix(text, "[u]"):
			flush(&segments, &cur)
			cur.Underline = true
			text = text[3:]

		case strings.HasPrefix(text, "[/u]"):
			flush(&segments, &cur)
			cur.Underline = false
			text = text[4:]

		default:
			cur.Text += string(text[0])
			text = text[1:]
		}
	}

	flush(&segments, &cur)
	return segments
}

func flush(list *[]ast.Segment, cur *ast.Segment) {
	if cur.Text != "" {
		*list = append(*list, *cur)
	}
	cur.Text = ""
}
