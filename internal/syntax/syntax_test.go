package syntax

import (
	"strings"
	"testing"
)

func flatten(tokens []Token) string {
	var sb strings.Builder
	for _, t := range tokens {
		sb.WriteString(t.Text)
	}
	return sb.String()
}

func findKind(tokens []Token, k TokenKind) (Token, bool) {
	for _, t := range tokens {
		if t.Kind == k {
			return t, true
		}
	}
	return Token{}, false
}

func containsKind(tokens []Token, k TokenKind) bool {
	_, ok := findKind(tokens, k)
	return ok
}

func containsText(tokens []Token, text string) bool {
	for _, t := range tokens {
		if strings.Contains(t.Text, text) {
			return true
		}
	}
	return false
}

func TestTokenize_EmptySource(t *testing.T) {
	tok := Tokenize("go", "")
	if len(tok) != 1 {
		t.Fatalf("want 1 token for empty source, got %d", len(tok))
	}
	if tok[0].Text != "" {
		t.Errorf("want empty text, got %q", tok[0].Text)
	}
}

func TestTokenize_UnknownLang_Fallback(t *testing.T) {
	src := "some random text 123"
	tok := Tokenize("brainfuck", src)
	if len(tok) != 1 {
		t.Fatalf("unknown lang: want 1 plain token, got %d", len(tok))
	}
	if tok[0].Kind != KindPlain {
		t.Errorf("unknown lang: want KindPlain, got %v", tok[0].Kind)
	}
	if tok[0].Text != src {
		t.Errorf("unknown lang: want original text, got %q", tok[0].Text)
	}
}

func TestTokenize_EmptyLang_Fallback(t *testing.T) {
	src := "hello world"
	tok := Tokenize("", src)
	if len(tok) != 1 || tok[0].Kind != KindPlain {
		t.Errorf("empty lang: want single plain token, got %v", tok)
	}
}

func TestTokenize_LangCaseInsensitive(t *testing.T) {
	src := "func main() {}"
	tok1 := Tokenize("go", src)
	tok2 := Tokenize("Go", src)
	tok3 := Tokenize("GO", src)
	if flatten(tok1) != flatten(tok2) || flatten(tok1) != flatten(tok3) {
		t.Error("lang matching should be case-insensitive")
	}
}

func TestTokenize_ReconstructsSource(t *testing.T) {
	cases := []struct {
		lang string
		src  string
	}{
		{"go", `func main() { fmt.Println("hello") // comment`},
		{"python", `def foo(x): # comment`},
		{"js", `const x = "hello"; // comment`},
		{"bash", `if [ -f "$file" ]; then`},
		{"json", `{"key": 42, "flag": true}`},
		{"rust", `fn main() { println!("hi"); // comment`},
		{"unknown", `whatever text goes here`},
	}
	for _, c := range cases {
		tok := Tokenize(c.lang, c.src)
		got := flatten(tok)
		if got != c.src {
			t.Errorf("lang=%q: reconstructed %q != original %q", c.lang, got, c.src)
		}
	}
}

func TestGo_Keyword(t *testing.T) {
	tok := Tokenize("go", "func main() {}")
	if !containsKind(tok, KindKeyword) {
		t.Error("expected KindKeyword token for 'func'")
	}

	for _, tk := range tok {
		if tk.Text == "func" && tk.Kind != KindKeyword {
			t.Errorf("'func' should be KindKeyword, got %v", tk.Kind)
		}
	}
}

func TestGo_String(t *testing.T) {
	tok := Tokenize("go", `fmt.Println("hello world")`)
	if !containsText(tok, "hello world") {
		t.Fatal("string content not found in tokens")
	}
	for _, tk := range tok {
		if strings.Contains(tk.Text, "hello") && tk.Kind != KindString {
			t.Errorf("string literal should be KindString, got %v", tk.Kind)
		}
	}
}

func TestGo_Comment(t *testing.T) {
	tok := Tokenize("go", "x := 1 // this is a comment")
	if !containsKind(tok, KindComment) {
		t.Error("expected KindComment token")
	}
	for _, tk := range tok {
		if tk.Kind == KindComment && !strings.Contains(tk.Text, "comment") {
			t.Errorf("comment token should contain 'comment', got %q", tk.Text)
		}
	}
}

func TestGo_Number(t *testing.T) {
	tok := Tokenize("go", "x := 42")
	if !containsKind(tok, KindNumber) {
		t.Error("expected KindNumber token for '42'")
	}
}

func TestGo_Type(t *testing.T) {
	tok := Tokenize("go", "var x string")
	found := false
	for _, tk := range tok {
		if tk.Text == "string" {
			if tk.Kind != KindType {
				t.Errorf("'string' should be KindType, got %v", tk.Kind)
			}
			found = true
		}
	}
	if !found {
		t.Error("'string' type not found in tokens")
	}
}

func TestGo_Builtin(t *testing.T) {
	tok := Tokenize("go", "n := len(s)")
	found := false
	for _, tk := range tok {
		if tk.Text == "len" {
			if tk.Kind != KindBuiltin {
				t.Errorf("'len' should be KindBuiltin, got %v", tk.Kind)
			}
			found = true
		}
	}
	if !found {
		t.Error("'len' builtin not found in tokens")
	}
}

func TestGo_RawString(t *testing.T) {
	tok := Tokenize("go", "s := `raw string`")
	for _, tk := range tok {
		if strings.Contains(tk.Text, "raw string") && tk.Kind != KindString {
			t.Errorf("raw string should be KindString, got %v", tk.Kind)
		}
	}
}

