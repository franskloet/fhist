package tui

import (
	"fmt"
	"strings"
	"testing"

	"fhist/history"
	tea "github.com/charmbracelet/bubbletea"
)

func TestHorizontalScrollingAndKeyNav(t *testing.T) {
	store := &history.HistoryStore{
		BashItems: []history.HistoryItem{
			{Index: 0, Command: "git commit -m 'very long command string line test'", Source: "bash"},
			{Index: 1, Command: "docker run -it ubuntu bash", Source: "bash"},
		},
		ActiveSource: "bash",
	}
	store.RebuildAll()

	m := NewModel(store, "", 5)

	// Test Right Key
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	model2 := m2.(Model)
	if model2.CommandOffset != 5 {
		t.Errorf("Expected CommandOffset 5 after Right key, got %d", model2.CommandOffset)
	}

	// Test Left Key
	m3, _ := model2.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model3 := m3.(Model)
	if model3.CommandOffset != 0 {
		t.Errorf("Expected CommandOffset 0 after Left key, got %d", model3.CommandOffset)
	}

	// Test Alt+Up Key for Context Radius
	m4, _ := model3.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'='}})
	model4 := m4.(Model)
	if model4.ContextRadius != 6 {
		t.Errorf("Expected ContextRadius 6 after '=' key, got %d", model4.ContextRadius)
	}

	fmt.Println("Horizontal scroll & key nav test passed cleanly!")
}

func TestChrootSearchAccuracy(t *testing.T) {
	store := &history.HistoryStore{
		BashItems: []history.HistoryItem{
			{Index: 0, Command: "cd /home/frans/Development", Source: "bash"},
			{Index: 1, Command: "chroot /mnt /bin/bash", Source: "bash"},
			{Index: 2, Command: "sudo chroot /target", Source: "bash"},
			{Index: 3, Command: "cat /etc/hosts", Source: "bash"},
			{Index: 4, Command: "echo 'hello world'", Source: "bash"},
		},
		ActiveSource: "bash",
	}
	store.RebuildAll()

	m := NewModel(store, "chroot", 5)
	if len(m.FilteredIndexes) != 2 {
		t.Fatalf("Expected exactly 2 hits for 'chroot', got %d", len(m.FilteredIndexes))
	}

	hit1 := m.CurrentItems[m.FilteredIndexes[0]].Command
	hit2 := m.CurrentItems[m.FilteredIndexes[1]].Command

	if hit1 != "sudo chroot /target" {
		t.Errorf("Expected first hit 'sudo chroot /target', got '%s'", hit1)
	}
	if hit2 != "chroot /mnt /bin/bash" {
		t.Errorf("Expected second hit 'chroot /mnt /bin/bash', got '%s'", hit2)
	}

	fmt.Println("Chroot search test passed cleanly!")
}

func TestLongCommandLineWrapping(t *testing.T) {
	store := &history.HistoryStore{
		BashItems: []history.HistoryItem{
			{Index: 0, Command: "git commit -m 'this is a very very very long command line string that might cause wrapping if width is calculated incorrectly'", Source: "bash"},
			{Index: 1, Command: "docker run -it --rm -v /home/frans/Development/tools/fhist:/app -w /app ubuntu:latest bash -c 'echo test && ls -la'", Source: "bash"},
		},
		ActiveSource: "bash",
	}
	store.RebuildAll()

	m := NewModel(store, "", 5)
	m.Width = 80
	m.Height = 24

	m.SelectedIndex = 0
	view0 := m.View()
	lines0 := strings.Split(view0, "\n")
	fmt.Printf("SelectedIndex 0 line count: %d\n", len(lines0))

	m.SelectedIndex = 1
	view1 := m.View()
	lines1 := strings.Split(view1, "\n")
	fmt.Printf("SelectedIndex 1 line count: %d\n", len(lines1))

	if len(lines0) > 24 {
		t.Errorf("SelectedIndex 0 exceeded height 24! Got %d lines", len(lines0))
	}
	if len(lines1) > 24 {
		t.Errorf("SelectedIndex 1 exceeded height 24! Got %d lines", len(lines1))
	}
}
