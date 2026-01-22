package ast

// Align representa alinhamento de parágrafo
type Align string

const (
	AlignLeft   Align = "left"
	AlignCenter Align = "center"
	AlignRight  Align = "right"
	AlignBoth   Align = "both"
)
