package legacy

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	blockRe     = regexp.MustCompile(`\{\{BLOCK:([^}]+)\}\}`)
	parRe       = regexp.MustCompile(`(?s)(<w:p[^>]*>.*?</w:p>)`)
	alignFindRe = regexp.MustCompile(`(?si)\[align:([a-z]+)\](.*?)\[/align\]`)
	// ident regex (case-insensitive, accepts optional :N)
	identFindRe = regexp.MustCompile(`(?si)\[ident(?::\s*([0-9]+))?\](.*?)\[/ident\]`)
	bpgRe       = regexp.MustCompile(`(?si)\[bpg:([#a-z0-9]+)\](.*?)\[/bpg\]`)
)

// ======================================================================
//  PROCESSAMENTO PRINCIPAL
// ======================================================================

func ProcessPlaceholders(xml string, replaces map[string]string) string {
	if xml == "" {
		return xml
	}

	// 1 — Se o valor tiver estilos, promover automaticamente para BLOCK
	xml = promoteStyledPlaceholdersToBlock(xml, replaces)

	// 2 — Processar blocos ({{BLOCK:...}})
	xml = processBlock(xml, replaces)

	// 3 — Processar inline ({{FIELD}})
	xml = processInline(xml, replaces)

	return xml
}

// ======================================================================
//  PROMOÇÃO AUTOMÁTICA PARA BLOCK
// ======================================================================

func promoteStyledPlaceholdersToBlock(xml string, replaces map[string]string) string {
	for key, val := range replaces {

		// Se contiver marcadores, o placeholder precisa ser um BLOCK
		if strings.Contains(val, "[b]") || strings.Contains(val, "[/b]") ||
			strings.Contains(val, "[i]") || strings.Contains(val, "[/i]") ||
			strings.Contains(val, "[u]") || strings.Contains(val, "[/u]") ||
			strings.Contains(val, "[text:") || strings.Contains(val, "[/text]") ||
			strings.Contains(val, "[bg:") || strings.Contains(val, "[/bg]") ||
			strings.Contains(val, "[tab]") || strings.Contains(val, "[tab:") ||
			strings.Contains(val, "[br]") ||
			strings.Contains(val, "[align:") || strings.Contains(val, "[ident") {
			inline := "{{" + key + "}}"
			block := "{{BLOCK:" + key + "}}"
			xml = strings.ReplaceAll(xml, inline, block)
		}

		// Se o valor for muito longo, promovemos também para BLOCK (útil para jurisprudências)
		if len(val) > 800 {
			inline := "{{" + key + "}}"
			block := "{{BLOCK:" + key + "}}"
			xml = strings.ReplaceAll(xml, inline, block)
		}
	}
	return xml
}

func addParagraphBackground(pPr string, color string) string {
	if color == "" {
		return pPr
	}

	c := strings.TrimPrefix(color, "#")

	shd := fmt.Sprintf(
		`<w:shd w:val="clear" w:color="auto" w:fill="%s"/>`,
		strings.ToUpper(c),
	)

	if strings.TrimSpace(pPr) == "" {
		return "<w:pPr>" + shd + "</w:pPr>"
	}

	// remove shading existente
	re := regexp.MustCompile(`(?s)<w:shd[^>]*/>`)
	pPr = re.ReplaceAllString(pPr, "")

	idx := strings.Index(pPr, "</w:pPr>")
	if idx < 0 {
		return pPr + shd
	}

	return pPr[:idx] + shd + pPr[idx:]
}

// ======================================================================
//  PROCESSAMENTO DE BLOCK ({{BLOCK:CHAVE}})
//  - suporta [align:...] e [ident(:N)] em qualquer ponto do texto
// ======================================================================

