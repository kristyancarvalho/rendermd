package syntax

import "strings"

type TokenKind int

const (
	KindPlain TokenKind = iota
	KindKeyword
	KindString
	KindComment
	KindNumber
	KindType
	KindBuiltin
	KindOperator
)

type Token struct {
	Text string
	Kind TokenKind
}

func Tokenize(lang, src string) []Token {
	if src == "" {
		return []Token{{Text: "", Kind: KindPlain}}
	}
	switch strings.ToLower(strings.TrimSpace(lang)) {
	case "go":
		return tokenizeGo(src)
	case "python", "py":
		return tokenizePython(src)
	case "javascript", "js", "typescript", "ts":
		return tokenizeJS(src)
	case "bash", "sh", "shell", "zsh":
		return tokenizeBash(src)
	case "json":
		return tokenizeJSON(src)
	case "rust", "rs":
		return tokenizeRust(src)
	default:
		return []Token{{Text: src, Kind: KindPlain}}
	}
}

type scanner struct {
	src []byte
	pos int
	out []Token
}

func newScanner(src string) *scanner {
	return &scanner{src: []byte(src)}
}

func (s *scanner) done() bool { return s.pos >= len(s.src) }
func (s *scanner) peek() byte {
	if s.pos >= len(s.src) {
		return 0
	}
	return s.src[s.pos]
}
func (s *scanner) peekAt(offset int) byte {
	i := s.pos + offset
	if i >= len(s.src) {
		return 0
	}
	return s.src[i]
}
func (s *scanner) advance() byte {
	b := s.src[s.pos]
	s.pos++
	return b
}
func (s *scanner) emit(start int, kind TokenKind) {
	text := string(s.src[start:s.pos])
	if text == "" {
		return
	}

	if len(s.out) > 0 && s.out[len(s.out)-1].Kind == kind {
		s.out[len(s.out)-1].Text += text
		return
	}
	s.out = append(s.out, Token{Text: text, Kind: kind})
}

func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}
func isDigit(b byte) bool { return b >= '0' && b <= '9' }
func isAlNum(b byte) bool { return isLetter(b) || isDigit(b) }
func isHexDigit(b byte) bool {
	return isDigit(b) || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
}

func (s *scanner) consumeString(q byte) {
	for !s.done() {
		b := s.advance()
		if b == '\\' && !s.done() {
			s.advance()
			continue
		}
		if b == q {
			return
		}
	}
}

func (s *scanner) consumeNumber() {
	if s.peek() == 'x' || s.peek() == 'X' {
		s.advance()
		for !s.done() && (isHexDigit(s.peek()) || s.peek() == '_') {
			s.advance()
		}
		return
	}
	for !s.done() && (isDigit(s.peek()) || s.peek() == '.' || s.peek() == '_' ||
		s.peek() == 'e' || s.peek() == 'E' || s.peek() == '+' || s.peek() == '-') {
		s.advance()
	}
}

var goKeywords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true,
	"interface": true, "map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true, "var": true,
}

var goTypes = map[string]bool{
	"bool": true, "byte": true, "complex64": true, "complex128": true,
	"error": true, "float32": true, "float64": true,
	"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
	"rune": true, "string": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"uintptr": true,
}

var goBuiltins = map[string]bool{
	"append": true, "cap": true, "close": true, "complex": true, "copy": true,
	"delete": true, "imag": true, "len": true, "make": true, "new": true,
	"panic": true, "print": true, "println": true, "real": true, "recover": true,
	"true": true, "false": true, "nil": true, "iota": true,
}

