package layout

import (
	"strings"
	"testing"

	"github.com/kristyancarvalho/mdp/internal/model"
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
	return LayoutConfig{Width: width, Padding: 0, SoftWrap: true}
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
	lines := Layout(doc(para(long)), LayoutConfig{Width: 40, Padding: 0, SoftWrap: true})
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

func TestLayout_CodeBlock_StyleIsCodeBlock(t *testing.T) {
	lines := Layout(doc(code("go", "x := 1", "y := 2")), defaultCfg(80))
	found := 0
	for _, l := range lines {
		for _, seg := range l.Segments {
			if strings.Contains(seg.Text, "x := 1") || strings.Contains(seg.Text, "y := 2") {
				if seg.Style != StyleCodeBlock {
					t.Errorf("code line should have StyleCodeBlock, got %d", seg.Style)
				}
				found++
			}
		}
	}
	if found != 2 {
		t.Errorf("expected 2 code lines, found %d", found)
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

func TestLayout_OrderedList_NumbersPresent(t *testing.T) {
	lst := &model.List{
		Ordered: true,
		Items: []model.ListItem{
			{Blocks: []model.Block{para("first")}},
			{Blocks: []model.Block{para("second")}},
		},
	}
	lines := Layout(doc(lst), defaultCfg(80))
	flat := ""
	for _, l := range lines {
		for _, seg := range l.Segments {
			flat += seg.Text
		}
	}
	if !strings.Contains(flat, "1.") || !strings.Contains(flat, "2.") {
		t.Errorf("ordered list should contain '1.' and '2.', got: %q", flat)
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
	flat := ""
	for _, l := range lines {
		for _, seg := range l.Segments {
			flat += seg.Text
		}
	}
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
	flat := ""
	for _, l := range lines {
		for _, seg := range l.Segments {
			flat += seg.Text
		}
	}
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

	flat1, flat2 := "", ""
	for _, l := range l1 {
		for _, s := range l.Segments {
			flat1 += s.Text
		}
	}
	for _, l := range l2 {
		for _, s := range l.Segments {
			flat2 += s.Text
		}
	}
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