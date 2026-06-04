package theme

import (
	"strings"
	"testing"
)

func TestIsValidColor_EmptyString(t *testing.T) {
	if !isValidColor("") {
		t.Error("empty string should be valid (means use theme default)")
	}
}

func TestIsValidColor_ValidHex6(t *testing.T) {
	cases := []string{"#1e1e2e", "#ffffff", "#000000", "#AABBCC", "#89b4fa"}
	for _, c := range cases {
		if !isValidColor(c) {
			t.Errorf("expected %q to be valid", c)
		}
	}
}

func TestIsValidColor_ValidHex3(t *testing.T) {
	cases := []string{"#fff", "#000", "#abc", "#ABC"}
	for _, c := range cases {
		if !isValidColor(c) {
			t.Errorf("expected %q to be valid", c)
		}
	}
}

func TestIsValidColor_ValidANSI256(t *testing.T) {
	cases := []string{"0", "1", "128", "255"}
	for _, c := range cases {
		if !isValidColor(c) {
			t.Errorf("expected ANSI 256 value %q to be valid", c)
		}
	}
}

func TestIsValidColor_InvalidHex_WrongLength(t *testing.T) {
	cases := []string{"#ff", "#ffff", "#fffff", "#fffffff"}
	for _, c := range cases {
		if isValidColor(c) {
			t.Errorf("expected %q to be invalid (wrong hex length)", c)
		}
	}
}

func TestIsValidColor_InvalidHex_BadChars(t *testing.T) {
	cases := []string{"#gggggg", "#zzzzzz", "#12345g"}
	for _, c := range cases {
		if isValidColor(c) {
			t.Errorf("expected %q to be invalid (bad hex chars)", c)
		}
	}
}

func TestIsValidColor_InvalidANSI_OutOfRange(t *testing.T) {
	cases := []string{"256", "999", "-1"}
	for _, c := range cases {
		if isValidColor(c) {
			t.Errorf("expected %q to be invalid (ANSI out of range)", c)
		}
	}
}

func TestIsValidColor_InvalidStrings(t *testing.T) {
	cases := []string{"red", "blue", "rgb(0,0,0)", "not-a-color", "hsl(0,0%,0%)"}
	for _, c := range cases {
		if isValidColor(c) {
			t.Errorf("expected %q to be invalid", c)
		}
	}
}

func TestIsValidThemeName_Valid(t *testing.T) {
	if !IsValidThemeName("default") {
		t.Error("'default' should be a valid theme name")
	}
	if !IsValidThemeName("light") {
		t.Error("'light' should be a valid theme name")
	}
}

func TestIsValidThemeName_Invalid(t *testing.T) {
	cases := []string{"dark", "monokai", "solarized", "", "DEFAULT", "Light"}
	for _, c := range cases {
		if IsValidThemeName(c) {
			t.Errorf("expected %q to be an invalid theme name", c)
		}
	}
}

func TestValidateThemeName_ValidName_NoWarnings(t *testing.T) {
	name, warnings := ValidateThemeName("light")
	if name != "light" {
		t.Errorf("want 'light', got %q", name)
	}
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %d", len(warnings))
	}
}

func TestValidateThemeName_InvalidName_FallsBackToDefault(t *testing.T) {
	name, warnings := ValidateThemeName("solarized")
	if name != "default" {
		t.Errorf("invalid theme name should fall back to 'default', got %q", name)
	}
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].Field != "theme.name" {
		t.Errorf("warning Field: want 'theme.name', got %q", warnings[0].Field)
	}
	if warnings[0].Value != "solarized" {
		t.Errorf("warning Value: want 'solarized', got %q", warnings[0].Value)
	}
	if !strings.Contains(warnings[0].String(), "theme.name") {
		t.Errorf("warning string should mention 'theme.name', got: %q", warnings[0].String())
	}
}

func TestValidateThemeName_EmptyName_FallsBackToDefault(t *testing.T) {
	name, warnings := ValidateThemeName("")
	if name != "default" {
		t.Errorf("empty theme name should fall back to 'default', got %q", name)
	}
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning for empty name, got %d", len(warnings))
	}
}