func processBlock(xml string, replaces map[string]string) string {
	return parRe.ReplaceAllStringFunc(xml, func(par string) string {

		if !blockRe.MatchString(par) {
			return par
		}

		match := blockRe.FindStringSubmatch(par)
		if len(match) != 2 {
			return par
		}

		key := match[1]
		blockText, ok := lookupReplaceValue(key, replaces)
		if !ok {
			return blockRe.ReplaceAllString(par, "")
		}

		pPr, baseRPr := extractParagraphStyle(par)

		// partes separadas por parágrafos lógicos
		parts := strings.Split(blockText, "\n\n")

		var out strings.Builder

		for _, part := range parts {
			part = strings.TrimSpace(part)

			// -------- bpg (background de parágrafo) --------
			bpgColor := ""
			if m := bpgRe.FindStringSubmatch(part); len(m) == 3 {
				bpgColor = strings.TrimSpace(m[1])
				part = strings.TrimSpace(m[2])
			}

			if part == "" {
				// gera linha em branco quando a controller pediu
				out.WriteString("<w:p><w:r><w:br/></w:r></w:p>")
				continue
			}

			// Primeiro processamos ocorrências de ident inline: [ident(:N)]...[/ident]
			// Se houver, iteramos sobre matches e escrevemos antes/identificado/depois.
			lastIdx := 0
			identMatches := identFindRe.FindAllStringSubmatchIndex(part, -1)
			if len(identMatches) > 0 {
				for _, mi := range identMatches {
					// mi: [fullStart, fullEnd, grp1Start, grp1End (optional number), grp2Start, grp2End]
					fullStart, fullEnd := mi[0], mi[1]
					numStart, numEnd := mi[2], mi[3]
					grp2Start, grp2End := mi[4], mi[5]

					// texto antes do match
					if fullStart > lastIdx {
						before := strings.TrimSpace(part[lastIdx:fullStart])
						if before != "" {
							writeParagraphsFromText(&out, before, pPr, baseRPr)
						}
					}

					// número (opcional)
					n := 1
					if numStart >= 0 && numEnd > numStart {
						if v, err := strconv.Atoi(strings.TrimSpace(part[numStart:numEnd])); err == nil && v > 0 {
							n = v
						}
					}

					inner := strings.TrimSpace(part[grp2Start:grp2End])

					// calcula pPr com indent (120 twips por nível, limite 10)
					indentPPr := addParaIndent(pPr, n)
					// escreve inner com indent pPr
					writeParagraphsFromText(&out, inner, indentPPr, baseRPr)

					lastIdx = fullEnd
				}

				// texto depois do último match
				if lastIdx < len(part) {
					tail := strings.TrimSpace(part[lastIdx:])
					if tail != "" {
						writeParagraphsFromText(&out, tail, pPr, baseRPr)
					}
				}
				continue // parte processada
			}

			// se não houver ident, tratamos align (como antes)
			alignMatches := alignFindRe.FindAllStringSubmatchIndex(part, -1)
			if len(alignMatches) == 0 {
				// sem alinhamento dentro da parte → escreve normalmente (mantendo quebras internas)
				finalPPr := pPr
				if bpgColor != "" {
					finalPPr = addParagraphBackground(finalPPr, bpgColor)
				}

				writeParagraphsFromText(&out, part, finalPPr, baseRPr)

				continue
			}

			lastIdx = 0
			for _, mi := range alignMatches {
				fullStart, fullEnd := mi[0], mi[1]
				grp1Start, grp1End := mi[2], mi[3] // alignment value
				grp2Start, grp2End := mi[4], mi[5] // inner content

				// texto antes do match
				if fullStart > lastIdx {
					before := strings.TrimSpace(part[lastIdx:fullStart])
					if before != "" {
						writeParagraphsFromText(&out, before, pPr, baseRPr)
					}
				}

				// alignment value (case-insensitive)
				alignVal := strings.ToLower(strings.TrimSpace(part[grp1Start:grp1End]))
				inner := strings.TrimSpace(part[grp2Start:grp2End])

				finalPPr := pPr
				if bpgColor != "" {
					finalPPr = addParagraphBackground(finalPPr, bpgColor)
				}

				alignedPPr := addParaAlignment(finalPPr, alignVal)
				writeParagraphsFromText(&out, inner, alignedPPr, baseRPr)

				lastIdx = fullEnd
			}

			// texto depois do último match
			if lastIdx < len(part) {
				tail := strings.TrimSpace(part[lastIdx:])
				if tail != "" {
					writeParagraphsFromText(&out, tail, pPr, baseRPr)
				}
			}
		}

		return out.String()
	})
}

