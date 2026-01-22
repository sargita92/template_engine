package main

import (
	"fmt"

	"github.com/sargita92/template_engine/parser"
)

func main() {
	text := `
[b]Olá[/b] mundo

[i]Segunda[/i] linha
`

	doc := parser.Parse(text)

	fmt.Printf("%#v\n", doc)
}
