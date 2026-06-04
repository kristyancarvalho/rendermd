package config

import "github.com/kristyancarvalho/rendermd/internal/theme"

func defaults() Config {
	return Config{
		UI: UIConfig{
			Padding:         2,
			LineSpacing:     0,
			Scrolloff:       4,
			SoftWrap:        true,
			MaxWidth:        96,
			ShowLineNumbers: false,
			ShowURLs:        false,
			Mouse:           true,
		},
		Theme: ThemeConfig{
			Name: "default",
		},
		Markdown: MarkdownConfig{
			HideSyntax:      true,
			RenderEmphasis:  true,
			RenderStrong:    true,
			RenderLinks:     true,
			RenderImages:    false,
			RenderTables:    true,
			RenderTaskLists: true,
		},
		Keys: KeysConfig{
			Up:       "k",
			Down:     "j",
			HalfUp:   "ctrl+u",
			HalfDown: "ctrl+d",
			Top:      "g",
			Bottom:   "G",
			Search:   "/",
			NextHit:  "n",
			PrevHit:  "N",
			Reload:   "r",
			Quit:     "q",
			Help:     "?",
		},
		Watch: WatchConfig{
			Enabled:    true,
			DebounceMs: 150,
		},
		resolved: theme.Default,
	}
}
