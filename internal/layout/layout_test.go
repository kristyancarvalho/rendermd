package layout

import (
	"strings"
	"testing"

	"github.com/kristyancarvalho/rendermd/internal/model"
)

func doc(blocks ...model.Block) model.Document {
	return model.Document{Blocks: blocks}
}

func heading(level int, text string) *model.Heading {
	return &model.Heading{Level: level, Spans: []model.Span{&model.Text{Value: text}}}
}

func para(text string) *model.Paragraph {
	return &model.Paragraph{Spans: []model.Span{&model.Text{Value: text}}}
}

func code(lang string, lines ...string) *model.CodeBlock {
	return &model.CodeBlock{Lang: lang, Lines: lines}
}

func defaultCfg(width int) LayoutConfig {
	return LayoutConfig{Width: width, Padding: 0, SoftWrap: true, HideSyntax: true}
}

func flattenLines(lines []Line) string {
	var sb strings.Builder
	for _, l := range lines {
		for _, seg := range l.Segments {
			sb.WriteString(seg.Text)
		}
	}
	return sb.String()
}

func flattenLine(l Line) string {
	var sb strings.Builder
	for _, seg := range l.Segments {
		sb.WriteString(seg.Text)
	}
	return sb.String()
}

func TestLayout_H1_HasBlankLinesBefore(t *testing.T) {
	lines := Layout(doc(heading(1, "Title")), defaultCfg(80))
	if len(lines) < 3 {
		t.Fatalf("H1 should produce at least 3 lines (2 blank + text), got %d", len(lines))
	}
	if !isEmptyLine(lines[0]) || !isEmptyLine(lines[1]) {
		t.Error("H1 should be preceded by two blank lines")
	}
}

func TestLayout_H2_HasBlankLineBefore(t *testing.T) {
	lines := Layout(doc(heading(2, "Section")), defaultCfg(80))
	if len(lines) < 2 {
		t.Fatalf("H2 should produce at least 2 lines, got %d", len(lines))
	}
	if !isEmptyLine(lines[0]) {
		t.Error("H2 should be preceded by one blank line")
	}
}

func TestLayout_HeadingStyle(t *testing.T) {
	for level := 1; level <= 6; level++ {
		lines := Layout(doc(heading(level, "X")), defaultCfg(80))
		found := false
		for _, l := range lines {
			for _, seg := range l.Segments {
				if strings.Contains(seg.Text, "X") {
					st := headingStyle(level)
					if seg.Style != st {
						t.Errorf("H%d: want style %d, got %d", level, st, seg.Style)
					}
					found = true
				}
			}
		}
		if !found {
			t.Errorf("H%d: heading text not found in output", level)
		}
	}
}

func TestLayout_Paragraph_TrailingBlankLine(t *testing.T) {
	lines := Layout(doc(para("hello world")), defaultCfg(80))
	if len(lines) < 2 {
		t.Fatal("paragraph should have at least 2 lines (text + blank)")
	}
	if !isEmptyLine(lines[len(lines)-1]) {
		t.Error("paragraph should end with a blank line")
	}
}

func TestLayout_Paragraph_ContainsText(t *testing.T) {
	lines := Layout(doc(para("hello world")), defaultCfg(80))
	found := false
	for _, l := range lines {
		for _, seg := range l.Segments {
			if strings.Contains(seg.Text, "hello") {
				found = true
			}
		}
	}
	if !found {
		t.Error("paragraph text not found in layout output")
	}
}

func TestLayout_Paragraph_Wrapping(t *testing.T) {
	long := strings.Repeat("word ", 30)
	lines := Layout(doc(para(long)), LayoutConfig{Width: 40, Padding: 0, SoftWrap: true, HideSyntax: true})
	wrapped := 0
	for _, l := range lines {
		if !isEmptyLine(l) {
			wrapped++
		}
	}
	if wrapped < 2 {
		t.Errorf("expected wrapping to produce multiple lines, got %d content lines", wrapped)
	}
}