func TestGo_MultipleKeywords(t *testing.T) {
	tok := Tokenize("go", "if err != nil { return err }")
	keywords := 0
	for _, tk := range tok {
		if tk.Kind == KindKeyword {
			keywords++
		}
	}
	if keywords < 2 {
		t.Errorf("expected at least 2 keywords (if, return), got %d", keywords)
	}
}

func TestGo_BlockComment(t *testing.T) {
	tok := Tokenize("go", "/* block comment */ x := 1")
	if !containsKind(tok, KindComment) {
		t.Error("expected block comment to produce KindComment")
	}
}

func TestPython_Keyword(t *testing.T) {
	tok := Tokenize("python", "def foo(x):")
	found := false
	for _, tk := range tok {
		if tk.Text == "def" && tk.Kind == KindKeyword {
			found = true
		}
	}
	if !found {
		t.Error("'def' should be KindKeyword")
	}
}

func TestPython_Comment(t *testing.T) {
	tok := Tokenize("py", "x = 1 # inline comment")
	if !containsKind(tok, KindComment) {
		t.Error("expected KindComment for Python '#' comment")
	}
}

func TestPython_String(t *testing.T) {
	tok := Tokenize("python", `x = "hello"`)
	if !containsKind(tok, KindString) {
		t.Error("expected KindString")
	}
}

func TestPython_Builtin(t *testing.T) {
	tok := Tokenize("python", "print(len(x))")
	builtins := 0
	for _, tk := range tok {
		if tk.Kind == KindBuiltin {
			builtins++
		}
	}
	if builtins < 1 {
		t.Error("expected at least one KindBuiltin (print or len)")
	}
}

func TestJS_Keyword(t *testing.T) {
	tok := Tokenize("js", "const x = 42;")
	found := false
	for _, tk := range tok {
		if tk.Text == "const" && tk.Kind == KindKeyword {
			found = true
		}
	}
	if !found {
		t.Error("'const' should be KindKeyword")
	}
}

func TestTS_Type(t *testing.T) {
	tok := Tokenize("ts", "let x: string = 'hi';")
	found := false
	for _, tk := range tok {
		if tk.Text == "string" && tk.Kind == KindType {
			found = true
		}
	}
	if !found {
		t.Error("'string' should be KindType in TypeScript")
	}
}

func TestJS_TemplateLiteral(t *testing.T) {
	tok := Tokenize("js", "const s = `hello ${name}`;")
	if !containsKind(tok, KindString) {
		t.Error("template literal should produce KindString")
	}
}

func TestJS_Comment(t *testing.T) {
	tok := Tokenize("javascript", "const x = 1; // comment")
	if !containsKind(tok, KindComment) {
		t.Error("expected KindComment for JS '//' comment")
	}
}

func TestBash_Comment(t *testing.T) {
	tok := Tokenize("bash", "# this is a comment")
	if !containsKind(tok, KindComment) {
		t.Error("expected KindComment for bash '#' comment")
	}
}

func TestBash_Keyword(t *testing.T) {
	tok := Tokenize("sh", "if [ -f file ]; then")
	found := false
	for _, tk := range tok {
		if (tk.Text == "if" || tk.Text == "then") && tk.Kind == KindKeyword {
			found = true
		}
	}
	if !found {
		t.Error("expected KindKeyword for bash 'if'/'then'")
	}
}

func TestJSON_String(t *testing.T) {
	tok := Tokenize("json", `"name": "Alice"`)
	strings_ := 0
	for _, tk := range tok {
		if tk.Kind == KindString {
			strings_++
		}
	}
	if strings_ < 2 {
		t.Errorf("expected at least 2 string tokens in JSON, got %d", strings_)
	}
}

func TestJSON_Number(t *testing.T) {
	tok := Tokenize("json", `"age": 30`)
	if !containsKind(tok, KindNumber) {
		t.Error("expected KindNumber for JSON number")
	}
}

func TestJSON_Keyword(t *testing.T) {
	tok := Tokenize("json", `"flag": true`)
	found := false
	for _, tk := range tok {
		if tk.Text == "true" && tk.Kind == KindKeyword {
			found = true
		}
	}
	if !found {
		t.Error("'true' should be KindKeyword in JSON")
	}
}

func TestRust_Keyword(t *testing.T) {
	tok := Tokenize("rust", "fn main() {}")
	found := false
	for _, tk := range tok {
		if tk.Text == "fn" && tk.Kind == KindKeyword {
			found = true
		}
	}
	if !found {
		t.Error("'fn' should be KindKeyword in Rust")
	}
}

func TestRust_Comment(t *testing.T) {
	tok := Tokenize("rs", "let x = 1; // comment")
	if !containsKind(tok, KindComment) {
		t.Error("expected KindComment in Rust")
	}
}

func TestRust_MacroBuiltin(t *testing.T) {
	tok := Tokenize("rust", `println!("hello");`)
	found := false
	for _, tk := range tok {
		if strings.HasPrefix(tk.Text, "println") && tk.Kind == KindBuiltin {
			found = true
		}
	}
	if !found {
		t.Error("'println!' should be KindBuiltin in Rust")
	}
}

func TestRust_Type(t *testing.T) {
	tok := Tokenize("rust", "let x: String = String::new();")
	found := false
	for _, tk := range tok {
		if tk.Text == "String" && tk.Kind == KindType {
			found = true
		}
	}
	if !found {
		t.Error("'String' should be KindType in Rust")
	}
}