// helper que escreve um bloco de texto (possivelmente com quebras '\n')
// transformando cada linha (quebradas por '\n') em runs com <w:br/> entre elas,
// e monta os <w:p> aplicando o pPr fornecido.
// OBS: este helper NÃO divide por '\n\n' — cada chamada corresponde a um "part".
func writeParagraphsFromText(out *strings.Builder, text, pPr, baseRPr string) {
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimRight(line, " \t")

		runs, firstLine := styledTextToParagraphRuns(line, baseRPr)

		usedPPr := pPr
		if firstLine > 0 {
			usedPPr = addParaFirstLineIndent(pPr, firstLine)
		}

		out.WriteString("<w:p>")
		out.WriteString(usedPPr)
		out.WriteString(runs)
		out.WriteString("</w:p>")
	}
}

func styledTextToParagraphRuns(text, baseRPr string) (string, int) {
	segs := parseStyledSegments(text)

	firstLineIndent := 0
	var out strings.Builder

	for _, s := range segs {
		// captura o marcador
		if s.FirstLineIndent > 0 {
			firstLineIndent = s.FirstLineIndent
			continue
		}

		// TAB
		if s.Tab {
			out.WriteString("<w:r><w:tab/></w:r>")
			continue
		}

		// BR
		if s.Br {
			out.WriteString("<w:r><w:br/></w:r>")
			continue
		}

		out.WriteString("<w:r>")
		finalRPr := mergeRPr(baseRPr, s.Bold, s.Italic, s.Underline, s.Color, s.Highlight)
		if finalRPr != "" {
			out.WriteString(finalRPr)
		}

		out.WriteString(`<w:t xml:space="preserve">`)
		out.WriteString(escapeXML(s.Text))
		out.WriteString("</w:t></w:r>")
	}

	return out.String(), firstLineIndent
}

func addParaFirstLineIndent(pPr string, twips int) string {
	if twips <= 0 {
		twips = 360 // padrão Word (~0,63 cm)
	}

	ind := fmt.Sprintf(`<w:ind w:firstLine="%d"/>`, twips)

	// Se não houver pPr, cria um novo
	if strings.TrimSpace(pPr) == "" {
		return "<w:pPr>" + ind + "</w:pPr>"
	}

	// Remove qualquer <w:ind> existente (left ou firstLine)
	indRe := regexp.MustCompile(`(?s)<w:ind[^>]*\/?>`)
	pPr = indRe.ReplaceAllString(pPr, "")

	// Injeta antes do fechamento </w:pPr>
	idx := strings.Index(pPr, "</w:pPr>")
	if idx < 0 {
		// fallback seguro
		return pPr + ind
	}

	return pPr[:idx] + ind + pPr[idx:]
}

// ======================================================================
//  EXTRAÇÃO DE ESTILO DO PARÁGRAFO (<w:pPr>) E DO RUN (<w:rPr>)
// ======================================================================

// Retorna:
//   - pPr → estilo do parágrafo
//   - baseRPr → estilo padrão dos runs
func extractParagraphStyle(par string) (string, string) {

	// ---- Extrair <w:pPr> ----

	pPr := ""
	pPrRe := regexp.MustCompile(`(?s)(<w:pPr[^>]*>)(.*?)(</w:pPr>)`)

	if m := pPrRe.FindStringSubmatch(par); len(m) == 4 {
		inner := m[2]

		// Remove <w:rPr> internos do <w:pPr>
		inner = regexp.MustCompile(`(?s)<w:rPr[^>]*>.*?</w:rPr>`).ReplaceAllString(inner, "")

		pPr = m[1] + inner + m[3]
	}

	// ---- Extrair <w:rPr> do primeiro run ----

	baseRPr := ""
	rPrRe := regexp.MustCompile(`(?s)<w:r[^>]*>\s*(<w:rPr[^>]*>.*?</w:rPr>)`)

	if m := rPrRe.FindStringSubmatch(par); len(m) == 2 {
		baseRPr = m[1]
	}

	return pPr, baseRPr
}

