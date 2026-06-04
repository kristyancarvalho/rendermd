package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/kristyancarvalho/mdp/internal/theme"
)

type Config struct {
	UI       UIConfig       `toml:"ui"`
	Theme    ThemeConfig    `toml:"theme"`
	Markdown MarkdownConfig `toml:"markdown"`
	Keys     KeysConfig     `toml:"keys"`
	Watch    WatchConfig    `toml:"watch"`
	resolved theme.Theme
}

type UIConfig struct {
	Padding         int  `toml:"padding"`
	LineSpacing     int  `toml:"line_spacing"`
	Scrolloff       int  `toml:"scrolloff"`
	SoftWrap        bool `toml:"soft_wrap"`
	MaxWidth        int  `toml:"max_width"`
	ShowLineNumbers bool `toml:"show_line_numbers"`
	ShowURLs        bool `toml:"show_urls"`
}

type ThemeConfig struct {
	Name       string `toml:"name"`
	Background string `toml:"background"`
	Text       string `toml:"text"`
	Muted      string `toml:"muted"`
	Heading    string `toml:"heading"`
	Accent     string `toml:"accent"`
	Link       string `toml:"link"`
	LinkURL    string `toml:"link_url"`
	CodeBg     string `toml:"code_bg"`
	QuoteBg    string `toml:"quote_bg"`
	Border     string `toml:"border"`

	SyntaxKeyword  string `toml:"syntax_keyword"`
	SyntaxString   string `toml:"syntax_string"`
	SyntaxComment  string `toml:"syntax_comment"`
	SyntaxNumber   string `toml:"syntax_number"`
	SyntaxType     string `toml:"syntax_type"`
	SyntaxBuiltin  string `toml:"syntax_builtin"`
	SyntaxOperator string `toml:"syntax_operator"`
}

type MarkdownConfig struct {
	HideSyntax      bool `toml:"hide_syntax"`
	RenderEmphasis  bool `toml:"render_emphasis"`
	RenderStrong    bool `toml:"render_strong"`
	RenderLinks     bool `toml:"render_links"`
	RenderImages    bool `toml:"render_images"`
	RenderTables    bool `toml:"render_tables"`
	RenderTaskLists bool `toml:"render_task_lists"`
}

type KeysConfig struct {
	Up       string `toml:"up"`
	Down     string `toml:"down"`
	HalfUp   string `toml:"half_up"`
	HalfDown string `toml:"half_down"`
	Top      string `toml:"top"`
	Bottom   string `toml:"bottom"`
	Search   string `toml:"search"`
	NextHit  string `toml:"next_hit"`
	PrevHit  string `toml:"prev_hit"`
	Reload   string `toml:"reload"`
	Quit     string `toml:"quit"`
	Help     string `toml:"help"`
}

type WatchConfig struct {
	Enabled    bool `toml:"enabled"`
	DebounceMs int  `toml:"debounce_ms"`
}

func (c *Config) ResolvedTheme() theme.Theme {
	return c.resolved
}

func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "mdp", "config.toml")
}

func Load(path string) Config {
	cfg := defaults()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}

	var overlay Config
	meta, err2 := toml.Decode(string(data), &overlay)
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "mdp: config parse error: %v\n", err2)
		return cfg
	}

	cfg = merge(cfg, overlay, meta)

	warnings := validateThemeConfig(&cfg.Theme)
	for _, w := range warnings {
		fmt.Fprintln(os.Stderr, w.String())
	}

	cfg.resolved = resolveTheme(cfg.Theme)
	return cfg
}

