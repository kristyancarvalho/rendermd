package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun_Version(t *testing.T) {
	err := Run("1.2.3", []string{"--version"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestRun_Help(t *testing.T) {
	err := Run("dev", []string{"--help"})
	if err != nil {
		t.Fatalf("expected no error on --help, got: %v", err)
	}
}

func TestRun_FileNotFound(t *testing.T) {
	err := Run("dev", []string{"/nonexistent/file.md"})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRun_FileRead(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.md")
	if err := os.WriteFile(f, []byte("# Hello\n\nworld\n"), 0644); err != nil {
		t.Fatal(err)
	}

	err := Run("dev", []string{"--version"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_UnknownFlag(t *testing.T) {
	err := Run("dev", []string{"--no-such-flag"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestRun_ThemeFlag(t *testing.T) {
	err := Run("dev", []string{"--theme", "light", "--version"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}