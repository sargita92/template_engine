package ast

// Segment representa um trecho de texto com estilo
type Segment struct {
	Text string

	Bold      bool
	Italic    bool
	Underline bool

	TextColor string // #RRGGBB
	Highlight string // #RRGGBB ou nome

	Tab bool
	Br  bool
}
