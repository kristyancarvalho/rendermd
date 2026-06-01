package model

type Document struct {
	Blocks []Block
}

type Block interface {
	blockNode()
}

type Heading struct {
	Level int
	Spans []Span
}

type Paragraph struct {
	Spans []Span
}

type CodeBlock struct {
	Lang  string
	Lines []string
}

type Quote struct {
	Blocks []Block
}

type List struct {
	Ordered bool
	Items   []ListItem
}

type ListItem struct {
	Blocks  []Block
	Checked *bool
}

type Table struct {
	Headers [][]Span
	Rows    [][]Span
	Align   []TableAlign
}

type TableAlign int

const (
	AlignNone TableAlign = iota
	AlignLeft
	AlignCenter
	AlignRight
)

type ThematicBreak struct{}

type ImagePlaceholder struct {
	AltText string
	URL     string
}

func (h *Heading) blockNode()          {}
func (p *Paragraph) blockNode()        {}
func (c *CodeBlock) blockNode()        {}
func (q *Quote) blockNode()            {}
func (l *List) blockNode()             {}
func (t *Table) blockNode()            {}
func (r *ThematicBreak) blockNode()    {}
func (i *ImagePlaceholder) blockNode() {}