func TestValidateColorOverrides_AllValid_NoWarnings(t *testing.T) {
	overrides := &Theme{
		Background: "#1e1e2e",
		Text:       "#cdd6f4",
		Heading:    "#89b4fa",
		Accent:     "128",
	}
	warnings := ValidateColorOverrides(overrides)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings for valid colors, got %d: %v", len(warnings), warnings)
	}
	if overrides.Background != "#1e1e2e" {
		t.Error("valid color should not be cleared")
	}
}

func TestValidateColorOverrides_AllEmpty_NoWarnings(t *testing.T) {
	overrides := &Theme{}
	warnings := ValidateColorOverrides(overrides)
	if len(warnings) != 0 {
		t.Errorf("empty overrides should produce no warnings, got %d", len(warnings))
	}
}

func TestValidateColorOverrides_InvalidColor_ClearedToEmpty(t *testing.T) {
	overrides := &Theme{
		Heading: "not-a-color",
	}
	warnings := ValidateColorOverrides(overrides)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if overrides.Heading != "" {
		t.Errorf("invalid color should be cleared to empty, got %q", overrides.Heading)
	}
	if warnings[0].Field != "theme.heading" {
		t.Errorf("warning Field: want 'theme.heading', got %q", warnings[0].Field)
	}
	if warnings[0].Value != "not-a-color" {
		t.Errorf("warning Value: want 'not-a-color', got %q", warnings[0].Value)
	}
}

func TestValidateColorOverrides_MultipleInvalid(t *testing.T) {
	overrides := &Theme{
		Background: "notvalid",
		Text:       "#zzzzzz",
		Heading:    "#89b4fa",
	}
	warnings := ValidateColorOverrides(overrides)
	if len(warnings) != 2 {
		t.Errorf("expected 2 warnings, got %d", len(warnings))
	}
	if overrides.Background != "" {
		t.Error("invalid background should be cleared")
	}
	if overrides.Text != "" {
		t.Error("invalid text should be cleared")
	}
	if overrides.Heading != "#89b4fa" {
		t.Error("valid heading should not be cleared")
	}
}

func TestValidateColorOverrides_AllColorFields_AreChecked(t *testing.T) {
	overrides := &Theme{
		Background:     "BAD",
		Text:           "BAD",
		Muted:          "BAD",
		Heading:        "BAD",
		Accent:         "BAD",
		Link:           "BAD",
		LinkURL:        "BAD",
		CodeBg:         "BAD",
		QuoteBg:        "BAD",
		Border:         "BAD",
		SyntaxKeyword:  "BAD",
		SyntaxString:   "BAD",
		SyntaxComment:  "BAD",
		SyntaxNumber:   "BAD",
		SyntaxType:     "BAD",
		SyntaxBuiltin:  "BAD",
		SyntaxOperator: "BAD",
	}
	warnings := ValidateColorOverrides(overrides)
	if len(warnings) != 17 {
		t.Errorf("expected 17 warnings (one per color field), got %d", len(warnings))
	}
}

func TestValidationWarning_String_ContainsFieldAndValue(t *testing.T) {
	w := ValidationWarning{
		Field:   "theme.heading",
		Value:   "bad-value",
		Message: "some explanation",
	}
	s := w.String()
	if !strings.Contains(s, "theme.heading") {
		t.Errorf("warning string should contain field name, got: %q", s)
	}
	if !strings.Contains(s, "bad-value") {
		t.Errorf("warning string should contain invalid value, got: %q", s)
	}
	if !strings.Contains(s, "some explanation") {
		t.Errorf("warning string should contain message, got: %q", s)
	}
	if !strings.HasPrefix(s, "rendermd:") {
		t.Errorf("warning string should start with 'rendermd:', got: %q", s)
	}
}

func TestValidateColorOverrides_ValidANSI256Boundary(t *testing.T) {
	overrides := &Theme{
		Background: "0",
		Text:       "255",
	}
	warnings := ValidateColorOverrides(overrides)
	if len(warnings) != 0 {
		t.Errorf("ANSI boundary values 0 and 255 should be valid, got warnings: %v", warnings)
	}
}

func TestValidateColorOverrides_InvalidANSI256_OutOfRange(t *testing.T) {
	overrides := &Theme{
		Background: "256",
	}
	warnings := ValidateColorOverrides(overrides)
	if len(warnings) != 1 {
		t.Errorf("ANSI value 256 should be invalid, got %d warnings", len(warnings))
	}
	if overrides.Background != "" {
		t.Error("out-of-range ANSI value should be cleared")
	}
}
