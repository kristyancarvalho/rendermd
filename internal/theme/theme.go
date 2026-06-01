package theme

type Theme struct {
	Background string
	Text       string
	Muted      string
	Heading    string
	Accent     string
	Link       string
	LinkURL    string
	CodeBg     string
	QuoteBg    string
	Border     string

	SyntaxKeyword  string
	SyntaxString   string
	SyntaxComment  string
	SyntaxNumber   string
	SyntaxType     string
	SyntaxBuiltin  string
	SyntaxOperator string
}