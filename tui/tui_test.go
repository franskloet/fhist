package tui

import (
	"fmt"
	"strings"
	"testing"

	"fhist/history"
)

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

	// Test index 0
	m.SelectedIndex = 0
	view0 := m.View()
	lines0 := strings.Split(view0, "\n")
	fmt.Printf("SelectedIndex 0 line count: %d\n", len(lines0))

	// Test index 1
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
