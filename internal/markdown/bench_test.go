package markdown

import (
	"fmt"
	"strings"
	"testing"
)

func makeMarkdown(sections, listsPerSection, codeBlocksPerSection int) []byte {
	var sb strings.Builder
	for i := 0; i < sections; i++ {
		fmt.Fprintf(&sb, "# Section %d\n\n", i+1)
		sb.WriteString("Introductory paragraph with **bold**, _italic_, `code`, and [a link](https://example.com).\n\n")

		for j := 0; j < listsPerSection; j++ {
			fmt.Fprintf(&sb, "## Subsection %d.%d\n\n", i+1, j+1)
			sb.WriteString("- Item one in the unordered list\n")
			sb.WriteString("- Item two with `inline code`\n")
			sb.WriteString("- Item three with **bold text**\n\n")
			sb.WriteString("1. First ordered item\n")
			sb.WriteString("2. Second ordered item\n")
			sb.WriteString("3. Third ordered item\n\n")
			sb.WriteString("- [x] Completed task\n")
			sb.WriteString("- [ ] Pending task\n\n")
			sb.WriteString("| Header A | Header B |\n")
			sb.WriteString("|----------|----------|\n")
			fmt.Fprintf(&sb, "| cell%d-1  | cell%d-2  |\n\n", j, j)
			sb.WriteString("> A blockquote paragraph with enough text to be representative.\n\n")
		}

		for k := 0; k < codeBlocksPerSection; k++ {
			sb.WriteString("```go\n")
			fmt.Fprintf(&sb, "func function%d(x int) (string, error) {\n", k)
			sb.WriteString("\tif x < 0 {\n")
			sb.WriteString("\t\treturn \"\", fmt.Errorf(\"negative: %d\", x)\n")
			sb.WriteString("\t}\n")
			sb.WriteString("\treturn fmt.Sprintf(\"%d\", x*x), nil\n")
			sb.WriteString("}\n")
			sb.WriteString("```\n\n")
		}

		sb.WriteString("---\n\n")
	}
	return []byte(sb.String())
}

func BenchmarkParse_Small(b *testing.B) {
	src := makeMarkdown(5, 1, 1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Parse(src)
	}
}

func BenchmarkParse_Medium(b *testing.B) {
	src := makeMarkdown(50, 2, 2)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Parse(src)
	}
}

func BenchmarkParse_Large(b *testing.B) {
	src := makeMarkdown(200, 3, 3)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Parse(src)
	}
}

func BenchmarkParse_HeavyCode(b *testing.B) {
	src := makeMarkdown(20, 1, 10)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Parse(src)
	}
}

func BenchmarkParse_HeavyTables(b *testing.B) {
	var sb strings.Builder
	for i := 0; i < 200; i++ {
		sb.WriteString("| Col1 | Col2 | Col3 | Col4 | Col5 |\n")
		sb.WriteString("|------|------|------|------|------|\n")
		for r := 0; r < 10; r++ {
			fmt.Fprintf(&sb, "| val%d | val%d | val%d | val%d | val%d |\n", r, r, r, r, r)
		}
		sb.WriteString("\n")
	}
	src := []byte(sb.String())
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Parse(src)
	}
}