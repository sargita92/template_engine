package engine

import "github.com/sargita92/template_engine/internal/legacy"

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

// WithReplaces define os placeholders
func (e *Engine) WithReplaces(r map[string]string) *Engine {
	e.replaces = r
	return e
}

// ProcessXML processa um XML usando o engine legado
func (e *Engine) ProcessXML(xml string) string {
	return legacy.ProcessPlaceholders(xml, e.replaces)
}

// FillDocx preenche um DOCX usando o engine legado
func (e *Engine) FillDocx(inputPath, outputPath string) error {
	return legacy.FillDocx(inputPath, outputPath, e.replaces)
}
