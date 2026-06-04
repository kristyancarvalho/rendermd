package theme

import (
	"fmt"
	"strconv"
	"strings"
)

type ValidationWarning struct {
	Field   string
	Value   string
	Message string
}

func (w ValidationWarning) String() string {
	return fmt.Sprintf("mdp: config warning: %s: invalid value %q — %s", w.Field, w.Value, w.Message)
}

func isValidColor(s string) bool {
	if s == "" {
		return true
	}
	if strings.HasPrefix(s, "#") {
		hex := s[1:]
		if len(hex) != 3 && len(hex) != 6 {
			return false
		}
		for _, c := range hex {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return false
			}
		}
		return true
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return n >= 0 && n <= 255
}

func IsValidThemeName(name string) bool {
	switch name {
	case "default", "light":
		return true
	default:
		return false
	}
}

func ValidateThemeName(name string) (string, []ValidationWarning) {
	if IsValidThemeName(name) {
		return name, nil
	}
	return "default", []ValidationWarning{{
		Field:   "theme.name",
		Value:   name,
		Message: `unknown theme name; valid values are "default" and "light", falling back to "default"`,
	}}
}

func ValidateColorOverrides(overrides *Theme) []ValidationWarning {
	type entry struct {
		field string
		ptr   *string
	}
	entries := []entry{
		{"theme.background", &overrides.Background},
		{"theme.text", &overrides.Text},
		{"theme.muted", &overrides.Muted},
		{"theme.heading", &overrides.Heading},
		{"theme.accent", &overrides.Accent},
		{"theme.link", &overrides.Link},
		{"theme.link_url", &overrides.LinkURL},
		{"theme.code_bg", &overrides.CodeBg},
		{"theme.quote_bg", &overrides.QuoteBg},
		{"theme.border", &overrides.Border},
		{"theme.syntax_keyword", &overrides.SyntaxKeyword},
		{"theme.syntax_string", &overrides.SyntaxString},
		{"theme.syntax_comment", &overrides.SyntaxComment},
		{"theme.syntax_number", &overrides.SyntaxNumber},
		{"theme.syntax_type", &overrides.SyntaxType},
		{"theme.syntax_builtin", &overrides.SyntaxBuiltin},
		{"theme.syntax_operator", &overrides.SyntaxOperator},
	}
	var warnings []ValidationWarning
	for _, e := range entries {
		if !isValidColor(*e.ptr) {
			warnings = append(warnings, ValidationWarning{
				Field:   e.field,
				Value:   *e.ptr,
				Message: "expected a hex color (#rrggbb or #rgb) or ANSI 256 index (0-255), falling back to theme default",
			})
			*e.ptr = ""
		}
	}
	return warnings
}