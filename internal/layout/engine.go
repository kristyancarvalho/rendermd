package layout

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/kristyancarvalho/rendermd/internal/model"
	"github.com/kristyancarvalho/rendermd/internal/syntax"
	"github.com/mattn/go-runewidth"
)

type StyleID int

const (
	StyleNormal StyleID = iota
	StyleHeading1
	StyleHeading2
	StyleHeading3
	StyleHeading4
	StyleHeading5
	StyleHeading6
	StyleStrong
	StyleEmphasis
	StyleInlineCode
	StyleCodeBlock
	StyleQuote
	StyleLink
	StyleLinkURL
	StyleMuted
	StyleRule
	StyleAccent
	StyleSyntaxKeyword
	StyleSyntaxString
	StyleSyntaxComment
	StyleSyntaxNumber
	StyleSyntaxType
	StyleSyntaxBuiltin
	StyleSyntaxOperator
)

type Segment struct {
	Text  string
	Style StyleID
}

type Line struct {
	Segments []Segment
	Indent   int
}

type LayoutConfig struct {
	Width      int
	Padding    int
	SoftWrap   bool
	ShowURLs   bool
	HideSyntax bool
}

type Engine struct {
	lastHash  uint64
	lastLines []Line
}

func (e *Engine) Render(doc model.Document, cfg LayoutConfig) []Line {
	h := hashDoc(doc)
	if h == e.lastHash && e.lastLines != nil {
		return e.lastLines
	}
	lines := Layout(doc, cfg)
	e.lastHash = h
	e.lastLines = lines
	return lines
}

func Layout(doc model.Document, cfg LayoutConfig) []Line {
	effWidth := cfg.Width
	if cfg.Width <= 0 {
		effWidth = 80
	}
	avail := effWidth - 2*cfg.Padding
	if avail < 10 {
		avail = 10
	}

	var out []Line
	for _, block := range doc.Blocks {
		out = append(out, renderBlock(block, avail, cfg)...)
	}
	return out
}

func renderBlock(b model.Block, width int, cfg LayoutConfig) []Line {
	switch v := b.(type) {
	case *model.Heading:
		return renderHeading(v, width, cfg)
	case *model.Paragraph:
		return renderParagraph(v, width, cfg)
	case *model.CodeBlock:
		return renderCodeBlock(v, width, cfg)
	case *model.Quote:
		return renderQuote(v, width, cfg)
	case *model.List:
		return renderList(v, width, cfg)
	case *model.Table:
		return renderTable(v, width, cfg)
	case *model.ThematicBreak:
		return renderThematicBreak(width)
	case *model.ImagePlaceholder:
		return []Line{{Segments: []Segment{{Text: "[image: " + v.AltText + "]", Style: StyleMuted}}}}
	}
	return nil
}

func headingStyle(level int) StyleID {
	switch level {
	case 1:
		return StyleHeading1
	case 2:
		return StyleHeading2
	case 3:
		return StyleHeading3
	case 4:
		return StyleHeading4
	case 5:
		return StyleHeading5
	default:
		return StyleHeading6
	}
}

func renderHeading(h *model.Heading, width int, _ LayoutConfig) []Line {
	var lines []Line
	if h.Level == 1 {
		lines = append(lines, emptyLine(), emptyLine())
	} else {
		lines = append(lines, emptyLine())
	}
	st := headingStyle(h.Level)
	text := spansText(h.Spans)
	for _, wl := range wrapText(text, width) {
		lines = append(lines, Line{Segments: []Segment{{Text: wl, Style: st}}})
	}
	lines = append(lines, emptyLine())
	return lines
}

func renderParagraph(p *model.Paragraph, width int, cfg LayoutConfig) []Line {
	lines := wrapSpans(p.Spans, width, cfg)
	lines = append(lines, emptyLine())
	return lines
}

func kindToStyle(k syntax.TokenKind) StyleID {
	switch k {
	case syntax.KindKeyword:
		return StyleSyntaxKeyword
	case syntax.KindString:
		return StyleSyntaxString
	case syntax.KindComment:
		return StyleSyntaxComment
	case syntax.KindNumber:
		return StyleSyntaxNumber
	case syntax.KindType:
		return StyleSyntaxType
	case syntax.KindBuiltin:
		return StyleSyntaxBuiltin
	case syntax.KindOperator:
		return StyleSyntaxOperator
	default:
		return StyleCodeBlock
	}
}

