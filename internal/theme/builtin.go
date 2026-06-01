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
	return base
}
