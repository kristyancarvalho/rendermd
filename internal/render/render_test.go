package render

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/kristyancarvalho/mdp/internal/layout"
	"github.com/kristyancarvalho/mdp/internal/theme"
)

func testStyles() map[layout.StyleID]lipgloss.Style {
	return map[layout.StyleID]lipgloss.Style{
		layout.StyleNormal:    lipgloss.NewStyle(),
		layout.StyleHeading1:  lipgloss.NewStyle(),
		layout.StyleCodeBlock: lipgloss.NewStyle(),
		layout.StyleMuted:     lipgloss.NewStyle(),
	}
}

func lines(segs ...string) []layout.Line {
	var out []layout.Line
	for _, s := range segs {
		out = append(out, layout.Line{
			Segments: []layout.Segment{{Text: s, Style: layout.StyleNormal}},
		})
	}
	return out
}

func TestRender_FullWindow(t *testing.T) {
	ls := lines("line0", "line1", "line2")
	vp := Viewport{Width: 40, Height: 3, Offset: 0}
	out := Render(ls, testStyles(), vp)
	for _, want := range []string{"line0", "line1", "line2"} {
		if !strings.Contains(out, want) {
			t.Errorf("want %q in output, got:\n%s", want, out)
		}
	}
}

func TestRender_Offset(t *testing.T) {
	ls := lines("line0", "line1", "line2", "line3")
	vp := Viewport{Width: 40, Height: 2, Offset: 2}
	out := Render(ls, testStyles(), vp)
	if strings.Contains(out, "line0") || strings.Contains(out, "line1") {
		t.Error("offset=2 should skip first two lines")
	}
	if !strings.Contains(out, "line2") {
		t.Errorf("want 'line2' in output, got:\n%s", out)
	}
}

func TestRender_OffsetPastEnd(t *testing.T) {
	ls := lines("only")
	vp := Viewport{Width: 40, Height: 5, Offset: 99}
	out := Render(ls, testStyles(), vp)	
	_ = out
}

func TestRender_EmptyLines(t *testing.T) {
	vp := Viewport{Width: 40, Height: 3, Offset: 0}
	out := Render(nil, testStyles(), vp)
	
	if strings.TrimSpace(out) != "" {
		t.Errorf("empty line list should produce blank output, got %q", out)
	}
}

func TestRender_HeightPadding(t *testing.T) {
	ls := lines("only one line")
	vp := Viewport{Width: 40, Height: 5, Offset: 0}
	out := Render(ls, testStyles(), vp)
	count := strings.Count(out, "\n")
	if count < 4 {
		t.Errorf("output should be padded to height with newlines, got %d newlines", count)
	}
}

func TestNew_AllStylesInitialised(t *testing.T) {
	r := New(theme.Default)
	if r == nil {
		t.Fatal("New should not return nil")
	}
	
	ids := []layout.StyleID{
		layout.StyleNormal,
		layout.StyleHeading1,
		layout.StyleHeading2,
		layout.StyleHeading3,
		layout.StyleStrong,
		layout.StyleEmphasis,
		layout.StyleInlineCode,
		layout.StyleCodeBlock,
		layout.StyleQuote,
		layout.StyleLink,
		layout.StyleMuted,
		layout.StyleRule,
	}
	for _, id := range ids {
		if _, ok := r.styles[id]; !ok {
			t.Errorf("renderer missing style for StyleID %d", id)
		}
	}
}

func TestNew_LightTheme(t *testing.T) {
	r := New(theme.Light)
	if r == nil {
		t.Fatal("New should not return nil for light theme")
	}
}

func TestRenderer_Render_DoesNotPanic(t *testing.T) {
	r := New(theme.Default)
	ls := lines("hello", "world")
	vp := Viewport{Width: 80, Height: 10, Offset: 0}
	out := r.Render(ls, vp)
	if !strings.Contains(out, "hello") {
		t.Errorf("rendered output should contain 'hello', got: %q", out)
	}
}

func TestRender_Indent(t *testing.T) {
	ls := []layout.Line{
		{
			Indent:   4,
			Segments: []layout.Segment{{Text: "indented", Style: layout.StyleNormal}},
		},
	}
	vp := Viewport{Width: 40, Height: 1, Offset: 0}
	out := Render(ls, testStyles(), vp)
	if !strings.HasPrefix(out, "    ") {
		t.Errorf("expected 4-space indent prefix, got: %q", out)
	}
}

func TestRender_MultiSegmentLine(t *testing.T) {
	ls := []layout.Line{{
		Segments: []layout.Segment{
			{Text: "hello", Style: layout.StyleNormal},
			{Text: " world", Style: layout.StyleMuted},
		},
	}}
	vp := Viewport{Width: 40, Height: 1, Offset: 0}
	out := Render(ls, testStyles(), vp)
	if !strings.Contains(out, "hello") || !strings.Contains(out, "world") {
		t.Errorf("multi-segment line should contain all segments, got: %q", out)
	}
}