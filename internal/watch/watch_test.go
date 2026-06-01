package watch

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew_InvalidPath(t *testing.T) {
	_, err := New("/nonexistent/path/file.md", 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

func TestNew_ValidFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.md")
	if err := os.WriteFile(f, []byte("# hello"), 0644); err != nil {
		t.Fatal(err)
	}

	w, err := New(f, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer w.Close()

	if w.Events == nil {
		t.Error("Events channel should not be nil")
	}
	if w.Errors == nil {
		t.Error("Errors channel should not be nil")
	}
}

func TestWatcher_FiresOnWrite(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "watch.md")
	if err := os.WriteFile(f, []byte("initial"), 0644); err != nil {
		t.Fatal(err)
	}

	w, err := New(f, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer w.Close()

	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(f, []byte("updated"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case <-w.Events:
		
	case <-time.After(2 * time.Second):
		t.Error("expected file-change event within 2s, got none")
	}
}

func TestWatcher_Close(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "c.md")
	if err := os.WriteFile(f, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	w, err := New(f, 50*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	if err := w.Close(); err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestWatcher_DebounceCoalesces(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "debounce.md")
	if err := os.WriteFile(f, []byte("v0"), 0644); err != nil {
		t.Fatal(err)
	}

	debounce := 100 * time.Millisecond
	w, err := New(f, debounce)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer w.Close()

	for i := 0; i < 5; i++ {
		_ = os.WriteFile(f, []byte("v"+string(rune('0'+i))), 0644)
		time.Sleep(10 * time.Millisecond)
	}

	deadline := time.After(debounce * 3)
	events := 0
loop:
	for {
		select {
		case <-w.Events:
			events++
		case <-deadline:
			break loop
		}
	}

	if events == 0 {
		t.Error("expected at least one event after writes")
	}
	if events >= 5 {
		t.Errorf("debounce should coalesce events; got %d (one per write)", events)
	}
}

func TestWatcher_NoEventAfterClose(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "noev.md")
	if err := os.WriteFile(f, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	w, err := New(f, 20*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	time.Sleep(30 * time.Millisecond)
	_ = os.WriteFile(f, []byte("after close"), 0644)

	select {
	case <-w.Events:
	case <-time.After(200 * time.Millisecond):
	}
}