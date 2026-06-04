package render

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/kristyancarvalho/rendermd/internal/layout"
	"github.com/kristyancarvalho/rendermd/internal/theme"
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
		layout.StyleSyntaxKeyword,
		layout.StyleSyntaxString,
		layout.StyleSyntaxComment,
		layout.StyleSyntaxNumber,
		layout.StyleSyntaxType,
		layout.StyleSyntaxBuiltin,
		layout.StyleSyntaxOperator,
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

func TestRender_SyntaxStyles_DoNotPanic(t *testing.T) {
	r := New(theme.Default)
	ls := []layout.Line{{
		Segments: []layout.Segment{
			{Text: "func", Style: layout.StyleSyntaxKeyword},
			{Text: " ", Style: layout.StyleSyntaxOperator},
			{Text: "main", Style: layout.StyleNormal},
			{Text: `"hello"`, Style: layout.StyleSyntaxString},
			{Text: "42", Style: layout.StyleSyntaxNumber},
			{Text: "string", Style: layout.StyleSyntaxType},
			{Text: "len", Style: layout.StyleSyntaxBuiltin},
			{Text: "// note", Style: layout.StyleSyntaxComment},
		},
	}}
	vp := Viewport{Width: 80, Height: 1, Offset: 0}
	out := r.Render(ls, vp)
	if !strings.Contains(out, "func") {
		t.Errorf("syntax keyword should appear in rendered output, got: %q", out)
	}
}

func TestSanitize_StripNullByte(t *testing.T) {
	if got := Sanitize("abc\x00def"); got != "abcdef" {
		t.Errorf("want %q, got %q", "abcdef", got)
	}
}

func TestSanitize_StripESC(t *testing.T) {
	got := Sanitize("\x1b[31mred\x1b[0m")
	if strings.ContainsRune(got, '\x1b') {
		t.Errorf("Sanitize should strip ESC, got %q", got)
	}
	if !strings.Contains(got, "red") {
		t.Errorf("Sanitize should preserve printable text, got %q", got)
	}
}

func TestSanitize_StripC0Controls(t *testing.T) {
	input := "a\x01\x02\x03\x04\x05\x06\x07\x08\x0b\x0c\x0e\x0f\x10\x1fb"
	got := Sanitize(input)
	if got != "ab" {
		t.Errorf("want %q, got %q", "ab", got)
	}
}

func TestSanitize_StripDEL(t *testing.T) {
	if got := Sanitize("abc\x7fdef"); got != "abcdef" {
		t.Errorf("want %q, got %q", "abcdef", got)
	}
}

func TestSanitize_StripC1Controls(t *testing.T) {
	got := Sanitize("a\u0080\u009bb")
	if got != "ab" {
		t.Errorf("C1 controls should be stripped, want %q, got %q", "ab", got)
	}
}

func TestSanitize_PreservesTab(t *testing.T) {
	if got := Sanitize("a\tb"); got != "a\tb" {
		t.Errorf("Sanitize should preserve tab, got %q", got)
	}
}

func TestSanitize_PreservesValidASCII(t *testing.T) {
	cases := []string{
		"hello world",
		"func main() { fmt.Println(\"hi\") }",
		"https://example.com/path?q=1&r=2",
		"# Heading — Paragraph text.",
	}
	for _, tc := range cases {
		if got := Sanitize(tc); got != tc {
			t.Errorf("Sanitize changed valid ASCII: want %q, got %q", tc, got)
		}
	}
}

func TestSanitize_PreservesUnicode(t *testing.T) {
	cases := []string{
		"café résumé",
		"你好世界",
		"▎ blockquote",
		"─────────",
	}
	for _, tc := range cases {
		if got := Sanitize(tc); got != tc {
			t.Errorf("Sanitize changed valid Unicode: want %q, got %q", tc, got)
		}
	}
}

func TestSanitize_EmptyString(t *testing.T) {
	if got := Sanitize(""); got != "" {
		t.Errorf("Sanitize empty string: want %q, got %q", "", got)
	}
}

func TestRender_ESCSequenceStripped(t *testing.T) {
	ls := []layout.Line{{
		Segments: []layout.Segment{
			{Text: "hello\x1b[31m world", Style: layout.StyleNormal},
		},
	}}
	vp := Viewport{Width: 80, Height: 1, Offset: 0}
	out := Render(ls, testStyles(), vp)
	if strings.ContainsRune(out, '\x1b') {
		t.Errorf("rendered output must not contain ESC, got: %q", out)
	}
	if !strings.Contains(out, "hello") || !strings.Contains(out, "world") {
		t.Errorf("rendered output should contain printable text, got: %q", out)
	}
}

func TestRender_NullByteStripped(t *testing.T) {
	ls := []layout.Line{{
		Segments: []layout.Segment{
			{Text: "abc\x00def", Style: layout.StyleNormal},
		},
	}}
	vp := Viewport{Width: 80, Height: 1, Offset: 0}
	out := Render(ls, testStyles(), vp)
	if strings.ContainsRune(out, '\x00') {
		t.Errorf("rendered output must not contain null byte, got: %q", out)
	}
	if !strings.Contains(out, "abc") || !strings.Contains(out, "def") {
		t.Errorf("rendered output should contain printable text, got: %q", out)
	}
}

func TestRender_MultipleControlCharsStripped(t *testing.T) {
	ls := []layout.Line{{
		Segments: []layout.Segment{
			{Text: "\x07bell\x08backspace\x0cformfeed", Style: layout.StyleNormal},
		},
	}}
	vp := Viewport{Width: 80, Height: 1, Offset: 0}
	out := Render(ls, testStyles(), vp)
	for _, r := range []rune{'\x07', '\x08', '\x0c'} {
		if strings.ContainsRune(out, r) {
			t.Errorf("rendered output must not contain control char %q, got: %q", r, out)
		}
	}
	if !strings.Contains(out, "bell") || !strings.Contains(out, "backspace") || !strings.Contains(out, "formfeed") {
		t.Errorf("rendered output should contain printable text, got: %q", out)
	}
}

func TestRender_ValidContentUnchanged(t *testing.T) {
	text := "Hello, 世界! café résumé ▎ ─"
	ls := []layout.Line{{
		Segments: []layout.Segment{{Text: text, Style: layout.StyleNormal}},
	}}
	vp := Viewport{Width: 80, Height: 1, Offset: 0}
	out := Render(ls, testStyles(), vp)
	if !strings.Contains(out, text) {
		t.Errorf("valid Unicode content should be preserved, got: %q", out)
	}
}
