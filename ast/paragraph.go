package ast

// Paragraph representa um parágrafo lógico
type Paragraph struct {
	Align           Align
	IndentLevel     int    // nível de recuo (ident)
	FirstLineIndent int    // recuo da primeira linha (twips)
	BackgroundColor string // fundo do parágrafo (hex)

	Segments []Segment
}