func TestLayout_InlineFormatting_PunctuationSpacing(t *testing.T) {
	p := &model.Paragraph{Spans: []model.Span{
		&model.Strong{Children: []model.Span{&model.Text{Value: "bold"}}},
		&model.Text{Value: ", "},
		&model.Emphasis{Children: []model.Span{&model.Text{Value: "italic"}}},
		&model.Text{Value: ", "},
		&model.InlineCode{Value: "code"},
		&model.Text{Value: ", "},
		&model.Link{Label: []model.Span{&model.Text{Value: "link"}}, URL: "https://example.com"},
		&model.Text{Value: "."},
	}}
	lines := Layout(doc(p), defaultCfg(80))
	flat := flattenLines(lines)
	want := "bold, italic, code, link."
	if !strings.Contains(flat, want) {
		t.Errorf("inline punctuation should remain attached: want %q in %q", want, flat)
	}
	if strings.Contains(flat, " ,") || strings.Contains(flat, " .") {
		t.Errorf("inline punctuation should not have leading spaces, got %q", flat)
	}
}

func TestLayout_InlineFormatting_SpanBoundarySpacing(t *testing.T) {
	p := &model.Paragraph{Spans: []model.Span{
		&model.Text{Value: "plain "},
		&model.Strong{Children: []model.Span{&model.Text{Value: "bold"}}},
		&model.Text{Value: " text"},
	}}
	lines := Layout(doc(p), defaultCfg(80))
	flat := flattenLines(lines)
	want := "plain bold text"
	if !strings.Contains(flat, want) {
		t.Errorf("inline whitespace should survive span boundaries: want %q in %q", want, flat)
	}
}

func TestLayout_LinkSegmentsKeepURL(t *testing.T) {
	p := &model.Paragraph{Spans: []model.Span{
		&model.Link{Label: []model.Span{&model.Text{Value: "site"}}, URL: "https://example.com"},
	}}
	lines := Layout(doc(p), defaultCfg(80))
	for _, l := range lines {
		for _, seg := range l.Segments {
			if seg.Style == StyleLink {
				if seg.URL != "https://example.com" {
					t.Errorf("link segment URL: want %q, got %q", "https://example.com", seg.URL)
				}
				return
			}
		}
	}
	t.Fatal("link segment not found")
}

func TestLayout_CodeBlock_ContainsCodeText(t *testing.T) {
	lines := Layout(doc(code("go", "x := 1", "y := 2")), defaultCfg(80))
	full := flattenLines(lines)
	if !strings.Contains(full, "x") || !strings.Contains(full, "1") {
		t.Error("code block output should contain 'x' and '1'")
	}
	if !strings.Contains(full, "y") || !strings.Contains(full, "2") {
		t.Error("code block output should contain 'y' and '2'")
	}
}

func TestLayout_CodeBlock_LinesPreserved(t *testing.T) {
	lines := Layout(doc(code("go", "x := 1", "y := 2")), defaultCfg(80))
	found := 0
	for _, l := range lines {
		trimmed := strings.TrimSpace(flattenLine(l))
		if trimmed == "x := 1" || trimmed == "y := 2" {
			found++
		}
	}
	if found != 2 {
		t.Errorf("expected 2 reconstructed code lines, found %d", found)
	}
}

func TestLayout_CodeBlock_PaddingUsesCodeStyle(t *testing.T) {
	lines := Layout(doc(code("go", "x := 1")), defaultCfg(20))
	var codeLine Line
	for _, l := range lines {
		if strings.Contains(flattenLine(l), "x := 1") {
			codeLine = l
			break
		}
	}
	if len(codeLine.Segments) == 0 {
		t.Fatal("code line not found")
	}
	if codeLine.Indent != 0 {
		t.Errorf("code padding should be emitted as styled text, got indent %d", codeLine.Indent)
	}
	if segmentsWidth(codeLine.Segments) != 20 {
		t.Errorf("code line should fill width 20, got %d", segmentsWidth(codeLine.Segments))
	}
	if codeLine.Segments[0].Style != StyleCodeBlock {
		t.Errorf("leading code padding should use StyleCodeBlock, got %d", codeLine.Segments[0].Style)
	}
	if codeLine.Segments[len(codeLine.Segments)-1].Style != StyleCodeBlock {
		t.Errorf("trailing code padding should use StyleCodeBlock, got %d", codeLine.Segments[len(codeLine.Segments)-1].Style)
	}
}

func TestLayout_CodeBlock_SurroundedByBlanks(t *testing.T) {
	lines := Layout(doc(code("", "line")), defaultCfg(80))
	if len(lines) < 3 {
		t.Fatal("code block should have blank before and after")
	}
	if !isEmptyLine(lines[0]) {
		t.Error("code block should start with blank line")
	}
	if !isEmptyLine(lines[len(lines)-1]) {
		t.Error("code block should end with blank line")
	}
}

