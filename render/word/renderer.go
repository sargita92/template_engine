package word

import (
	"strings"

	"github.com/sargita92/template_engine/ast"
)

// Renderer renderiza AST para WordprocessingML
type Renderer struct{}

// New cria um renderer Word
func New() *Renderer {
	return &Renderer{}
}

// RenderDocument converte um documento AST em XML Word
func (r *Renderer) RenderDocument(doc ast.Document) string {
	var out strings.Builder

	out.WriteString(`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">`)
	out.WriteString(`<w:body>`)

	for _, p := range doc.Paragraphs {
		out.WriteString(r.renderParagraph(p))
	}

	out.WriteString(`</w:body></w:document>`)

	return out.String()
}
