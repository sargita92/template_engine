package main

import (
	"fmt"

	"github.com/sargita92/template_engine/parser"
	"github.com/sargita92/template_engine/render/word"
)

func main() {
	text := `
[b]Olá[/b] mundo
[i]Segunda[/i] linha
`

	doc := parser.Parse(text)

	renderer := word.New()
	xml := renderer.RenderDocument(doc)

	fmt.Println(xml)
}