func TestLayout_CodeBlock_UnknownLang_FallsBack(t *testing.T) {
	lines := Layout(doc(code("brainfuck", "some code")), defaultCfg(80))
	found := false
	for _, l := range lines {
		for _, seg := range l.Segments {
			if strings.Contains(seg.Text, "some code") {
				if seg.Style != StyleCodeBlock {
					t.Errorf("unknown lang: segment should have StyleCodeBlock, got %d", seg.Style)
				}
				found = true
			}
		}
	}
	if !found {
		t.Error("unknown lang: code text not found in output")
	}
}

func TestLayout_CodeBlock_EmptyLang_FallsBack(t *testing.T) {
	lines := Layout(doc(code("", "plain text")), defaultCfg(80))
	found := false
	for _, l := range lines {
		for _, seg := range l.Segments {
			if strings.Contains(seg.Text, "plain text") {
				if seg.Style != StyleCodeBlock {
					t.Errorf("empty lang: segment should have StyleCodeBlock, got %d", seg.Style)
				}
				found = true
			}
		}
	}
	if !found {
		t.Error("empty lang: code text not found in output")
	}
}

func TestLayout_CodeBlock_Go_HasSyntaxStyles(t *testing.T) {
	lines := Layout(doc(code("go", `func main() { // comment`)), defaultCfg(80))
	hasKeyword := false
	hasComment := false
	for _, l := range lines {
		for _, seg := range l.Segments {
			if seg.Style == StyleSyntaxKeyword {
				hasKeyword = true
			}
			if seg.Style == StyleSyntaxComment {
				hasComment = true
			}
		}
	}
	if !hasKeyword {
		t.Error("go code block should produce at least one StyleSyntaxKeyword segment")
	}
	if !hasComment {
		t.Error("go code block should produce at least one StyleSyntaxComment segment")
	}
}

func TestLayout_CodeBlock_Go_StringHighlighted(t *testing.T) {
	lines := Layout(doc(code("go", `fmt.Println("hello")`)), defaultCfg(80))
	hasString := false
	for _, l := range lines {
		for _, seg := range l.Segments {
			if seg.Style == StyleSyntaxString && strings.Contains(seg.Text, "hello") {
				hasString = true
			}
		}
	}
	if !hasString {
		t.Error("go string literal should produce StyleSyntaxString segment")
	}
}

func TestLayout_CodeBlock_Go_NumberHighlighted(t *testing.T) {
	lines := Layout(doc(code("go", "x := 42")), defaultCfg(80))
	hasNumber := false
	for _, l := range lines {
		for _, seg := range l.Segments {
			if seg.Style == StyleSyntaxNumber && strings.Contains(seg.Text, "42") {
				hasNumber = true
			}
		}
	}
	if !hasNumber {
		t.Error("go number literal should produce StyleSyntaxNumber segment")
	}
}

func TestLayout_CodeBlock_HideSyntax_False_ShowsLangHeader(t *testing.T) {
	cfg := LayoutConfig{Width: 80, HideSyntax: false}
	lines := Layout(doc(code("go", "x := 1")), cfg)
	found := false
	for _, l := range lines {
		for _, seg := range l.Segments {
			if strings.Contains(seg.Text, "go") && seg.Style == StyleMuted {
				found = true
			}
		}
	}
	if !found {
		t.Error("HideSyntax=false with lang should show a muted lang header line")
	}
}

func TestLayout_CodeBlock_HideSyntax_True_NoLangHeader(t *testing.T) {
	cfg := LayoutConfig{Width: 80, HideSyntax: true}
	lines := Layout(doc(code("go", "x := 1")), cfg)
	for _, l := range lines {
		for _, seg := range l.Segments {
			if seg.Text == " go" && seg.Style == StyleMuted {
				t.Error("HideSyntax=true should not show lang header line")
			}
		}
	}
}

func TestLayout_CodeBlock_HideSyntax_NoLang_NoHeader(t *testing.T) {
	cfg := LayoutConfig{Width: 80, HideSyntax: false}
	lines := Layout(doc(code("", "x := 1")), cfg)
	count := 0
	for _, l := range lines {
		if !isEmptyLine(l) {
			count++
		}
	}
	if count != 1 {
		t.Errorf("empty lang with HideSyntax=false should produce exactly 1 content line, got %d", count)
	}
}