func tokenizeGo(src string) []Token {
	s := newScanner(src)
	for !s.done() {
		start := s.pos
		b := s.peek()

		if b == '/' && s.peekAt(1) == '/' {
			s.pos = len(s.src)
			s.emit(start, KindComment)
			continue
		}

		if b == '/' && s.peekAt(1) == '*' {
			s.advance()
			s.advance()
			for !s.done() {
				if s.peek() == '*' && s.peekAt(1) == '/' {
					s.advance()
					s.advance()
					break
				}
				s.advance()
			}
			s.emit(start, KindComment)
			continue
		}

		if b == '"' || b == '\'' || b == '`' {
			q := s.advance()
			if b == '`' {
				for !s.done() && s.peek() != '`' {
					s.advance()
				}
				if !s.done() {
					s.advance()
				}
			} else {
				s.consumeString(q)
			}
			s.emit(start, KindString)
			continue
		}

		if isDigit(b) || (b == '.' && isDigit(s.peekAt(1))) {
			s.advance()
			s.consumeNumber()
			s.emit(start, KindNumber)
			continue
		}

		if isLetter(b) {
			for !s.done() && isAlNum(s.peek()) {
				s.advance()
			}
			word := string(s.src[start:s.pos])
			switch {
			case goKeywords[word]:
				s.emit(start, KindKeyword)
			case goTypes[word]:
				s.emit(start, KindType)
			case goBuiltins[word]:
				s.emit(start, KindBuiltin)
			default:
				s.emit(start, KindPlain)
			}
			continue
		}

		s.advance()
		s.emit(start, KindOperator)
	}
	return s.out
}

var pythonKeywords = map[string]bool{
	"False": true, "None": true, "True": true,
	"and": true, "as": true, "assert": true, "async": true, "await": true,
	"break": true, "class": true, "continue": true, "def": true, "del": true,
	"elif": true, "else": true, "except": true, "finally": true, "for": true,
	"from": true, "global": true, "if": true, "import": true, "in": true,
	"is": true, "lambda": true, "nonlocal": true, "not": true, "or": true,
	"pass": true, "raise": true, "return": true, "try": true, "while": true,
	"with": true, "yield": true,
}

var pythonBuiltins = map[string]bool{
	"abs": true, "all": true, "any": true, "bin": true, "bool": true,
	"bytes": true, "callable": true, "chr": true, "dict": true, "dir": true,
	"divmod": true, "enumerate": true, "eval": true, "exec": true, "filter": true,
	"float": true, "format": true, "frozenset": true, "getattr": true,
	"globals": true, "hasattr": true, "hash": true, "help": true, "hex": true,
	"id": true, "input": true, "int": true, "isinstance": true, "issubclass": true,
	"iter": true, "len": true, "list": true, "locals": true, "map": true,
	"max": true, "min": true, "next": true, "object": true, "oct": true,
	"open": true, "ord": true, "pow": true, "print": true, "property": true,
	"range": true, "repr": true, "reversed": true, "round": true, "set": true,
	"setattr": true, "slice": true, "sorted": true, "staticmethod": true,
	"str": true, "sum": true, "super": true, "tuple": true, "type": true,
	"vars": true, "zip": true,
}

func tokenizePython(src string) []Token {
	s := newScanner(src)
	for !s.done() {
		start := s.pos
		b := s.peek()

		if b == '#' {
			s.pos = len(s.src)
			s.emit(start, KindComment)
			continue
		}

		if (b == '"' || b == '\'') && s.peekAt(1) == b && s.peekAt(2) == b {
			s.advance()
			s.advance()
			s.advance()
			q := b
			for !s.done() {
				if s.peek() == q && s.peekAt(1) == q && s.peekAt(2) == q {
					s.advance()
					s.advance()
					s.advance()
					break
				}
				s.advance()
			}
			s.emit(start, KindString)
			continue
		}

		if (b == 'f' || b == 'F' || b == 'r' || b == 'R' || b == 'b' || b == 'B') &&
			(s.peekAt(1) == '"' || s.peekAt(1) == '\'') {
			s.advance()
			q := s.advance()
			s.consumeString(q)
			s.emit(start, KindString)
			continue
		}

		if b == '"' || b == '\'' {
			q := s.advance()
			s.consumeString(q)
			s.emit(start, KindString)
			continue
		}

		if isDigit(b) || (b == '.' && isDigit(s.peekAt(1))) {
			s.advance()
			s.consumeNumber()
			s.emit(start, KindNumber)
			continue
		}

		if isLetter(b) {
			for !s.done() && isAlNum(s.peek()) {
				s.advance()
			}
			word := string(s.src[start:s.pos])
			switch {
			case pythonKeywords[word]:
				s.emit(start, KindKeyword)
			case pythonBuiltins[word]:
				s.emit(start, KindBuiltin)
			default:
				s.emit(start, KindPlain)
			}
			continue
		}

		s.advance()
		s.emit(start, KindOperator)
	}
	return s.out
}