// ======================================================================
//  PROCESSAMENTO INLINE: {{PLACEHOLDER}}
// ======================================================================

func processInline(xml string, replaces map[string]string) string {
	tRe := regexp.MustCompile(`(?s)<w:t[^>]*>(.*?)</w:t>`)

	return tRe.ReplaceAllStringFunc(xml, func(t string) string {
		m := tRe.FindStringSubmatch(t)
		if len(m) != 2 {
			return t
		}

		content := m[1]
		content = applyInlineReplaces(content, replaces)

		return strings.Replace(t, m[1], content, 1)
	})
}

var phRe = regexp.MustCompile(`\{\{([^{}]+)\}\}`)

func applyInlineReplaces(text string, replaces map[string]string) string {
	return phRe.ReplaceAllStringFunc(text, func(ph string) string {
		match := phRe.FindStringSubmatch(ph)
		if len(match) != 2 {
			return ph
		}

		key := match[1]

		// BLOCK não é tratado inline
		if strings.HasPrefix(key, "BLOCK:") {
			return ph
		}

		var fn, field string
		parts := strings.SplitN(key, ":", 2)
		if len(parts) == 2 {
			fn = strings.ToUpper(parts[0])
			field = parts[1]
		} else {
			field = key
		}

		value, ok := lookupReplaceValue(field, replaces)
		if !ok {
			return ph
		}

		if fn != "" {
			value = applyFunction(fn, value)
		}

		return escapeXML(value)
	})
}

// ======================================================================
//  ESTILIZAÇÃO — RUNS COM [b], [i], [u], [text:#hex], [bg:color], [tab], [br]
// ======================================================================

type segment struct {
	Text      string
	Bold      bool
	Italic    bool
	Underline bool
	Color     string
	Highlight string
	Tab       bool
	Br        bool

	FirstLineIndent int // twips; >0 indica marcador
}

// Converte texto estilizado em múltiplos <w:r>
func styledTextToRunsWithStyle(text, baseRPr string) string {
	segs := parseStyledSegments(text)

	var out strings.Builder
	for _, s := range segs {
		// Se for tab, emitir <w:r><w:tab/></w:r>
		if s.Tab {
			out.WriteString("<w:r><w:tab/></w:r>")
			continue
		}

		// Se for quebra de linha (dentro do parágrafo), emitir <w:r><w:br/></w:r>
		if s.Br {
			out.WriteString("<w:r><w:br/></w:r>")
			continue
		}

		out.WriteString("<w:r>")

		finalRPr := mergeRPr(baseRPr, s.Bold, s.Italic, s.Underline, s.Color, s.Highlight)
		if finalRPr != "" {
			out.WriteString(finalRPr)
		}

		out.WriteString(`<w:t xml:space="preserve">`)
		out.WriteString(escapeXML(s.Text))
		out.WriteString("</w:t></w:r>")
	}

	return out.String()
}