func validateThemeConfig(tc *ThemeConfig) []theme.ValidationWarning {
	var warnings []theme.ValidationWarning

	validName, nameWarnings := theme.ValidateThemeName(tc.Name)
	warnings = append(warnings, nameWarnings...)
	tc.Name = validName

	overrides := &theme.Theme{
		Background:     tc.Background,
		Text:           tc.Text,
		Muted:          tc.Muted,
		Heading:        tc.Heading,
		Accent:         tc.Accent,
		Link:           tc.Link,
		LinkURL:        tc.LinkURL,
		CodeBg:         tc.CodeBg,
		QuoteBg:        tc.QuoteBg,
		Border:         tc.Border,
		SyntaxKeyword:  tc.SyntaxKeyword,
		SyntaxString:   tc.SyntaxString,
		SyntaxComment:  tc.SyntaxComment,
		SyntaxNumber:   tc.SyntaxNumber,
		SyntaxType:     tc.SyntaxType,
		SyntaxBuiltin:  tc.SyntaxBuiltin,
		SyntaxOperator: tc.SyntaxOperator,
	}

	colorWarnings := theme.ValidateColorOverrides(overrides)
	warnings = append(warnings, colorWarnings...)

	tc.Background = overrides.Background
	tc.Text = overrides.Text
	tc.Muted = overrides.Muted
	tc.Heading = overrides.Heading
	tc.Accent = overrides.Accent
	tc.Link = overrides.Link
	tc.LinkURL = overrides.LinkURL
	tc.CodeBg = overrides.CodeBg
	tc.QuoteBg = overrides.QuoteBg
	tc.Border = overrides.Border
	tc.SyntaxKeyword = overrides.SyntaxKeyword
	tc.SyntaxString = overrides.SyntaxString
	tc.SyntaxComment = overrides.SyntaxComment
	tc.SyntaxNumber = overrides.SyntaxNumber
	tc.SyntaxType = overrides.SyntaxType
	tc.SyntaxBuiltin = overrides.SyntaxBuiltin
	tc.SyntaxOperator = overrides.SyntaxOperator

	return warnings
}

