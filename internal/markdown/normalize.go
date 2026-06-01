package markdown

import (
	"bytes"

	"github.com/kristyancarvalho/mdp/internal/model"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

func normalize(root ast.Node, src []byte) model.Document {
	var doc model.Document
	for n := root.FirstChild(); n != nil; n = n.NextSibling() {
		if b := toBlock(n, src); b != nil {
			doc.Blocks = append(doc.Blocks, b)
		}
	}
	return doc
}

func toBlock(n ast.Node, src []byte) model.Block {
	switch v := n.(type) {
	case *ast.Heading:
		return &model.Heading{
			Level: v.Level,
			Spans: childSpans(v, src),
		}
	case *ast.Paragraph:
		return &model.Paragraph{
			Spans: childSpans(v, src),
		}
	case *ast.FencedCodeBlock:
		lang := string(v.Language(src))
		var lines []string
		l := v.Lines()
		for i := 0; i < l.Len(); i++ {
			seg := l.At(i)
			lines = append(lines, string(bytes.TrimRight(seg.Value(src), "\n")))
		}
		return &model.CodeBlock{Lang: lang, Lines: lines}
	case *ast.CodeBlock:
		var lines []string
		l := v.Lines()
		for i := 0; i < l.Len(); i++ {
			seg := l.At(i)
			lines = append(lines, string(bytes.TrimRight(seg.Value(src), "\n")))
		}
		return &model.CodeBlock{Lines: lines}
	case *ast.Blockquote:
		return &model.Quote{Blocks: childBlocks(v, src)}
	case *ast.List:
		return toList(v, src)
	case *extast.Table:
		return toTable(v, src)
	case *ast.ThematicBreak:
		return &model.ThematicBreak{}
	case *ast.HTMLBlock:
		raw := nodeRawText(v, src)
		return &model.Paragraph{Spans: []model.Span{&model.Text{Value: raw}}}
	default:
		raw := nodeRawText(n, src)
		if raw == "" {
			return nil
		}
		return &model.Paragraph{Spans: []model.Span{&model.Text{Value: raw}}}
	}
}

func toList(n *ast.List, src []byte) *model.List {
	l := &model.List{Ordered: n.IsOrdered()}
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if li, ok := child.(*ast.ListItem); ok {
			item := model.ListItem{Blocks: childBlocks(li, src)}
			if li.FirstChild() != nil {
				if taskItem, ok := li.FirstChild().(*extast.TaskCheckBox); ok {
					checked := taskItem.IsChecked
					item.Checked = &checked
				}
			}
			l.Items = append(l.Items, item)
		}
	}
	return l
}

func toTable(n *extast.Table, src []byte) *model.Table {
	t := &model.Table{}
	for align := range n.Alignments {
		_ = align
		switch n.Alignments[align] {
		case extast.AlignLeft:
			t.Align = append(t.Align, model.AlignLeft)
		case extast.AlignRight:
			t.Align = append(t.Align, model.AlignRight)
		case extast.AlignCenter:
			t.Align = append(t.Align, model.AlignCenter)
		default:
			t.Align = append(t.Align, model.AlignNone)
		}
	}
	for row := n.FirstChild(); row != nil; row = row.NextSibling() {
		switch r := row.(type) {
		case *extast.TableHeader:
			for cell := r.FirstChild(); cell != nil; cell = cell.NextSibling() {
				t.Headers = append(t.Headers, childSpans(cell, src))
			}
		case *extast.TableRow:
			var cells [][]model.Span
			for cell := r.FirstChild(); cell != nil; cell = cell.NextSibling() {
				cells = append(cells, childSpans(cell, src))
			}
			t.Rows = append(t.Rows, cells...)
		}
	}
	return t
}

func childBlocks(n ast.Node, src []byte) []model.Block {
	var blocks []model.Block
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if b := toBlock(child, src); b != nil {
			blocks = append(blocks, b)
		}
	}
	return blocks
}

func childSpans(n ast.Node, src []byte) []model.Span {
	var spans []model.Span
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if s := toSpan(child, src); s != nil {
			spans = append(spans, s)
		}
	}
	return spans
}

func toSpan(n ast.Node, src []byte) model.Span {
	switch v := n.(type) {
	case *ast.Text:
		seg := v.Segment
		val := string(seg.Value(src))
		if v.SoftLineBreak() {
			val += " "
		}
		if v.HardLineBreak() {
			return &model.HardBreak{}
		}
		return &model.Text{Value: val}
	case *ast.String:
		return &model.Text{Value: string(v.Value)}
	case *ast.Emphasis:
		if v.Level == 2 {
			return &model.Strong{Children: childSpans(v, src)}
		}
		return &model.Emphasis{Children: childSpans(v, src)}
	case *ast.CodeSpan:
		return &model.InlineCode{Value: string(v.Text(src))}
	case *ast.Link:
		return &model.Link{
			Label: childSpans(v, src),
			URL:   string(v.Destination),
		}
	case *ast.AutoLink:
		url := string(v.URL(src))
		return &model.Link{
			Label: []model.Span{&model.Text{Value: url}},
			URL:   url,
		}
	case *ast.Image:
		return &model.Text{Value: "[image: " + nodeRawText(v, src) + "]"}
	case *extast.Strikethrough:
		return &model.Emphasis{Children: childSpans(v, src)}
	default:
		raw := nodeRawText(n, src)
		if raw == "" {
			return nil
		}
		return &model.Text{Value: raw}
	}
}

func nodeRawText(n ast.Node, src []byte) string {
	var buf bytes.Buffer
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if t, ok := child.(*ast.Text); ok {
			seg := t.Segment
			buf.Write(seg.Value(src))
		}
	}
	if buf.Len() == 0 {
		if lines := n.Lines(); lines != nil {
			for i := 0; i < lines.Len(); i++ {
				seg := lines.At(i)
				buf.Write(seg.Value(src))
			}
		}
	}
	return buf.String()
}