var jsKeywords = map[string]bool{
	"break": true, "case": true, "catch": true, "class": true, "const": true,
	"continue": true, "debugger": true, "default": true, "delete": true,
	"do": true, "else": true, "export": true, "extends": true, "false": true,
	"finally": true, "for": true, "function": true, "if": true, "import": true,
	"in": true, "instanceof": true, "let": true, "new": true, "null": true,
	"of": true, "return": true, "static": true, "super": true, "switch": true,
	"this": true, "throw": true, "true": true, "try": true, "typeof": true,
	"undefined": true, "var": true, "void": true, "while": true, "with": true,
	"yield": true, "async": true, "await": true,
	"abstract": true, "as": true, "declare": true, "enum": true, "from": true,
	"implements": true, "interface": true, "namespace": true, "never": true,
	"readonly": true, "type": true,
}

var jsTypes = map[string]bool{
	"any": true, "boolean": true, "number": true, "object": true, "string": true,
	"symbol": true, "unknown": true,
}

func tokenizeJS(src string) []Token {
	s := newScanner(src)
	for !s.done() {
		start := s.pos
		b := s.peek()

		if b == '/' && s.peekAt(1) == '/' {
			s.pos = len(s.src)
			s.emit(start, KindComment)
			continue
		}

		if b == '/' && s.peekAt(1) == '*' {
			s.advance()
			s.advance()
			for !s.done() {
				if s.peek() == '*' && s.peekAt(1) == '/' {
					s.advance()
					s.advance()
					break
				}
				s.advance()
			}
			s.emit(start, KindComment)
			continue
		}

		if b == '`' {
			s.advance()
			for !s.done() && s.peek() != '`' {
				if s.peek() == '\\' {
					s.advance()
				}
				if !s.done() {
					s.advance()
				}
			}
			if !s.done() {
				s.advance()
			}
			s.emit(start, KindString)
			continue
		}

		if b == '"' || b == '\'' {
			q := s.advance()
			s.consumeString(q)
			s.emit(start, KindString)
			continue
		}

		if isDigit(b) || (b == '.' && isDigit(s.peekAt(1))) {
			s.advance()
			s.consumeNumber()
			s.emit(start, KindNumber)
			continue
		}

		if isLetter(b) {
			for !s.done() && isAlNum(s.peek()) {
				s.advance()
			}
			word := string(s.src[start:s.pos])
			switch {
			case jsKeywords[word]:
				s.emit(start, KindKeyword)
			case jsTypes[word]:
				s.emit(start, KindType)
			default:
				s.emit(start, KindPlain)
			}
			continue
		}

		s.advance()
		s.emit(start, KindOperator)
	}
	return s.out
}

var bashKeywords = map[string]bool{
	"case": true, "do": true, "done": true, "elif": true, "else": true,
	"esac": true, "fi": true, "for": true, "function": true, "if": true,
	"in": true, "select": true, "then": true, "until": true, "while": true,
	"return": true, "exit": true, "export": true, "local": true, "readonly": true,
	"source": true, "unset": true,
}

func tokenizeBash(src string) []Token {
	s := newScanner(src)
	for !s.done() {
		start := s.pos
		b := s.peek()

		if b == '#' {
			s.pos = len(s.src)
			s.emit(start, KindComment)
			continue
		}

		if b == '"' || b == '\'' {
			q := s.advance()
			s.consumeString(q)
			s.emit(start, KindString)
			continue
		}

		if isDigit(b) {
			s.advance()
			s.consumeNumber()
			s.emit(start, KindNumber)
			continue
		}

		if isLetter(b) {
			for !s.done() && (isAlNum(s.peek()) || s.peek() == '-') {
				s.advance()
			}
			word := string(s.src[start:s.pos])
			if bashKeywords[word] {
				s.emit(start, KindKeyword)
			} else {
				s.emit(start, KindPlain)
			}
			continue
		}

		s.advance()
		s.emit(start, KindOperator)
	}
	return s.out
}

