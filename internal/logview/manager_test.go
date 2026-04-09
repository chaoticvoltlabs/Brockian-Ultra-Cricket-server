package logview

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestManagerWritesAndBroadcasts(t *testing.T) {
	dir := t.TempDir()

	m, err := NewManager(dir)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	id, ch, recent := m.Subscribe()
	defer m.Unsubscribe(id)

	if len(recent) != 0 {
		t.Fatalf("expected empty recent buffer, got %d", len(recent))
	}

	if _, err := m.Write([]byte("hello world\n")); err != nil {
		t.Fatalf("Write: %v", err)
	}

	select {
	case got := <-ch:
		if got != "hello world\n" {
			t.Fatalf("unexpected broadcast: %q", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for broadcast")
	}

	files, err := m.Files()
	if err != nil {
		t.Fatalf("Files: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 log file, got %d", len(files))
	}

	data, err := os.ReadFile(filepath.Join(dir, files[0].Name))
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	if string(data) != "hello world\n" {
		t.Fatalf("unexpected log file contents: %q", string(data))
	}
}
