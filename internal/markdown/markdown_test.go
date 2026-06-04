package markdown

import (
	"testing"

	"github.com/kristyancarvalho/rendermd/internal/model"
)

func mustParse(t *testing.T, src string) model.Document {
	t.Helper()
	return Parse([]byte(src))
}

func firstBlock(t *testing.T, doc model.Document) model.Block {
	t.Helper()
	if len(doc.Blocks) == 0 {
		t.Fatal("document has no blocks")
	}
	return doc.Blocks[0]
}

func spansToText(spans []model.Span) string {
	s := ""
	for _, sp := range spans {
		switch v := sp.(type) {
		case *model.Text:
			s += v.Value
		case *model.InlineCode:
			s += v.Value
		case *model.Strong:
			s += spansToText(v.Children)
		case *model.Emphasis:
			s += spansToText(v.Children)
		case *model.Link:
			s += spansToText(v.Label)
		}
	}
	return s
}

func TestParse_H1(t *testing.T) {
	doc := mustParse(t, "# Hello World\n")
	h, ok := firstBlock(t, doc).(*model.Heading)
	if !ok {
		t.Fatalf("expected *model.Heading, got %T", firstBlock(t, doc))
	}
	if h.Level != 1 {
		t.Errorf("want level 1, got %d", h.Level)
	}
	if got := spansToText(h.Spans); got != "Hello World" {
		t.Errorf("want 'Hello World', got %q", got)
	}
}

func TestParse_H2_Through_H6(t *testing.T) {
	for level := 2; level <= 6; level++ {
		prefix := ""
		for i := 0; i < level; i++ {
			prefix += "#"
		}
		doc := mustParse(t, prefix+" Section\n")
		h, ok := firstBlock(t, doc).(*model.Heading)
		if !ok {
			t.Fatalf("H%d: expected Heading, got %T", level, firstBlock(t, doc))
		}
		if h.Level != level {
			t.Errorf("H%d: want level %d, got %d", level, level, h.Level)
		}
	}
}

func TestParse_Paragraph(t *testing.T) {
	doc := mustParse(t, "Simple paragraph.\n")
	p, ok := firstBlock(t, doc).(*model.Paragraph)
	if !ok {
		t.Fatalf("expected *model.Paragraph, got %T", firstBlock(t, doc))
	}
	if len(p.Spans) == 0 {
		t.Error("paragraph should have spans")
	}
}

func TestParse_BoldText(t *testing.T) {
	doc := mustParse(t, "**bold**\n")
	p := firstBlock(t, doc).(*model.Paragraph)
	var found bool
	for _, sp := range p.Spans {
		if _, ok := sp.(*model.Strong); ok {
			found = true
		}
	}
	if !found {
		t.Error("expected Strong span for **bold**")
	}
}

func TestParse_ItalicText(t *testing.T) {
	doc := mustParse(t, "_italic_\n")
	p := firstBlock(t, doc).(*model.Paragraph)
	var found bool
	for _, sp := range p.Spans {
		if _, ok := sp.(*model.Emphasis); ok {
			found = true
		}
	}
	if !found {
		t.Error("expected Emphasis span for _italic_")
	}
}

func TestParse_InlineCode(t *testing.T) {
	doc := mustParse(t, "Use `fmt.Println` here.\n")
	p := firstBlock(t, doc).(*model.Paragraph)
	var code *model.InlineCode
	for _, sp := range p.Spans {
		if ic, ok := sp.(*model.InlineCode); ok {
			code = ic
		}
	}
	if code == nil {
		t.Fatal("expected InlineCode span")
	}
	if code.Value != "fmt.Println" {
		t.Errorf("want 'fmt.Println', got %q", code.Value)
	}
}

func TestParse_Link(t *testing.T) {
	doc := mustParse(t, "[click here](https://example.com)\n")
	p := firstBlock(t, doc).(*model.Paragraph)
	var link *model.Link
	for _, sp := range p.Spans {
		if l, ok := sp.(*model.Link); ok {
			link = l
		}
	}
	if link == nil {
		t.Fatal("expected Link span")
	}
	if link.URL != "https://example.com" {
		t.Errorf("want URL 'https://example.com', got %q", link.URL)
	}
	if spansToText(link.Label) != "click here" {
		t.Errorf("want label 'click here', got %q", spansToText(link.Label))
	}
}

