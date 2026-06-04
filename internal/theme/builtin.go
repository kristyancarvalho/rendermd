package theme

var Default = Theme{
	Background: "#1e1e2e",
	Text:       "#cdd6f4",
	Muted:      "#7f849c",
	Heading:    "#89b4fa",
	Accent:     "#a6e3a1",
	Link:       "#74c7ec",
	LinkURL:    "#7f849c",
	CodeBg:     "#313244",
	QuoteBg:    "#45475a",
	Border:     "#585b70",

	SyntaxKeyword:  "#cba6f7",
	SyntaxString:   "#a6e3a1",
	SyntaxComment:  "#6c7086",
	SyntaxNumber:   "#fab387",
	SyntaxType:     "#89dceb",
	SyntaxBuiltin:  "#f38ba8",
	SyntaxOperator: "#89b4fa",
}

var Light = Theme{
	Background: "#fafafa",
	Text:       "#1a1a1a",
	Muted:      "#888888",
	Heading:    "#1a56db",
	Accent:     "#2a7a2a",
	Link:       "#1a56db",
	LinkURL:    "#888888",
	CodeBg:     "#f0f0f0",
	QuoteBg:    "#f5f5f5",
	Border:     "#cccccc",

	SyntaxKeyword:  "#7c3aed",
	SyntaxString:   "#166534",
	SyntaxComment:  "#9ca3af",
	SyntaxNumber:   "#b45309",
	SyntaxType:     "#0369a1",
	SyntaxBuiltin:  "#dc2626",
	SyntaxOperator: "#374151",
}

func Resolve(name string) Theme {
	switch name {
	case "light":
		return Light
	default:
		return Default
	}
}

func Merge(base Theme, overrides Theme) Theme {
	if overrides.Background != "" {
		base.Background = overrides.Background
	}
	if overrides.Text != "" {
		base.Text = overrides.Text
	}
	if overrides.Muted != "" {
		base.Muted = overrides.Muted
	}
	if overrides.Heading != "" {
		base.Heading = overrides.Heading
	}
	if overrides.Accent != "" {
		base.Accent = overrides.Accent
	}
	if overrides.Link != "" {
		base.Link = overrides.Link
	}
	if overrides.LinkURL != "" {
		base.LinkURL = overrides.LinkURL
	}
	if overrides.CodeBg != "" {
		base.CodeBg = overrides.CodeBg
	}
	if overrides.QuoteBg != "" {
		base.QuoteBg = overrides.QuoteBg
	}
	if overrides.Border != "" {
		base.Border = overrides.Border
	}
	if overrides.SyntaxKeyword != "" {
		base.SyntaxKeyword = overrides.SyntaxKeyword
	}
	if overrides.SyntaxString != "" {
		base.SyntaxString = overrides.SyntaxString
	}
	if overrides.SyntaxComment != "" {
		base.SyntaxComment = overrides.SyntaxComment
	}
	if overrides.SyntaxNumber != "" {
		base.SyntaxNumber = overrides.SyntaxNumber
	}
	if overrides.SyntaxType != "" {
		base.SyntaxType = overrides.SyntaxType
	}
	if overrides.SyntaxBuiltin != "" {
		base.SyntaxBuiltin = overrides.SyntaxBuiltin
	}
	if overrides.SyntaxOperator != "" {
		base.SyntaxOperator = overrides.SyntaxOperator
	}
	return base
}
