package main

import (
	"log"

	"github.com/sargita92/template_engine/engine"
)

func main() {
	e := engine.New()

	text := `
[b]Olá[/b] mundo

[i]Documento[/i] gerado pela Engine usando AST
`

	err := e.FillDocxFromText("engine_ast_output.docx", text)
	if err != nil {
		log.Fatal(err)
	}
}
