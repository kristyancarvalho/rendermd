package render

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kristyancarvalho/mdp/internal/layout"
	"github.com/kristyancarvalho/mdp/internal/theme"
)

type Viewport struct {
	Width  int
	Height int
	Offset int
}

type Renderer struct {
	styles map[layout.StyleID]lipgloss.Style
}

func New(thm theme.Theme) *Renderer {
	r := &Renderer{
		styles: make(map[layout.StyleID]lipgloss.Style),
	}
	r.styles[layout.StyleNormal] = lipgloss.NewStyle().Foreground(lipgloss.Color(thm.Text))
	r.styles[layout.StyleHeading1] = lipgloss.NewStyle().Foreground(lipgloss.Color(thm.Heading)).Bold(true)
	r.styles[layout.StyleHeading2] = lipgloss.NewStyle().Foreground(lipgloss.Color(thm.Heading)).Bold(true)
	r.styles[layout.StyleHeading3] = lipgloss.NewStyle().Foreground(lipgloss.Color(thm.Heading))
	r.styles[layout.StyleHeading4] = lipgloss.NewStyle().Foreground(lipgloss.Color(thm.Heading))
	r.styles[layout.StyleHeading5] = lipgloss.NewStyle().Foreground(lipgloss.Color(thm.Heading))
	r.styles[layout.StyleHeading6] = lipgloss.NewStyle().Foreground(lipgloss.Color(thm.Heading))
	r.styles[layout.StyleStrong] = lipgloss.NewStyle().Bold(true)
	r.styles[layout.StyleEmphasis] = lipgloss.NewStyle().Italic(true)
	r.styles[layout.StyleInlineCode] = lipgloss.NewStyle().
		Foreground(lipgloss.Color(thm.Text)).
		Background(lipgloss.Color(thm.CodeBg))
	r.styles[layout.StyleCodeBlock] = lipgloss.NewStyle().
		Foreground(lipgloss.Color(thm.Text)).
		Background(lipgloss.Color(thm.CodeBg))
	r.styles[layout.StyleQuote] = lipgloss.NewStyle().
		Foreground(lipgloss.Color(thm.Muted)).
		Background(lipgloss.Color(thm.QuoteBg))
	r.styles[layout.StyleLink] = lipgloss.NewStyle().
		Foreground(lipgloss.Color(thm.Link)).
		Underline(true)
	r.styles[layout.StyleLinkURL] = lipgloss.NewStyle().Foreground(lipgloss.Color(thm.LinkURL))
	r.styles[layout.StyleMuted] = lipgloss.NewStyle().Foreground(lipgloss.Color(thm.Muted))
	r.styles[layout.StyleRule] = lipgloss.NewStyle().Foreground(lipgloss.Color(thm.Border))
	r.styles[layout.StyleAccent] = lipgloss.NewStyle().Foreground(lipgloss.Color(thm.Accent)).Bold(true)
	return r
}

func (r *Renderer) Render(lines []layout.Line, vp Viewport) string {
	return Render(lines, r.styles, vp)
}

func Render(lines []layout.Line, styles map[layout.StyleID]lipgloss.Style, vp Viewport) string {
	start := vp.Offset
	end := vp.Offset + vp.Height
	if start < 0 {
		start = 0
	}
	if end > len(lines) {
		end = len(lines)
	}
	if start >= len(lines) {
		return strings.Repeat("\n", vp.Height)
	}

	visible := lines[start:end]
	var sb strings.Builder
	for i, line := range visible {
		indentStr := strings.Repeat(" ", line.Indent)
		if line.Indent > 0 {
			sb.WriteString(indentStr)
		}
		for _, seg := range line.Segments {
			st, ok := styles[seg.Style]
			if !ok {
				st = lipgloss.NewStyle()
			}
			sb.WriteString(st.Render(seg.Text))
		}
		if i < len(visible)-1 {
			sb.WriteByte('\n')
		}
	}

	for len(visible) < vp.Height {
		sb.WriteByte('\n')
		visible = append(visible, layout.Line{})
	}
	return sb.String()
}
