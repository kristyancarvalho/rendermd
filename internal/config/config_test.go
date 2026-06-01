package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaults(t *testing.T) {
	cfg := defaults()

	if cfg.UI.Padding != 2 {
		t.Errorf("Padding: want 2, got %d", cfg.UI.Padding)
	}
	if cfg.UI.MaxWidth != 96 {
		t.Errorf("MaxWidth: want 96, got %d", cfg.UI.MaxWidth)
	}
	if !cfg.UI.SoftWrap {
		t.Error("SoftWrap: want true")
	}
	if cfg.UI.Scrolloff != 4 {
		t.Errorf("Scrolloff: want 4, got %d", cfg.UI.Scrolloff)
	}
	if cfg.Theme.Name != "default" {
		t.Errorf("Theme.Name: want 'default', got %q", cfg.Theme.Name)
	}
	if cfg.Keys.Up != "k" {
		t.Errorf("Keys.Up: want 'k', got %q", cfg.Keys.Up)
	}
	if cfg.Keys.Down != "j" {
		t.Errorf("Keys.Down: want 'j', got %q", cfg.Keys.Down)
	}
	if cfg.Watch.DebounceMs != 150 {
		t.Errorf("Watch.DebounceMs: want 150, got %d", cfg.Watch.DebounceMs)
	}
	if !cfg.Watch.Enabled {
		t.Error("Watch.Enabled: want true")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	cfg := Load("/nonexistent/config.toml")
	if cfg.UI.Padding != 2 {
		t.Errorf("expected default padding, got %d", cfg.UI.Padding)
	}
}

func TestLoad_ValidTOML(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")

	content := `
[ui]
padding   = 4
max_width = 120
soft_wrap = false

[theme]
name = "light"

[keys]
quit = "x"

[watch]
enabled     = false
debounce_ms = 300
`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := Load(cfgPath)

	if cfg.UI.Padding != 4 {
		t.Errorf("Padding: want 4, got %d", cfg.UI.Padding)
	}
	if cfg.UI.MaxWidth != 120 {
		t.Errorf("MaxWidth: want 120, got %d", cfg.UI.MaxWidth)
	}
	if cfg.UI.SoftWrap {
		t.Error("SoftWrap: want false")
	}
	if cfg.Theme.Name != "light" {
		t.Errorf("Theme.Name: want 'light', got %q", cfg.Theme.Name)
	}
	if cfg.Keys.Quit != "x" {
		t.Errorf("Keys.Quit: want 'x', got %q", cfg.Keys.Quit)
	}
	if cfg.Watch.Enabled {
		t.Error("Watch.Enabled: want false")
	}
	if cfg.Watch.DebounceMs != 300 {
		t.Errorf("Watch.DebounceMs: want 300, got %d", cfg.Watch.DebounceMs)
	}
}

func TestLoad_InvalidTOML(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "bad.toml")
	if err := os.WriteFile(cfgPath, []byte(":::invalid toml:::"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := Load(cfgPath)
	if cfg.UI.Padding != 2 {
		t.Errorf("expected default padding after parse error, got %d", cfg.UI.Padding)
	}
}

func TestLoad_ThemeOverride(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	content := `
[theme]
name    = "default"
heading = "#ff0000"
`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := Load(cfgPath)
	
	thm := cfg.ResolvedTheme()
	if thm.Heading != "#ff0000" {
		t.Errorf("Heading override: want '#ff0000', got %q", thm.Heading)
	}
}

func TestLoad_PartialKeys(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	content := `
[keys]
up = "w"
`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := Load(cfgPath)
	if cfg.Keys.Up != "w" {
		t.Errorf("Keys.Up: want 'w', got %q", cfg.Keys.Up)
	}
	if cfg.Keys.Down != "j" {
		t.Errorf("Keys.Down: want default 'j', got %q", cfg.Keys.Down)
	}
}

func TestMerge_PreservesUnsetFields(t *testing.T) {
	base := defaults()
	overlay := Config{}
	result := merge(base, overlay)

	if result.UI.Padding != base.UI.Padding {
		t.Errorf("Padding should be preserved: want %d, got %d", base.UI.Padding, result.UI.Padding)
	}
	if result.Keys.Up != base.Keys.Up {
		t.Errorf("Keys.Up should be preserved: want %q, got %q", base.Keys.Up, result.Keys.Up)
	}
}

func TestDefaultPath(t *testing.T) {
	p := DefaultPath()
	if p == "" {
		t.Error("DefaultPath should not be empty")
	}
}