func merge(base, overlay Config, meta toml.MetaData) Config {
	if overlay.UI.Padding != 0 {
		base.UI.Padding = overlay.UI.Padding
	}
	if overlay.UI.LineSpacing != 0 {
		base.UI.LineSpacing = overlay.UI.LineSpacing
	}
	if overlay.UI.Scrolloff != 0 {
		base.UI.Scrolloff = overlay.UI.Scrolloff
	}
	if overlay.UI.MaxWidth != 0 {
		base.UI.MaxWidth = overlay.UI.MaxWidth
	}
	if meta.IsDefined("ui", "soft_wrap") {
		base.UI.SoftWrap = overlay.UI.SoftWrap
	}
	if meta.IsDefined("ui", "show_line_numbers") {
		base.UI.ShowLineNumbers = overlay.UI.ShowLineNumbers
	}
	if meta.IsDefined("ui", "show_urls") {
		base.UI.ShowURLs = overlay.UI.ShowURLs
	}

	if overlay.Theme.Name != "" {
		base.Theme.Name = overlay.Theme.Name
	}
	if overlay.Theme.Background != "" {
		base.Theme.Background = overlay.Theme.Background
	}
	if overlay.Theme.Text != "" {
		base.Theme.Text = overlay.Theme.Text
	}
	if overlay.Theme.Muted != "" {
		base.Theme.Muted = overlay.Theme.Muted
	}
	if overlay.Theme.Heading != "" {
		base.Theme.Heading = overlay.Theme.Heading
	}
	if overlay.Theme.Accent != "" {
		base.Theme.Accent = overlay.Theme.Accent
	}
	if overlay.Theme.Link != "" {
		base.Theme.Link = overlay.Theme.Link
	}
	if overlay.Theme.LinkURL != "" {
		base.Theme.LinkURL = overlay.Theme.LinkURL
	}
	if overlay.Theme.CodeBg != "" {
		base.Theme.CodeBg = overlay.Theme.CodeBg
	}
	if overlay.Theme.QuoteBg != "" {
		base.Theme.QuoteBg = overlay.Theme.QuoteBg
	}
	if overlay.Theme.Border != "" {
		base.Theme.Border = overlay.Theme.Border
	}
	if overlay.Theme.SyntaxKeyword != "" {
		base.Theme.SyntaxKeyword = overlay.Theme.SyntaxKeyword
	}
	if overlay.Theme.SyntaxString != "" {
		base.Theme.SyntaxString = overlay.Theme.SyntaxString
	}
	if overlay.Theme.SyntaxComment != "" {
		base.Theme.SyntaxComment = overlay.Theme.SyntaxComment
	}
	if overlay.Theme.SyntaxNumber != "" {
		base.Theme.SyntaxNumber = overlay.Theme.SyntaxNumber
	}
	if overlay.Theme.SyntaxType != "" {
		base.Theme.SyntaxType = overlay.Theme.SyntaxType
	}
	if overlay.Theme.SyntaxBuiltin != "" {
		base.Theme.SyntaxBuiltin = overlay.Theme.SyntaxBuiltin
	}
	if overlay.Theme.SyntaxOperator != "" {
		base.Theme.SyntaxOperator = overlay.Theme.SyntaxOperator
	}

	if meta.IsDefined("markdown", "hide_syntax") {
		base.Markdown.HideSyntax = overlay.Markdown.HideSyntax
	}
	if meta.IsDefined("markdown", "render_emphasis") {
		base.Markdown.RenderEmphasis = overlay.Markdown.RenderEmphasis
	}
	if meta.IsDefined("markdown", "render_strong") {
		base.Markdown.RenderStrong = overlay.Markdown.RenderStrong
	}
	if meta.IsDefined("markdown", "render_links") {
		base.Markdown.RenderLinks = overlay.Markdown.RenderLinks
	}
	if meta.IsDefined("markdown", "render_images") {
		base.Markdown.RenderImages = overlay.Markdown.RenderImages
	}
	if meta.IsDefined("markdown", "render_tables") {
		base.Markdown.RenderTables = overlay.Markdown.RenderTables
	}
	if meta.IsDefined("markdown", "render_task_lists") {
		base.Markdown.RenderTaskLists = overlay.Markdown.RenderTaskLists
	}

	if overlay.Keys.Up != "" {
		base.Keys.Up = overlay.Keys.Up
	}
	if overlay.Keys.Down != "" {
		base.Keys.Down = overlay.Keys.Down
	}
	if overlay.Keys.HalfUp != "" {
		base.Keys.HalfUp = overlay.Keys.HalfUp
	}
	if overlay.Keys.HalfDown != "" {
		base.Keys.HalfDown = overlay.Keys.HalfDown
	}
	if overlay.Keys.Top != "" {
		base.Keys.Top = overlay.Keys.Top
	}
	if overlay.Keys.Bottom != "" {
		base.Keys.Bottom = overlay.Keys.Bottom
	}
	if overlay.Keys.Search != "" {
		base.Keys.Search = overlay.Keys.Search
	}
	if overlay.Keys.NextHit != "" {
		base.Keys.NextHit = overlay.Keys.NextHit
	}
	if overlay.Keys.PrevHit != "" {
		base.Keys.PrevHit = overlay.Keys.PrevHit
	}
	if overlay.Keys.Reload != "" {
		base.Keys.Reload = overlay.Keys.Reload
	}
	if overlay.Keys.Quit != "" {
		base.Keys.Quit = overlay.Keys.Quit
	}
	if overlay.Keys.Help != "" {
		base.Keys.Help = overlay.Keys.Help
	}

	if overlay.Watch.DebounceMs != 0 {
		base.Watch.DebounceMs = overlay.Watch.DebounceMs
	}
	if meta.IsDefined("watch", "enabled") {
		base.Watch.Enabled = overlay.Watch.Enabled
	}

	return base
}

func resolveTheme(tc ThemeConfig) theme.Theme {
	base := theme.Resolve(tc.Name)
	overrides := theme.Theme{
		Background:     tc.Background,
		Text:           tc.Text,
		Muted:          tc.Muted,
		Heading:        tc.Heading,
		Accent:         tc.Accent,
		Link:           tc.Link,
		LinkURL:        tc.LinkURL,
		CodeBg:         tc.CodeBg,
		QuoteBg:        tc.QuoteBg,
		Border:         tc.Border,
		SyntaxKeyword:  tc.SyntaxKeyword,
		SyntaxString:   tc.SyntaxString,
		SyntaxComment:  tc.SyntaxComment,
		SyntaxNumber:   tc.SyntaxNumber,
		SyntaxType:     tc.SyntaxType,
		SyntaxBuiltin:  tc.SyntaxBuiltin,
		SyntaxOperator: tc.SyntaxOperator,
	}
	return theme.Merge(base, overrides)
}