func renderCodeBlock(c *model.CodeBlock, _ int, cfg LayoutConfig) []Line {
	var lines []Line
	lines = append(lines, emptyLine())

	if !cfg.HideSyntax && c.Lang != "" {
		lines = append(lines, Line{
			Segments: []Segment{{Text: " " + c.Lang, Style: StyleMuted}},
		})
	}

	for _, l := range c.Lines {
		tokens := syntax.Tokenize(c.Lang, l)
		segs := make([]Segment, 0, len(tokens))
		for i, tok := range tokens {
			text := tok.Text
			if i == 0 {
				text = " " + text
			}
			segs = append(segs, Segment{
				Text:  text,
				Style: kindToStyle(tok.Kind),
			})
		}
		if len(segs) == 0 {
			segs = []Segment{{Text: " ", Style: StyleCodeBlock}}
		}
		lines = append(lines, Line{Segments: segs, Indent: 1})
	}

	lines = append(lines, emptyLine())
	return lines
}

func renderQuote(q *model.Quote, width int, cfg LayoutConfig) []Line {
	var lines []Line
	lines = append(lines, emptyLine())
	inner := width - 2
	if inner < 5 {
		inner = 5
	}
	for _, b := range q.Blocks {
		for _, l := range renderBlock(b, inner, cfg) {
			if len(l.Segments) == 0 {
				lines = append(lines, Line{Segments: []Segment{{Text: "\u258e ", Style: StyleQuote}}})
				continue
			}
			prefixed := Line{Indent: l.Indent}
			prefixed.Segments = append([]Segment{{Text: "\u258e ", Style: StyleQuote}}, l.Segments...)
			lines = append(lines, prefixed)
		}
	}
	lines = append(lines, emptyLine())
	return lines
}

func renderList(lst *model.List, width int, cfg LayoutConfig) []Line {
	var lines []Line
	for i, item := range lst.Items {
		var marker string
		if lst.Ordered {
			marker = padLeft(fmt.Sprintf("%d", i+1), 2) + ". "
		} else {
			marker = "• "
		}
		if item.Checked != nil {
			if *item.Checked {
				marker = "[x] "
			} else {
				marker = "[ ] "
			}
		}
		markerWidth := runewidth.StringWidth(marker)
		inner := width - markerWidth
		if inner < 5 {
			inner = 5
		}
		first := true
		for _, b := range item.Blocks {
			for _, l := range renderBlock(b, inner, cfg) {
				if isEmptyLine(l) {
					continue
				}
				if first {
					nl := Line{Indent: l.Indent + markerWidth}
					nl.Segments = append([]Segment{{Text: marker, Style: StyleNormal}}, l.Segments...)
					lines = append(lines, nl)
					first = false
				} else {
					nl := Line{Indent: l.Indent + markerWidth}
					nl.Segments = append([]Segment{{Text: strings.Repeat(" ", markerWidth), Style: StyleNormal}}, l.Segments...)
					lines = append(lines, nl)
				}
			}
		}
	}
	lines = append(lines, emptyLine())
	return lines
}

func renderThematicBreak(width int) []Line {
	if width <= 0 {
		width = 40
	}
	rule := strings.Repeat("─", width)
	return []Line{
		emptyLine(),
		{Segments: []Segment{{Text: rule, Style: StyleRule}}},
		emptyLine(),
	}
}

