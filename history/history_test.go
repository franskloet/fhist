package history

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBashHistoryLoader(t *testing.T) {
	tmpDir := t.TempDir()
	histFile := filepath.Join(tmpDir, ".bash_history")

	content := "#1724128500\ncd /tmp\n#1724128510\ngit status\nls -la\n"
	err := os.WriteFile(histFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write mock history file: %v", err)
	}

	items, err := LoadBashHistory(histFile)
	if err != nil {
		t.Fatalf("LoadBashHistory error: %v", err)
	}

	if len(items) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(items))
	}

	if items[0].Command != "cd /tmp" {
		t.Errorf("Expected 'cd /tmp', got '%s'", items[0].Command)
	}
	if items[0].Index != 0 {
		t.Errorf("Expected index 0, got %d", items[0].Index)
	}

	if items[1].Command != "git status" {
		t.Errorf("Expected 'git status', got '%s'", items[1].Command)
	}

	if items[2].Command != "ls -la" {
		t.Errorf("Expected 'ls -la', got '%s'", items[2].Command)
	}
	if items[2].Index != 2 {
		t.Errorf("Expected index 2, got %d", items[2].Index)
	}
}

func TestZshHistoryLoader(t *testing.T) {
	tmpDir := t.TempDir()
	histFile := filepath.Join(tmpDir, ".zsh_history")

	content := ": 1724128500:0;mkdir test\n: 1724128510:0;cd test\n"
	err := os.WriteFile(histFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write mock zsh history: %v", err)
	}

	items, err := LoadZshHistory(histFile)
	if err != nil {
		t.Fatalf("LoadZshHistory error: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(items))
	}
	if items[0].Command != "mkdir test" {
		t.Errorf("Expected 'mkdir test', got '%s'", items[0].Command)
	}
}