// Mescla o rPr original com novos estilos
func mergeRPr(base string, bold, italic, underline bool, color, hl string) string {

	hasStyle := bold || italic || underline || color != "" || hl != ""

	// Sem estilo algum
	if base == "" && !hasStyle {
		return ""
	}

	// Construção incremental dos estilos extras
	var extra strings.Builder

	if bold {
		extra.WriteString("<w:b/>")
	}
	if italic {
		extra.WriteString("<w:i/>")
	}
	if underline {
		extra.WriteString(`<w:u w:val="single"/>`)
	}
	if color != "" {
		// cor pode vir no formato #RRGGBB ou RRGGBB
		c := strings.TrimPrefix(color, "#")
		extra.WriteString(fmt.Sprintf(`<w:color w:val="%s"/>`, c))
	}
	if hl != "" {
		hl = strings.TrimSpace(strings.ToLower(hl))

		// HEX color → usar shading
		if strings.HasPrefix(hl, "#") {
			c := strings.TrimPrefix(hl, "#")

			// valida HEX básico (6 caracteres)
			if len(c) == 6 {
				extra.WriteString(
					fmt.Sprintf(
						`<w:shd w:val="clear" w:color="auto" w:fill="%s"/>`,
						strings.ToUpper(c),
					),
				)
			}

		} else {
			// nome de highlight válido do Word
			extra.WriteString(fmt.Sprintf(`<w:highlight w:val="%s"/>`, hl))
		}
	}

	// Se não existe estilo base → criar <w:rPr> novo
	if base == "" {
		return "<w:rPr>" + extra.String() + "</w:rPr>"
	}

	// Injeção no estilo existente
	idx := strings.Index(base, "</w:rPr>")
	if idx < 0 {
		return "<w:rPr>" + extra.String() + "</w:rPr>"
	}

	return base[:idx] + extra.String() + base[idx:]
}

// ======================================================================
//  PARSE DOS MARCADORES DE ESTILO
// ======================================================================

func parseStyledSegments(text string) []segment {
	var segments []segment
	cur := segment{}

	for len(text) > 0 {
		i := strings.Index(text, "[")
		if i < 0 {
			cur.Text += text
			break
		}

		if i > 0 {
			cur.Text += text[:i]
			text = text[i:]
		}

		switch {

		// ------------------- NEGRITO -------------------
		case strings.HasPrefix(text, "[b]"):
			segments = appendSegment(segments, &cur)
			cur.Bold = true
			text = text[3:]

		case strings.HasPrefix(text, "[/b]"):
			segments = appendSegment(segments, &cur)
			cur.Bold = false
			text = text[4:]

		// ------------------- ITÁLICO -------------------
		case strings.HasPrefix(text, "[i]"):
			segments = appendSegment(segments, &cur)
			cur.Italic = true
			text = text[3:]

		case strings.HasPrefix(text, "[/i]"):
			segments = appendSegment(segments, &cur)
			cur.Italic = false
			text = text[4:]

		// ------------------- SUBLINHADO -------------------
		case strings.HasPrefix(text, "[u]"):
			segments = appendSegment(segments, &cur)
			cur.Underline = true
			text = text[3:]

		case strings.HasPrefix(text, "[/u]"):
			segments = appendSegment(segments, &cur)
			cur.Underline = false
			text = text[4:]

		// ------------------- COR DO TEXTO -------------------
		case strings.HasPrefix(text, "[text:"):
			end := strings.Index(text, "]")
			if end > 0 {
				color := text[len("[text:"):end]
				segments = appendSegment(segments, &cur)
				cur.Color = color
				text = text[end+1:]
				continue
			}

		case strings.HasPrefix(text, "[/text]"):
			segments = appendSegment(segments, &cur)
			cur.Color = ""
			text = text[len("[/text]"):]
			continue

		// ------------------- FUNDO (HIGHLIGHT) -------------------
		case strings.HasPrefix(text, "[bg:"):
			end := strings.Index(text, "]")
			if end > 0 {
				hl := text[len("[bg:"):end]
				segments = appendSegment(segments, &cur)
				cur.Highlight = hl
				text = text[end+1:]
				continue
			}

		case strings.HasPrefix(text, "[/bg]"):
			segments = appendSegment(segments, &cur)
			cur.Highlight = ""
			text = text[len("[/bg]"):]
			continue

		// ------------------- TAB ([tab] or [tab:N]) -------------------
		case strings.HasPrefix(text, "[tab"):
			// pode ser [tab] ou [tab:3]
			end := strings.Index(text, "]")
			if end > 0 {
				inside := text[1:end] // exemplo: "tab" ou "tab:3"
				parts := strings.SplitN(inside, ":", 2)
				n := 1
				if len(parts) == 2 {
					if v, err := strconv.Atoi(parts[1]); err == nil && v > 0 {
						n = v
					}
				}
				segments = appendSegment(segments, &cur)
				for i := 0; i < n; i++ {
					segments = append(segments, segment{Tab: true})
				}
				text = text[end+1:]
				continue
			}

		// ------------------- QUEBRA DE LINHA [br] -------------------
		case strings.HasPrefix(text, "[br]"):
			segments = appendSegment(segments, &cur)
			segments = append(segments, segment{Br: true})
			text = text[4:]
			continue

		// ------------------- SPACE [space:N] (inserir N espaços)
		case strings.HasPrefix(text, "[space:"):
			end := strings.Index(text, "]")
			if end > 0 {
				numStr := text[len("[space:"):end]
				if n, err := strconv.Atoi(numStr); err == nil && n > 0 {
					// adiciona espaços ao texto corrente
					for i := 0; i < n; i++ {
						cur.Text += " "
					}
				}
				text = text[end+1:]
				continue
			}

		case strings.HasPrefix(text, "[firstLineIndent"):
			end := strings.Index(text, "]")
			if end > 0 {
				inside := text[len("[firstLineIndent"):end] // "" ou ":360"
				indent := 360                               // padrão Word

				if strings.HasPrefix(inside, ":") {
					if v, err := strconv.Atoi(strings.TrimSpace(inside[1:])); err == nil && v > 0 {
						indent = v
					}
				}

				segments = appendSegment(segments, &cur)

				segments = append(segments, segment{
					FirstLineIndent: indent,
				})

				text = text[end+1:]
				continue
			}

		default:
			cur.Text += string(text[0])
			text = text[1:]
		}
	}

	segments = appendSegment(segments, &cur)
	return segments
}

