package engine

import (
	"github.com/sargita92/template_engine/ast"
	"github.com/sargita92/template_engine/docx"
	"github.com/sargita92/template_engine/internal/legacy"
	"github.com/sargita92/template_engine/parser"
)

// Engine é a fachada pública da biblioteca
type Engine struct {
	replaces map[string]string
}

// New cria uma nova instância da engine
func New() *Engine {
	return &Engine{
		replaces: make(map[string]string),
	}
}

// WithReplaces define os placeholders (legacy)
func (e *Engine) WithReplaces(r map[string]string) *Engine {
	e.replaces = r
	return e
}

// ===============================
// LEGACY (mantido)
// ===============================

// ProcessXML processa XML usando o engine legado
func (e *Engine) ProcessXML(xml string) string {
	return legacy.ProcessPlaceholders(xml, e.replaces)
}

// FillDocx preenche um DOCX usando o engine legado
func (e *Engine) FillDocx(inputPath, outputPath string) error {
	return legacy.FillDocx(inputPath, outputPath, e.replaces)
}

// ===============================
// NOVA PIPELINE (AST)
// ===============================

// FillDocxFromText gera um DOCX a partir de texto puro usando AST
func (e *Engine) FillDocxFromText(outputPath string, text string) error {
	doc := parser.Parse(text)
	return docx.Write(outputPath, doc)
}

// FillDocxFromAST gera um DOCX diretamente a partir de um AST.Document
func (e *Engine) FillDocxFromAST(outputPath string, doc ast.Document) error {
	return docx.Write(outputPath, doc)
}