func TestParse_FencedCodeBlock(t *testing.T) {
	src := "```go\nfmt.Println(\"hello\")\n```\n"
	doc := mustParse(t, src)
	cb, ok := firstBlock(t, doc).(*model.CodeBlock)
	if !ok {
		t.Fatalf("expected CodeBlock, got %T", firstBlock(t, doc))
	}
	if cb.Lang != "go" {
		t.Errorf("want lang 'go', got %q", cb.Lang)
	}
	if len(cb.Lines) == 0 {
		t.Error("code block should have lines")
	}
	if cb.Lines[0] != `fmt.Println("hello")` {
		t.Errorf("want code line, got %q", cb.Lines[0])
	}
}

func TestParse_CodeBlock_NoLang(t *testing.T) {
	src := "```\nplain code\n```\n"
	doc := mustParse(t, src)
	cb, ok := firstBlock(t, doc).(*model.CodeBlock)
	if !ok {
		t.Fatalf("expected CodeBlock, got %T", firstBlock(t, doc))
	}
	if cb.Lang != "" {
		t.Errorf("want empty lang, got %q", cb.Lang)
	}
}

func TestParse_Blockquote(t *testing.T) {
	doc := mustParse(t, "> quoted text\n")
	q, ok := firstBlock(t, doc).(*model.Quote)
	if !ok {
		t.Fatalf("expected Quote, got %T", firstBlock(t, doc))
	}
	if len(q.Blocks) == 0 {
		t.Error("blockquote should have inner blocks")
	}
}

func TestParse_UnorderedList(t *testing.T) {
	src := "- apple\n- banana\n- cherry\n"
	doc := mustParse(t, src)
	lst, ok := firstBlock(t, doc).(*model.List)
	if !ok {
		t.Fatalf("expected List, got %T", firstBlock(t, doc))
	}
	if lst.Ordered {
		t.Error("want unordered list")
	}
	if len(lst.Items) != 3 {
		t.Errorf("want 3 items, got %d", len(lst.Items))
	}
}

func TestParse_OrderedList(t *testing.T) {
	src := "1. first\n2. second\n"
	doc := mustParse(t, src)
	lst, ok := firstBlock(t, doc).(*model.List)
	if !ok {
		t.Fatalf("expected List, got %T", firstBlock(t, doc))
	}
	if !lst.Ordered {
		t.Error("want ordered list")
	}
}

func TestParse_TaskList(t *testing.T) {
	src := "- [x] done\n- [ ] todo\n"
	doc := mustParse(t, src)
	lst, ok := firstBlock(t, doc).(*model.List)
	if !ok {
		t.Fatalf("expected List, got %T", firstBlock(t, doc))
	}
	if len(lst.Items) < 2 {
		t.Fatalf("want 2 items, got %d", len(lst.Items))
	}
	if lst.Items[0].Checked == nil || !*lst.Items[0].Checked {
		t.Error("first item should be checked")
	}
	if lst.Items[1].Checked == nil || *lst.Items[1].Checked {
		t.Error("second item should be unchecked")
	}
}

func TestParse_Table(t *testing.T) {
	src := "| Name | Age |\n|------|-----|\n| Alice | 30 |\n"
	doc := mustParse(t, src)
	tbl, ok := firstBlock(t, doc).(*model.Table)
	if !ok {
		t.Fatalf("expected Table, got %T", firstBlock(t, doc))
	}
	if len(tbl.Headers) != 2 {
		t.Errorf("want 2 headers, got %d", len(tbl.Headers))
	}
}

func TestParse_ThematicBreak(t *testing.T) {
	doc := mustParse(t, "---\n")
	_, ok := firstBlock(t, doc).(*model.ThematicBreak)
	if !ok {
		t.Fatalf("expected ThematicBreak, got %T", firstBlock(t, doc))
	}
}

func TestParse_Empty(t *testing.T) {
	doc := mustParse(t, "")
	if len(doc.Blocks) != 0 {
		t.Errorf("empty input should produce no blocks, got %d", len(doc.Blocks))
	}
}

func TestParse_MultipleBlocks(t *testing.T) {
	src := "# Title\n\nParagraph one.\n\nParagraph two.\n"
	doc := mustParse(t, src)
	if len(doc.Blocks) < 3 {
		t.Errorf("expected at least 3 blocks, got %d", len(doc.Blocks))
	}
}