func appendSegment(list []segment, cur *segment) []segment {
	if cur.Text != "" {
		list = append(list, *cur)
	}
	cur.Text = ""
	return list
}

// ======================================================================
//  FUNÇÕES UTILITÁRIAS (lookup, máscaras, datas, escape, title case)
// ======================================================================

// Busca valor no mapa de replaces
func lookupReplaceValue(key string, rep map[string]string) (string, bool) {
	if v, ok := rep[key]; ok {
		return v, true
	}
	if v, ok := rep["{{"+key+"}}"]; ok {
		return v, true
	}
	return "", false
}

// Aplica funções do tipo UPPER, LOWER, TITLE, MASK_CPF, DATE_BR
func applyFunction(fn, val string) string {
	switch fn {
	case "UPPER":
		return strings.ToUpper(val)
	case "LOWER":
		return strings.ToLower(val)
	case "TITLE":
		return toTitle(val)
	case "MASK_CPF":
		return maskCPF(val)
	case "DATE_BR":
		return dateBR(val)
	case "PHONE_BR":
		return FormatPhoneBR(val)
	default:
		return val
	}
}

// Extrai apenas dígitos
func onlyDigits(s string) string {
	var out strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			out.WriteRune(r)
		}
	}
	return out.String()
}

// Formata CPF como 000.000.000-00 se tiver 11 dígitos
func maskCPF(s string) string {
	d := onlyDigits(s)
	if len(d) != 11 {
		return s
	}
	return d[:3] + "." + d[3:6] + "." + d[6:9] + "-" + d[9:]
}

// Converte 8 dígitos em dd/mm/aaaa
func dateBR(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	// 1️⃣ Formatos explícitos com separador
	layouts := []string{
		"2006-01-02", // YYYY-MM-DD (HTML date)
		"2006/01/02",
		"02/01/2006",
		"02-01-2006",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.Format("02/01/2006")
		}
	}

	// 2️⃣ Somente dígitos
	d := onlyDigits(s)
	if len(d) == 8 {

		// tenta YYYYMMDD
		if t, err := time.Parse("20060102", d); err == nil {
			return t.Format("02/01/2006")
		}

		// tenta DDMMYYYY
		if t, err := time.Parse("02012006", d); err == nil {
			return t.Format("02/01/2006")
		}
	}

	// 3️⃣ fallback seguro
	return s
}

// Escape de caracteres XML
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, `'`, "&apos;")
	return s
}

