package markdown

import (
	"github.com/kristyancarvalho/rendermd/internal/model"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var gm = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.Typographer,
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

func goldmarkText(src []byte) text.Reader {
	return text.NewReader(src)
}

func Parse(src []byte) model.Document {
	node := gm.Parser().Parse(goldmarkText(src))
	return normalize(node, src)
}
