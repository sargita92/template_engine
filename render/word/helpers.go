package word

import (
	"strings"

	"github.com/sargita92/template_engine/ast"
)

func (r *Renderer) renderParagraph(p ast.Paragraph) string {
	var out strings.Builder

	out.WriteString("<w:p>")

	for _, s := range p.Segments {
		out.WriteString(renderSegment(s))
	}

	out.WriteString("</w:p>")

	return out.String()
}

func renderSegment(s ast.Segment) string {
	var out strings.Builder

	// TAB
	if s.Tab {
		return "<w:r><w:tab/></w:r>"
	}

	// BR
	if s.Br {
		return "<w:r><w:br/></w:r>"
	}

	out.WriteString("<w:r>")

	rPr := renderRunProperties(s)
	if rPr != "" {
		out.WriteString(rPr)
	}

	out.WriteString(`<w:t xml:space="preserve">`)
	out.WriteString(escapeXML(s.Text))
	out.WriteString("</w:t>")

	out.WriteString("</w:r>")

	return out.String()
}

func renderRunProperties(s ast.Segment) string {
	var out strings.Builder

	if !s.Bold && !s.Italic && !s.Underline {
		return ""
	}

	out.WriteString("<w:rPr>")

	if s.Bold {
		out.WriteString("<w:b/>")
	}
	if s.Italic {
		out.WriteString("<w:i/>")
	}
	if s.Underline {
		out.WriteString(`<w:u w:val="single"/>`)
	}

	out.WriteString("</w:rPr>")

	return out.String()
}

func escapeXML(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(s)
}
