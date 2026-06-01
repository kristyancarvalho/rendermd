package model

type Span interface {
	spanNode()
}

type Text struct {
	Value string
}

type Emphasis struct {
	Children []Span
}

type Strong struct {
	Children []Span
}

type InlineCode struct {
	Value string
}

type Link struct {
	Label []Span
	URL   string
}

type HardBreak struct{}

func (t *Text) spanNode()       {}
func (e *Emphasis) spanNode()   {}
func (s *Strong) spanNode()     {}
func (i *InlineCode) spanNode() {}
func (l *Link) spanNode()       {}
func (h *HardBreak) spanNode()  {}