// Converte string para "Title Case" simples
func toTitle(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	// Conectivos que DEVEM ficar em minúsculo
	lowerWords := map[string]bool{
		"de":    true,
		"da":    true,
		"do":    true,
		"das":   true,
		"dos":   true,
		"e":     true,
		"des":   true,
		"del":   true,
		"della": true,
		"di":    true,
		"van":   true,
		"von":   true,
	}

	words := strings.Fields(strings.ToLower(s))
	for i, w := range words {

		// Algarismo romano → tudo maiúsculo
		if isRomanNumeral(w) {
			words[i] = strings.ToUpper(w)
			continue
		}

		// Conectivos → sempre minúsculo
		if lowerWords[w] {
			words[i] = w
			continue
		}

		// Capitalize unicode-safe
		r := []rune(w)
		r[0] = unicode.ToUpper(r[0])
		words[i] = string(r)
	}

	return strings.Join(words, " ")
}

func isRomanNumeral(s string) bool {
	s = strings.ToUpper(s)

	romanRe := regexp.MustCompile(`^(?i)(M{0,4}(CM|CD|D?C{0,3})` +
		`(XC|XL|L?X{0,3})(IX|IV|V?I{0,3}))$`)

	return romanRe.MatchString(s)
}

// ======================================================================
//  ALINHAMENTO E INDENTAÇÃO DE PARÁGRAFO
// ======================================================================

// addParaAlignment injeta a definição de alinhamento no pPr.
// alignment: "right", "center", "left", "both" etc.
func addParaAlignment(pPr string, alignment string) string {
	jc := fmt.Sprintf(`<w:jc w:val="%s"/>`, alignment)
	if strings.TrimSpace(pPr) == "" {
		return "<w:pPr>" + jc + "</w:pPr>"
	}

	// substitui <w:jc> existente se houver
	jcRe := regexp.MustCompile(`(?s)<w:jc[^>]*>.*?</w:jc>`)
	if jcRe.MatchString(pPr) {
		return jcRe.ReplaceAllString(pPr, jc)
	}

	// injetar antes do fechamento </w:pPr>
	idx := strings.Index(pPr, "</w:pPr>")
	if idx < 0 {
		// fallback: concatena
		return pPr + jc
	}
	return pPr[:idx] + jc + pPr[idx:]
}

// addParaIndent injeta recuo de parágrafo no pPr.
// level: número de níveis de indent (1 => 120 twips)
func addParaIndent(pPr string, level int) string {
	if level <= 0 {
		level = 1
	}
	leftVal := 120 * level
	ind := fmt.Sprintf(`<w:ind w:left="%d"/>`, leftVal)

	if strings.TrimSpace(pPr) == "" {
		return "<w:pPr>" + ind + "</w:pPr>"
	}

	// substitui <w:ind ...> existente no pPr, se houver
	indRe := regexp.MustCompile(`(?s)<w:ind[^>]*>.*?</w:ind>`)
	if indRe.MatchString(pPr) {
		return indRe.ReplaceAllString(pPr, ind)
	}

	// injetar antes do fechamento </w:pPr>
	idx := strings.Index(pPr, "</w:pPr>")
	if idx < 0 {
		// fallback: concatena
		return pPr + ind
	}
	return pPr[:idx] + ind + pPr[idx:]
}

func FormatPhoneBR(input string) string {
	var digits []rune
	for _, r := range input {
		if unicode.IsDigit(r) {
			digits = append(digits, r)
		}
	}

	n := len(digits)
	if n == 0 {
		return ""
	}

	if n < 3 {
		return string(digits)
	}

	ddd := string(digits[:2])
	rest := digits[2:]

	if n == 11 {
		return "(" + ddd + ") " +
			string(rest[:5]) + "-" +
			string(rest[5:])
	}

	if n == 10 {
		return "(" + ddd + ") " +
			string(rest[:4]) + "-" +
			string(rest[4:])
	}

	if len(rest) <= 4 {
		return "(" + ddd + ") " + string(rest)
	}

	return "(" + ddd + ") " +
		string(rest[:4]) + "-" +
		string(rest[4:])
}