func TestLayout_ThematicBreak_ContainsDashes(t *testing.T) {
	lines := Layout(doc(&model.ThematicBreak{}), defaultCfg(80))
	found := false
	for _, l := range lines {
		for _, seg := range l.Segments {
			if strings.Contains(seg.Text, "─") {
				found = true
				if seg.Style != StyleRule {
					t.Errorf("thematic break should have StyleRule, got %d", seg.Style)
				}
			}
		}
	}
	if !found {
		t.Error("thematic break should render horizontal rule characters")
	}
}

func TestLayout_UnorderedList_BulletPresent(t *testing.T) {
	lst := &model.List{
		Ordered: false,
		Items: []model.ListItem{
			{Blocks: []model.Block{para("item one")}},
			{Blocks: []model.Block{para("item two")}},
		},
	}
	lines := Layout(doc(lst), defaultCfg(80))
	bullets := 0
	for _, l := range lines {
		for _, seg := range l.Segments {
			if strings.Contains(seg.Text, "•") {
				bullets++
			}
		}
	}
	if bullets < 2 {
		t.Errorf("expected at least 2 bullet markers, got %d", bullets)
	}
}

func TestLayout_List_MarkersDoNotUseLineIndent(t *testing.T) {
	lst := &model.List{
		Ordered: false,
		Items: []model.ListItem{
			{Blocks: []model.Block{para("item one")}},
		},
	}
	lines := Layout(doc(lst), defaultCfg(80))
	for _, l := range lines {
		if !strings.Contains(flattenLine(l), "item one") {
			continue
		}
		if l.Indent != 0 {
			t.Errorf("list marker alignment should be in text segments, got indent %d", l.Indent)
		}
		if !strings.HasPrefix(flattenLine(l), "• item one") {
			t.Errorf("list line should start with marker, got %q", flattenLine(l))
		}
		return
	}
	t.Fatal("list item line not found")
}

func TestLayout_OrderedList_NumbersPresent(t *testing.T) {
	lst := &model.List{
		Ordered: true,
		Items: []model.ListItem{
			{Blocks: []model.Block{para("first")}},
			{Blocks: []model.Block{para("second")}},
		},
	}
	lines := Layout(doc(lst), defaultCfg(80))
	flat := flattenLines(lines)
	if !strings.Contains(flat, "1.") || !strings.Contains(flat, "2.") {
		t.Errorf("ordered list should contain '1.' and '2.', got: %q", flat)
	}
}

func TestLayout_OrderedList_TenItems_NumbersPresent(t *testing.T) {
	items := make([]model.ListItem, 10)
	for i := range items {
		items[i] = model.ListItem{Blocks: []model.Block{para("item")}}
	}
	lst := &model.List{Ordered: true, Items: items}
	lines := Layout(doc(lst), defaultCfg(80))
	flat := flattenLines(lines)
	if !strings.Contains(flat, "10.") {
		t.Errorf("ordered list with 10 items should contain '10.', got: %q", flat)
	}
	if strings.Contains(flat, ":") {
		t.Errorf("ordered list item 10 must not render as ':' (rune arithmetic bug), got: %q", flat)
	}
}

func TestLayout_OrderedList_MultiDigit_Alignment(t *testing.T) {
	items := make([]model.ListItem, 12)
	for i := range items {
		items[i] = model.ListItem{Blocks: []model.Block{para("item")}}
	}
	lst := &model.List{Ordered: true, Items: items}
	lines := Layout(doc(lst), defaultCfg(80))
	flat := flattenLines(lines)
	for _, n := range []string{"1.", "2.", "9.", "10.", "11.", "12."} {
		if !strings.Contains(flat, n) {
			t.Errorf("ordered list should contain %q, got: %q", n, flat)
		}
	}
}

func TestLayout_TaskList_CheckboxPresent(t *testing.T) {
	checked := true
	unchecked := false
	lst := &model.List{
		Items: []model.ListItem{
			{Blocks: []model.Block{para("done")}, Checked: &checked},
			{Blocks: []model.Block{para("todo")}, Checked: &unchecked},
		},
	}
	lines := Layout(doc(lst), defaultCfg(80))
	flat := flattenLines(lines)
	if !strings.Contains(flat, "[x]") {
		t.Errorf("expected checked checkbox '[x]', got: %q", flat)
	}
	if !strings.Contains(flat, "[ ]") {
		t.Errorf("expected unchecked checkbox '[ ]', got: %q", flat)
	}
}

