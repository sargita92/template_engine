package main

import (
	"log"

	"github.com/sargita92/template_engine/docx"
	"github.com/sargita92/template_engine/parser"
)

func main() {
	text := `
[b]Olá[/b] mundo

[i]Documento[/i] gerado via AST
`

	doc := parser.Parse(text)

	err := docx.Write("./examples/basic/testdata/ast_output.docx", doc)
	if err != nil {
		log.Fatal(err)
	}
}