func tokenizeJSON(src string) []Token {
	s := newScanner(src)
	for !s.done() {
		start := s.pos
		b := s.peek()

		if b == '"' {
			s.advance()
			s.consumeString('"')
			s.emit(start, KindString)
			continue
		}

		if isDigit(b) || b == '-' {
			s.advance()
			s.consumeNumber()
			s.emit(start, KindNumber)
			continue
		}

		if isLetter(b) {
			for !s.done() && isLetter(s.peek()) {
				s.advance()
			}
			word := string(s.src[start:s.pos])
			if word == "true" || word == "false" || word == "null" {
				s.emit(start, KindKeyword)
			} else {
				s.emit(start, KindPlain)
			}
			continue
		}

		s.advance()
		s.emit(start, KindOperator)
	}
	return s.out
}

var rustKeywords = map[string]bool{
	"as": true, "async": true, "await": true, "break": true, "const": true,
	"continue": true, "crate": true, "dyn": true, "else": true, "enum": true,
	"extern": true, "false": true, "fn": true, "for": true, "if": true,
	"impl": true, "in": true, "let": true, "loop": true, "match": true,
	"mod": true, "move": true, "mut": true, "pub": true, "ref": true,
	"return": true, "self": true, "Self": true, "static": true, "struct": true,
	"super": true, "trait": true, "true": true, "type": true, "unsafe": true,
	"use": true, "where": true, "while": true,
}

var rustTypes = map[string]bool{
	"bool": true, "char": true, "f32": true, "f64": true,
	"i8": true, "i16": true, "i32": true, "i64": true, "i128": true, "isize": true,
	"str": true, "String": true,
	"u8": true, "u16": true, "u32": true, "u64": true, "u128": true, "usize": true,
}

var rustBuiltins = map[string]bool{
	"Some": true, "None": true, "Ok": true, "Err": true,
	"Box": true, "Vec": true, "Option": true, "Result": true,
	"println": true, "print": true, "eprintln": true, "eprint": true,
	"format": true, "panic": true, "todo": true, "unimplemented": true,
	"unreachable": true, "assert": true, "assert_eq": true, "assert_ne": true,
	"dbg": true, "vec": true,
}

func tokenizeRust(src string) []Token {
	s := newScanner(src)
	for !s.done() {
		start := s.pos
		b := s.peek()

		if b == '/' && (s.peekAt(1) == '/' || s.peekAt(1) == '!') {
			s.pos = len(s.src)
			s.emit(start, KindComment)
			continue
		}

		if b == '/' && s.peekAt(1) == '*' {
			s.advance()
			s.advance()
			for !s.done() {
				if s.peek() == '*' && s.peekAt(1) == '/' {
					s.advance()
					s.advance()
					break
				}
				s.advance()
			}
			s.emit(start, KindComment)
			continue
		}

		if b == '"' {
			s.advance()
			s.consumeString('"')
			s.emit(start, KindString)
			continue
		}

		if b == '\'' {
			s.advance()
			s.consumeString('\'')
			s.emit(start, KindString)
			continue
		}

		if isDigit(b) {
			s.advance()
			s.consumeNumber()
			s.emit(start, KindNumber)
			continue
		}

		if isLetter(b) {
			for !s.done() && isAlNum(s.peek()) {
				s.advance()
			}

			if !s.done() && s.peek() == '!' {
				s.advance()
			}
			word := string(s.src[start:s.pos])
			lookup := strings.TrimSuffix(word, "!")
			switch {
			case rustKeywords[lookup]:
				s.emit(start, KindKeyword)
			case rustTypes[lookup]:
				s.emit(start, KindType)
			case rustBuiltins[lookup] || strings.HasSuffix(word, "!"):
				s.emit(start, KindBuiltin)
			default:
				s.emit(start, KindPlain)
			}
			continue
		}

		s.advance()
		s.emit(start, KindOperator)
	}
	return s.out
}