func TestLayout_Quote_HasBarPrefix(t *testing.T) {
	q := &model.Quote{Blocks: []model.Block{para("quoted text")}}
	lines := Layout(doc(q), defaultCfg(80))
	found := false
	for _, l := range lines {
		for _, seg := range l.Segments {
			if strings.Contains(seg.Text, "▎") {
				found = true
				if seg.Style != StyleQuote {
					t.Errorf("quote bar should have StyleQuote, got %d", seg.Style)
				}
			}
		}
	}
	if !found {
		t.Error("blockquote should render a vertical bar prefix (▎)")
	}
}

func TestLayout_Quote_StylesCompleteLine(t *testing.T) {
	q := &model.Quote{Blocks: []model.Block{para("quoted text")}}
	lines := Layout(doc(q), defaultCfg(24))
	found := false
	for _, l := range lines {
		if !strings.Contains(flattenLine(l), "quoted text") {
			continue
		}
		found = true
		if segmentsWidth(l.Segments) != 24 {
			t.Errorf("quote line should fill width 24, got %d", segmentsWidth(l.Segments))
		}
		for _, seg := range l.Segments {
			if seg.Style != StyleQuote {
				t.Errorf("quote line segment should have StyleQuote, got %d", seg.Style)
			}
		}
	}
	if !found {
		t.Fatal("quote text line not found")
	}
}

func TestLayout_Quote_WrappedLinesUseQuoteStyle(t *testing.T) {
	q := &model.Quote{Blocks: []model.Block{para("one two three four five six seven")}}
	lines := Layout(doc(q), defaultCfg(14))
	styledWrapped := 0
	for _, l := range lines {
		if !strings.Contains(flattenLine(l), "▎") || isEmptyLine(l) {
			continue
		}
		allQuote := true
		for _, seg := range l.Segments {
			if seg.Style != StyleQuote {
				allQuote = false
			}
		}
		if allQuote {
			styledWrapped++
		}
	}
	if styledWrapped < 2 {
		t.Errorf("wrapped quote should produce multiple quote-styled lines, got %d", styledWrapped)
	}
}

func TestLayout_Quote_NestedQuotesKeepMarkers(t *testing.T) {
	q := &model.Quote{Blocks: []model.Block{
		&model.Quote{Blocks: []model.Block{para("nested")}},
	}}
	lines := Layout(doc(q), defaultCfg(30))
	found := false
	for _, l := range lines {
		flat := flattenLine(l)
		if strings.Contains(flat, "▎ ▎") && strings.Contains(flat, "nested") {
			found = true
			for _, seg := range l.Segments {
				if seg.Style != StyleQuote {
					t.Errorf("nested quote segment should have StyleQuote, got %d", seg.Style)
				}
			}
		}
	}
	if !found {
		t.Error("nested quote should retain both quote markers")
	}
}

func TestLayout_Table_HeadersPresent(t *testing.T) {
	tbl := &model.Table{
		Headers: [][]model.Span{
			{&model.Text{Value: "Name"}},
			{&model.Text{Value: "Age"}},
		},
		Rows: [][]model.Span{
			{&model.Text{Value: "Alice"}},
			{&model.Text{Value: "30"}},
		},
		Align: []model.TableAlign{model.AlignLeft, model.AlignRight},
	}
	lines := Layout(doc(tbl), defaultCfg(80))
	flat := flattenLines(lines)
	if !strings.Contains(flat, "Name") || !strings.Contains(flat, "Age") {
		t.Errorf("table should contain headers, got: %q", flat)
	}
	if !strings.Contains(flat, "Alice") {
		t.Errorf("table should contain row data, got: %q", flat)
	}
}

func TestEngine_Cache_ReturnsSameSlice(t *testing.T) {
	e := &Engine{}
	d := doc(para("hello"))
	cfg := defaultCfg(80)

	first := e.Render(d, cfg)
	second := e.Render(d, cfg)

	if len(first) > 0 && len(second) > 0 && &first[0] != &second[0] {
		t.Error("engine should return cached slice on identical input")
	}
}

