package app

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_Version(t *testing.T) {
	var out bytes.Buffer
	build := BuildInfo{
		Version: "v1.2.3",
		Commit:  "abc1234",
		Date:    "2026-06-03T00:00:00Z",
	}
	err := runWithStdout(&out, func() error {
		return Run(build, []string{"--version"})
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	for _, want := range []string{"mdp v1.2.3", "commit abc1234", "built 2026-06-03T00:00:00Z"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("version output missing %q: %q", want, out.String())
		}
	}
}

func TestRun_Help(t *testing.T) {
	err := Run(BuildInfo{Version: "dev"}, []string{"--help"})
	if err != nil {
		t.Fatalf("expected no error on --help, got: %v", err)
	}
}

func TestRun_FileNotFound(t *testing.T) {
	err := Run(BuildInfo{Version: "dev"}, []string{"/nonexistent/file.md"})
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

	err := Run(BuildInfo{Version: "dev"}, []string{"--version"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_UnknownFlag(t *testing.T) {
	err := Run(BuildInfo{Version: "dev"}, []string{"--no-such-flag"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestRun_ThemeFlag(t *testing.T) {
	err := Run(BuildInfo{Version: "dev"}, []string{"--theme", "light", "--version"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func runWithStdout(out io.Writer, fn func() error) error {
	original := os.Stdout
	read, write, err := os.Pipe()
	if err != nil {
		return err
	}
	os.Stdout = write
	defer func() {
		os.Stdout = original
	}()

	result := fn()
	if err := write.Close(); err != nil && result == nil {
		result = err
	}
	if _, err := io.Copy(out, read); err != nil && result == nil {
		result = err
	}
	_ = read.Close()
	return result
}