func renderTable(t *model.Table, width int, _ LayoutConfig) []Line {
	var lines []Line
	lines = append(lines, emptyLine())

	nCols := len(t.Headers)
	if nCols == 0 {
		return lines
	}

	colWidths := make([]int, nCols)
	measureSpans := func(spans []model.Span) int {
		return runewidth.StringWidth(spansText(spans))
	}
	for i, h := range t.Headers {
		if w := measureSpans(h); w > colWidths[i] {
			colWidths[i] = w
		}
	}
	for r := 0; r+nCols <= len(t.Rows); r += nCols {
		for i := 0; i < nCols; i++ {
			if w := measureSpans(t.Rows[r+i]); w > colWidths[i] {
				colWidths[i] = w
			}
		}
	}

	total := 0
	for _, w := range colWidths {
		total += w + 3
	}
	if total > width && nCols > 0 {
		shrink := (total - width) / nCols
		for i := range colWidths {
			colWidths[i] = colWidths[i] - shrink
			if colWidths[i] < 3 {
				colWidths[i] = 3
			}
		}
	}

	align := func(i int) model.TableAlign {
		if i < len(t.Align) {
			return t.Align[i]
		}
		return model.AlignNone
	}

	formatCell := func(spans []model.Span, col int) string {
		text := spansText(spans)
		w := runewidth.StringWidth(text)
		max := colWidths[col]
		if w > max {
			runes := []rune(text)
			if max > 1 {
				text = string(runes[:max-1]) + "…"
			} else {
				text = string(runes[:max])
			}
			w = max
		}
		switch align(col) {
		case model.AlignCenter:
			left := (max - w) / 2
			right := max - w - left
			return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
		case model.AlignRight:
			return strings.Repeat(" ", max-w) + text
		default:
			return text + strings.Repeat(" ", max-w)
		}
	}

	buildRow := func(cells [][]model.Span, style StyleID) Line {
		var segs []Segment
		for i := 0; i < nCols; i++ {
			var spans []model.Span
			if i < len(cells) {
				spans = cells[i]
			}
			segs = append(segs, Segment{Text: " " + formatCell(spans, i) + " ", Style: style})
			if i < nCols-1 {
				segs = append(segs, Segment{Text: "│", Style: StyleMuted})
			}
		}
		return Line{Segments: segs}
	}

	headerCells := make([][]model.Span, nCols)
	copy(headerCells, t.Headers)
	lines = append(lines, buildRow(headerCells, StyleStrong))

	sep := make([]string, nCols)
	for i, w := range colWidths {
		sep[i] = strings.Repeat("─", w+2)
	}
	sepLine := Line{Segments: []Segment{{Text: strings.Join(sep, "┼"), Style: StyleMuted}}}
	lines = append(lines, sepLine)

	for r := 0; r+nCols <= len(t.Rows); r += nCols {
		rowCells := make([][]model.Span, nCols)
		for i := 0; i < nCols; i++ {
			rowCells[i] = t.Rows[r+i]
		}
		lines = append(lines, buildRow(rowCells, StyleNormal))
	}

	lines = append(lines, emptyLine())
	return lines
}

func emptyLine() Line {
	return Line{Segments: []Segment{{Text: "", Style: StyleNormal}}}
}

func isEmptyLine(l Line) bool {
	for _, s := range l.Segments {
		if strings.TrimSpace(s.Text) != "" {
			return false
		}
	}
	return true
}

func spansText(spans []model.Span) string {
	var sb strings.Builder
	for _, s := range spans {
		switch v := s.(type) {
		case *model.Text:
			sb.WriteString(v.Value)
		case *model.Emphasis:
			sb.WriteString(spansText(v.Children))
		case *model.Strong:
			sb.WriteString(spansText(v.Children))
		case *model.InlineCode:
			sb.WriteString(v.Value)
		case *model.Link:
			sb.WriteString(spansText(v.Label))
		}
	}
	return sb.String()
}

func wrapText(text string, width int) []string {
	if width <= 0 {
		width = 40
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}
	var lines []string
	var cur strings.Builder
	curWidth := 0
	for _, w := range words {
		wl := runewidth.StringWidth(w)
		if cur.Len() == 0 {
			cur.WriteString(w)
			curWidth = wl
		} else if curWidth+1+wl <= width {
			cur.WriteByte(' ')
			cur.WriteString(w)
			curWidth += 1 + wl
		} else {
			lines = append(lines, cur.String())
			cur.Reset()
			cur.WriteString(w)
			curWidth = wl
		}
	}
	if cur.Len() > 0 {
		lines = append(lines, cur.String())
	}
	return lines
}