func TestEngine_Cache_InvalidatesOnChange(t *testing.T) {
	e := &Engine{}
	cfg := defaultCfg(80)

	d1 := doc(para("hello"))
	l1 := e.Render(d1, cfg)

	d2 := doc(para("different content"))
	l2 := e.Render(d2, cfg)

	flat1 := flattenLines(l1)
	flat2 := flattenLines(l2)
	if flat1 == flat2 {
		t.Error("engine should re-render when document content changes")
	}
}

func TestWrapText_ShortLine(t *testing.T) {
	lines := wrapText("hello world", 80)
	if len(lines) != 1 {
		t.Errorf("short text should not wrap: got %d lines", len(lines))
	}
	if lines[0] != "hello world" {
		t.Errorf("want 'hello world', got %q", lines[0])
	}
}

func TestWrapText_LongLine(t *testing.T) {
	words := strings.Repeat("word ", 20)
	lines := wrapText(strings.TrimSpace(words), 20)
	if len(lines) < 2 {
		t.Errorf("long text should wrap into multiple lines, got %d", len(lines))
	}
}

func TestWrapText_EmptyString(t *testing.T) {
	lines := wrapText("", 80)
	if len(lines) != 1 || lines[0] != "" {
		t.Errorf("empty string should produce one empty line, got %v", lines)
	}
}

func TestPadLeft(t *testing.T) {
	if got := padLeft("5", 2); got != " 5" {
		t.Errorf("padLeft: want ' 5', got %q", got)
	}
	if got := padLeft("10", 2); got != "10" {
		t.Errorf("padLeft: want '10', got %q", got)
	}
}

func TestWrapText_MultibyteLatin_FitsCorrectly(t *testing.T) {
	lines := wrapText("café naïve", 80)
	if len(lines) != 1 {
		t.Errorf("short accented text should not wrap, got %d lines", len(lines))
	}
	if lines[0] != "café naïve" {
		t.Errorf("want %q, got %q", "café naïve", lines[0])
	}
}

func TestWrapText_CJK_WrapsAtVisualWidth(t *testing.T) {
	lines := wrapText("abc 你好", 6)
	if len(lines) != 2 {
		t.Errorf("CJK: want 2 lines, got %d: %v", len(lines), lines)
	}
	if len(lines) >= 1 && lines[0] != "abc" {
		t.Errorf("CJK first line: want 'abc', got %q", lines[0])
	}
	if len(lines) >= 2 && lines[1] != "你好" {
		t.Errorf("CJK second line: want '你好', got %q", lines[1])
	}
}

func TestWrapText_CJK_EachWordOnOwnLine(t *testing.T) {
	lines := wrapText("你好 世界 中文", 5)
	if len(lines) != 3 {
		t.Errorf("want 3 lines, got %d: %v", len(lines), lines)
	}
	for i, l := range lines {
		if len(strings.Fields(l)) > 1 {
			t.Errorf("line %d should contain one word at width 5, got %q", i, l)
		}
	}
}

func TestWrapText_MixedASCIIAndWide(t *testing.T) {
	lines := wrapText("hi 你好", 5)
	if len(lines) != 2 {
		t.Errorf("mixed ASCII+CJK: want 2 lines, got %d: %v", len(lines), lines)
	}
}

func TestPadLeft_WideChars(t *testing.T) {
	if got := padLeft("你好", 6); got != "  你好" {
		t.Errorf("padLeft wide: want %q, got %q", "  你好", got)
	}
	if got := padLeft("你好", 4); got != "你好" {
		t.Errorf("padLeft exact width: want %q, got %q", "你好", got)
	}
	if got := padLeft("你好世界", 4); got != "你好世界" {
		t.Errorf("padLeft overflow: want %q, got %q", "你好世界", got)
	}
}

func TestLayout_Paragraph_CJK_WrapsCorrectly(t *testing.T) {
	lines := Layout(doc(para("你好 你好 你好")), LayoutConfig{Width: 10, Padding: 0, SoftWrap: true, HideSyntax: true})
	content := 0
	for _, l := range lines {
		if !isEmptyLine(l) {
			content++
		}
	}
	if content < 2 {
		t.Errorf("CJK paragraph should wrap to ≥2 content lines at width 10, got %d", content)
	}
}
