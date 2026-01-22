package main

import (
	"log"

	"github.com/sargita92/template_engine/engine"
)

func main() {
	e := engine.New().
		WithReplaces(map[string]string{
			"NOME": "João da Silva",
			"DATA": "2026-01-22",
			"TEXTO": `
					[b]Texto em negrito[/b]
					[i]Itálico[/i]
					
					[align:center][bg:#FFFF00]Parágrafo centralizado com fundo[/bg][/align]
					
					[ident:2]Parágrafo com recuo[/ident]
			`,
		})

	err := e.FillDocx("./testdata/input.docx", "./testdata/output.docx")

	if err != nil {
		log.Fatal(err)
	}
}