func wrapSpans(spans []model.Span, width int, cfg LayoutConfig) []Line {
	type chunk struct {
		text  string
		style StyleID
	}
	var chunks []chunk
	for _, s := range spans {
		switch v := s.(type) {
		case *model.Text:
			chunks = append(chunks, chunk{v.Value, StyleNormal})
		case *model.Emphasis:
			chunks = append(chunks, chunk{spansText(v.Children), StyleEmphasis})
		case *model.Strong:
			chunks = append(chunks, chunk{spansText(v.Children), StyleStrong})
		case *model.InlineCode:
			chunks = append(chunks, chunk{v.Value, StyleInlineCode})
		case *model.Link:
			chunks = append(chunks, chunk{spansText(v.Label), StyleLink})
			if cfg.ShowURLs && v.URL != spansText(v.Label) {
				chunks = append(chunks, chunk{" (" + v.URL + ")", StyleLinkURL})
			}
		case *model.HardBreak:
			chunks = append(chunks, chunk{"\n", StyleNormal})
		}
	}

	var lines []Line
	var curLine []Segment
	curWidth := 0
	pendingSpace := false

	flushLine := func() {
		lines = append(lines, Line{Segments: curLine})
		curLine = nil
		curWidth = 0
		pendingSpace = false
	}

	for _, ch := range chunks {
		if ch.text == "\n" {
			flushLine()
			continue
		}
		words, spaces, trailingSpace := inlineTokens(ch.text)
		for wi, word := range words {
			spaceNeeded := 0
			if (pendingSpace || spaces[wi]) && curWidth > 0 {
				spaceNeeded = 1
			}
			wl := runewidth.StringWidth(word)
			if curWidth > 0 && curWidth+spaceNeeded+wl > width {
				flushLine()
				spaceNeeded = 0
			}
			text := word
			if spaceNeeded > 0 && curWidth > 0 {
				text = " " + word
			}
			curLine = append(curLine, Segment{Text: text, Style: ch.style})
			curWidth += runewidth.StringWidth(text)
			pendingSpace = false
		}
		if trailingSpace {
			pendingSpace = true
		}
	}
	if len(curLine) > 0 {
		flushLine()
	}
	if len(lines) == 0 {
		lines = append(lines, Line{Segments: []Segment{{Text: "", Style: StyleNormal}}})
	}
	return lines
}

func inlineTokens(text string) ([]string, []bool, bool) {
	var words []string
	var spaces []bool
	var cur strings.Builder
	pendingSpace := false
	for _, r := range text {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if cur.Len() > 0 {
				words = append(words, cur.String())
				spaces = append(spaces, pendingSpace)
				cur.Reset()
			}
			pendingSpace = true
			continue
		}
		cur.WriteRune(r)
	}
	if cur.Len() > 0 {
		words = append(words, cur.String())
		spaces = append(spaces, pendingSpace)
		pendingSpace = false
	}
	return words, spaces, pendingSpace
}

func padLeft(s string, n int) string {
	w := runewidth.StringWidth(s)
	if w >= n {
		return s
	}
	return strings.Repeat(" ", n-w) + s
}

func hashDoc(doc model.Document) uint64 {
	h := fnv.New64a()
	for _, b := range doc.Blocks {
		writeBlockHash(h, b)
	}
	return h.Sum64()
}

func writeBlockHash(h interface{ Write([]byte) (int, error) }, b model.Block) {
	switch v := b.(type) {
	case *model.Heading:
		h.Write([]byte{byte(v.Level)})
		for _, s := range v.Spans {
			writeSpanHash(h, s)
		}
	case *model.Paragraph:
		h.Write([]byte{2})
		for _, s := range v.Spans {
			writeSpanHash(h, s)
		}
	case *model.CodeBlock:
		h.Write([]byte{3})
		h.Write([]byte(v.Lang))
		for _, l := range v.Lines {
			h.Write([]byte(l))
		}
	case *model.ThematicBreak:
		h.Write([]byte{4})
	default:
		h.Write([]byte{0})
	}
}

func writeSpanHash(h interface{ Write([]byte) (int, error) }, s model.Span) {
	switch v := s.(type) {
	case *model.Text:
		h.Write([]byte(v.Value))
	case *model.InlineCode:
		h.Write([]byte(v.Value))
	}
}
